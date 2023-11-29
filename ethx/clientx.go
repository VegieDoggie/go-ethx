package ethx

import (
	"context"
	"errors"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	"golang.org/x/time/rate"
	"log"
	"math/big"
	"reflect"
	"sync"
	"time"
)

// Clientx defines typed wrappers for the Ethereum RPC API of a set of the Ethereum Clients.
type Clientx struct {
	it             *Iterator[*ethclient.Client]
	ctx            context.Context
	rpcMap         map[*ethclient.Client]string
	rpcErrCountMap map[*ethclient.Client]uint
	notFoundBlocks uint64
	chainId        *big.Int
	latestHeader   *types.Header
	stop           chan bool
	miningInterval time.Duration
}

// NewSimpleClientx create *Clientx
// concurrency is the concurrency per seconds of any rpc, default 1/s
// see NewClientx
func NewSimpleClientx(rpcList []string, concurrency ...int) *Clientx {
	_concurrency := 1
	if len(concurrency) > 0 && concurrency[0] > 1 {
		_concurrency = concurrency[0]
	}
	var weights []int
	for range rpcList {
		weights = append(weights, _concurrency)
	}
	return NewClientx(rpcList, weights)
}

// NewClientx connects clients to the given URLs, to provide a reliable Ethereum RPC API call, includes
// a timer to regularly update block height(AutoBlockNumber).
// If weight <= 1, the weight is always 1.
// Note: If len(weightList) == 0, then default weight = 1 will be active.
func NewClientx(rpcList []string, weights []int, limiter ...*rate.Limiter) *Clientx {
	iterator, rpcMap, chainId := newClientIteratorWithWeight(rpcList, weights, limiter...)
	c := &Clientx{
		ctx:            context.Background(),
		it:             iterator,
		rpcMap:         rpcMap,
		chainId:        chainId,
		rpcErrCountMap: make(map[*ethclient.Client]uint),
		notFoundBlocks: uint64(len(rpcMap) * 2),
		latestHeader:   &types.Header{Number: BigInt(0)},
	}
	c.startBackground()
	return c
}

// newClientIteratorWithWeight creates a clientIterator with wights.
func newClientIteratorWithWeight(rpcList []string, weightList []int, limiter ...*rate.Limiter) (clientIterator *Iterator[*ethclient.Client], rpcMap map[*ethclient.Client]string, latestChainId *big.Int) {
	if len(rpcList) != len(weightList) {
		tmp := weightList
		weightList = make([]int, len(rpcList))
		copy(weightList, tmp)
	}
	var reliableClients []*ethclient.Client
	rpcMap = make(map[*ethclient.Client]string)
	for i, rpc := range rpcList {
		client, err := ethclient.Dial(rpc)
		if err == nil {
			chainId, err := client.ChainID(context.TODO())
			if err == nil {
				if latestChainId == nil {
					latestChainId = chainId
				} else if latestChainId.Cmp(chainId) != 0 {
					panic(errors.New(fmt.Sprintf("[ERROR] rpc(%v) chainID is %v,but rpc(%v) chainId is %v\n", rpcList[i-1], latestChainId, rpc, chainId)))
				}
				reliableClients = append(reliableClients, client)
				for j := 1; j < weightList[i]; j++ {
					reliableClients = append(reliableClients, client)
				}
				rpcMap[client] = rpc
				continue
			}
		}
		log.Printf("[WARN] newClientIteratorWithWeight::Unreliable rpc: %v\n", rpc)
	}
	if len(reliableClients) == 0 {
		panic(errors.New(fmt.Sprintf("[ERROR] newClientIteratorWithWeight::Unreliable rpc List: %v\n", rpcList)))
	}
	clientIterator = NewIterator[*ethclient.Client](reliableClients, limiter...).Shuffle()
	return clientIterator, rpcMap, latestChainId
}

func (c *Clientx) NextClient() *ethclient.Client {
	return c.it.WaitNext()
}

func (c *Clientx) GetRPCs() (rpcList []string) {
	for _, v := range c.rpcMap {
		rpcList = append(rpcList, v)
	}
	return rpcList
}

func (c *Clientx) UpdateRPCs(newRPCs []string) {
	if len(newRPCs) == 0 {
		return
	}
	oldRPCs := c.GetRPCs()
	newSet := mapset.NewThreadUnsafeSet[string](newRPCs...)
	oldSet := mapset.NewThreadUnsafeSet[string](oldRPCs...)
	// rpc in newRPCs but not in oldRPCs: ADD
	for _, rpc := range newRPCs {
		if !oldSet.Contains(rpc) {
			client, err := ethclient.Dial(rpc)
			if err == nil {
				chainId, err := client.ChainID(context.TODO())
				if err == nil {
					if c.chainId.Cmp(chainId) != 0 {
						log.Println(fmt.Sprintf("[WARN] UpdateRPCs::required chainID is %v,but rpc(%v) chainId is %v\n", c.chainId, rpc, chainId))
						continue
					}
					c.rpcMap[client] = rpc
					c.it.Add(client)
					continue
				}
			}
			log.Printf("[WARN] UpdateRPCs::Unreliable rpc: %v\n", rpc)
		}
	}
	// rpc in oldRPCs but not in newRPCs: REMOVE
	for _, rpc := range oldRPCs {
		if !newSet.Contains(rpc) {
			var client *ethclient.Client
			for k, v := range c.rpcMap {
				if v == rpc {
					client = k
					break
				}
			}
			if c.it.Len() > 1 {
				delete(c.rpcMap, client)
				c.it.Remove(client)
			}
		}
	}
	log.Printf("[DONE] UpdateRPCs::from=%v, to=%v, params=%v\n", oldRPCs, c.GetRPCs(), newRPCs)
}

func (c *Clientx) record(f any, client *ethclient.Client, err error) {
	c.rpcErrCountMap[client]++
	log.Printf("%v [WARN] func=%v, rpc=%v #%v, err=%v\r\n", time.Now().Format(time.DateTime), getFuncName(f), c.rpcMap[client], c.rpcErrCountMap[client], err)
}

type MustContract struct {
	client          *Clientx
	constructor     any
	contractName    string
	contractAddress common.Address
	maxErrNum       int
	config          EventConfig
}

func (c *Clientx) NewMustContract(constructor any, addressLike any, config EventConfig, maxErrNum ...int) *MustContract {
	config.panicIfNotValid()
	_maxErrNum := 99
	if len(maxErrNum) > 0 {
		_maxErrNum = maxErrNum[0]
	}
	return &MustContract{
		client:          c,
		constructor:     constructor,
		contractName:    getFuncName(constructor)[3:],
		contractAddress: Address(addressLike),
		maxErrNum:       _maxErrNum,
		config:          config,
	}
}

func (m *MustContract) Run0(f any, args ...any) any {
	return m.Run(f, args...)[0]
}

func (m *MustContract) Run(f any, args ...any) []any {
	for i := 0; i < m.maxErrNum; i++ {
		ret, err := m.exec(f, args...)
		if err != nil {
			continue
		}
		return ret
	}
	panic(errors.New(fmt.Sprintf("%v::exceed maxErrNum(%v), contract=%v", getFuncName(f), m.maxErrNum, m.contractAddress)))
}

func (m *MustContract) exec(f any, args ...any) ([]any, error) {
	client := m.client.it.WaitNext()
	instanceRet := callFunc(m.constructor, m.contractAddress, client)
	if err, ok := instanceRet[len(instanceRet)-1].(error); ok {
		panic(err)
	}
	ret := callStructMethod(instanceRet[0], f, args...)
	if i := len(ret) - 1; i >= 0 {
		if err, ok := ret[i].(error); ok {
			m.client.record(f, client, err)
			return nil, err
		}
		return ret[:i], nil
	}
	return nil, nil
}

// Subscribe sink chan<- *MarketRouterUpdatePosition
func (m *MustContract) Subscribe(from any, ch any, index ...[]any) (sub event.Subscription) {
	// `chan<- *MarketRouterUpdatePosition` => "UpdatePosition"
	eventName := reflect.TypeOf(ch).Elem().Elem().Name()[len(m.contractName):]
	filterFcName := "Filter" + eventName
	chMiddle := make(chan any, 128)
	ignoreCloseChannelPanic := func() {
		if err, ok := recover().(error); ok {
			if err.Error() != "send on closed channel" {
				panic(err)
			}
		}
	}
	filterFc := func(from, to uint64) {
		defer ignoreCloseChannelPanic()
		log.Println("filterFc:", from, to)
		opts := &bind.FilterOpts{
			Start:   from,
			End:     &to,
			Context: m.client.ctx,
		}
		// ethereum/go-ethereum@v1.13.5/accounts/abi/bind/base.go:434:FilterLogs(opts, name, query...)
		var iterator any
		switch len(index) {
		case 0:
			iterator = m.Run(filterFcName, opts)[0]
		case 1:
			iterator = m.Run(filterFcName, opts, index[0])[0]
		case 2:
			iterator = m.Run(filterFcName, opts, index[0], index[1])[0]
		case 3:
			iterator = m.Run(filterFcName, opts, index[0], index[1], index[2])[0]
		default:
			panic(errors.New("Subscribe::the number of event-indexes don't exceed 3.\n"))
		}
		// *TestLogIndex0Iterator, error
		itValue := reflect.ValueOf(iterator)
		itNext := itValue.MethodByName("Next")
		itEvent := itValue.Elem().FieldByName("Event")
		for {
			//tl := itEvent.FieldByName("Raw").Interface().(types.Log)
			if itNext.Call(nil)[0].Interface().(bool) {
				chMiddle <- itEvent.Interface()
			} else {
				return
			}
		}
	}
	stop := make(chan bool, 1)
	go func() {
		tick := time.NewTicker(m.client.miningInterval)
		defer tick.Stop()
		_from, _to := BigInt(from).Uint64(), m.client.BlockNumber()
		for {
			_from = segmentCallback(_from, _to, m.config, filterFc)
			_to = m.client.BlockNumber()
			log.Println("from:", _from, "_to:", _to)
			select {
			case <-tick.C:
			case <-stop:
				return
			}
		}
	}()
	chValue := reflect.ValueOf(ch)
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer ignoreCloseChannelPanic()
		defer close(chMiddle)
		for {
			select {
			case l := <-chMiddle:
				chValue.Send(reflect.ValueOf(l))
			case <-quit:
				stop <- true
				return nil
			}
		}
	})
}

// TransactOpts create *bind.TransactOpts, and panic if privateKey err
// privateKeyLike eg: 0xf1...3, f1...3, []byte, *ecdsa.PrivateKey...
func (c *Clientx) TransactOpts(privateKeyLike any) *bind.TransactOpts {
	_privateKey := PrivateKey(privateKeyLike)
	opts, err := bind.NewKeyedTransactorWithChainID(_privateKey, c.chainId)
	if err != nil {
		panic(err)
	}
	if c.latestHeader.BaseFee != nil {
		gasTipCap := c.SuggestGasTipCap()
		opts.GasFeeCap = Add(Mul(c.latestHeader.BaseFee, 2), gasTipCap)
		opts.GasTipCap = gasTipCap
	} else {
		opts.GasPrice = c.SuggestGasPrice()
	}
	opts.GasLimit = 8000000
	opts.Nonce = BigInt(c.PendingNonceAt(opts.From))
	return opts
}

type TransferOption struct {
	data       []byte             // option
	AccessList types.AccessList   // option
	Opts       *bind.TransactOpts // option
}

// Transfer build transaction and send
// TransferOption is optional.
// see more: github.com/ethereum/go-ethereum/internal/ethapi/transaction_args.go:284
func (c *Clientx) Transfer(privateKeyLike, to, amount any, options ...TransferOption) (tx *types.Transaction, receipt *types.Receipt, err error) {
	_privateKey := PrivateKey(privateKeyLike)
	var data []byte
	var accessList types.AccessList
	var opts *bind.TransactOpts
	if len(options) > 0 {
		data = options[0].data
		accessList = options[0].AccessList
		opts = options[0].Opts
	}
	if opts == nil {
		opts = c.TransactOpts(_privateKey)
	}
	_to := Address(to)
	switch {
	case opts.GasFeeCap != nil:
		tx = types.NewTx(&types.DynamicFeeTx{
			To:         &_to,
			ChainID:    c.chainId,
			Nonce:      opts.Nonce.Uint64(),
			Gas:        opts.GasLimit,
			GasFeeCap:  opts.GasFeeCap,
			GasTipCap:  opts.GasTipCap,
			Value:      BigInt(amount),
			Data:       data,
			AccessList: accessList,
		})
	case accessList != nil:
		tx = types.NewTx(&types.AccessListTx{
			To:         &_to,
			ChainID:    c.chainId,
			Nonce:      opts.Nonce.Uint64(),
			Gas:        opts.GasLimit,
			GasPrice:   opts.GasPrice,
			Value:      BigInt(amount),
			Data:       data,
			AccessList: accessList,
		})
	default:
		tx = types.NewTx(&types.LegacyTx{
			To:       &_to,
			Nonce:    opts.Nonce.Uint64(),
			Gas:      opts.GasLimit,
			GasPrice: opts.GasPrice,
			Value:    BigInt(amount),
			Data:     data,
		})
	}
	if tx, err = opts.Signer(opts.From, tx); err != nil {
		return nil, nil, err
	}
	gasLimit, err := c.EstimateGas(CallMsg(opts.From, tx), len(c.rpcMap))
	if err != nil {
		return nil, nil, err
	}
	if opts.GasLimit < gasLimit || gasLimit == 0 {
		return tx, nil, errors.New(fmt.Sprintf("[ERROR] Transfer::Gas required %v, but %v.\n", gasLimit, opts.GasLimit))
	}
	if err = c.SendTransaction(tx, len(c.rpcMap)); err != nil {
		return nil, nil, err
	}
	if receipt, err = c.WaitMined(tx, 1); err != nil {
		return nil, nil, err
	}
	return tx, receipt, nil
}

func (c *Clientx) startBackground() {
	beforeNumber := c.BlockNumber()
	go func() {
		queryTicker := time.NewTicker(time.Second)
		defer queryTicker.Stop()
		beforeTime := time.Now().Add(-time.Second)
		for {
			header, err := c.HeaderByNumber(nil, 0)
			if err == nil && header.Number.Cmp(c.latestHeader.Number) > 0 {
				c.miningInterval = time.Since(beforeTime)
				c.latestHeader = header
				beforeTime = time.Now()
			}
			<-queryTicker.C
		}
	}()
	queryTicker := time.NewTicker(100 * time.Millisecond)
	defer queryTicker.Stop()
	for beforeNumber >= c.BlockNumber() {
		<-queryTicker.C
	}
}

// BlockNumber returns the most recent block number
func (c *Clientx) BlockNumber() (blockNumber uint64) {
	return c.latestHeader.Number.Uint64()
}

// ChainID retrieves the current chain ID for transaction replay protection.
func (c *Clientx) ChainID() (chainID *big.Int) {
	return c.chainId
}

// NetworkID returns the network ID.
func (c *Clientx) NetworkID() (networkID *big.Int) {
	var err error
	for {
		client := c.it.WaitNext()
		networkID, err = client.NetworkID(c.ctx)
		if err != nil {
			c.record(client.NetworkID, client, err)
			continue
		}
		return
	}
}

// BalanceAt returns the wei balance of the given account.
// The block number can be nil, in which case the balance is taken from the latest known block.
func (c *Clientx) BalanceAt(account any, blockNumber ...any) (balance *big.Int) {
	var _blockNumber *big.Int
	if len(blockNumber) > 0 {
		_blockNumber = BigInt(blockNumber[0])
	}
	_account := Address(account)
	for {
		client := c.it.WaitNext()
		balance, err := client.BalanceAt(c.ctx, _account, _blockNumber)
		if err != nil {
			c.record(client.BalanceAt, client, err)
			continue
		}
		return balance
	}
}

// PendingBalanceAt returns the wei balance of the given account in the pending state.
func (c *Clientx) PendingBalanceAt(account any) (balance *big.Int) {
	_account := Address(account)
	for {
		client := c.it.WaitNext()
		balance, err := client.PendingBalanceAt(c.ctx, _account)
		if err != nil {
			c.record(client.PendingBalanceAt, client, err)
			continue
		}
		return balance
	}
}

// NonceAt returns the account nonce of the given account.
// The block number can be nil, in which case the nonce is taken from the latest known block.
func (c *Clientx) NonceAt(account any, blockNumber ...any) (nonce uint64) {
	var _blockNumber *big.Int
	if len(blockNumber) > 0 {
		_blockNumber = BigInt(blockNumber[0])
	}
	_account := Address(account)
	for {
		client := c.it.WaitNext()
		nonce, err := client.NonceAt(c.ctx, _account, _blockNumber)
		if err != nil {
			c.record(client.NonceAt, client, err)
			continue
		}
		return nonce
	}
}

// PendingNonceAt returns the account nonce of the given account in the pending state.
// This is the nonce that should be used for the next transaction.
func (c *Clientx) PendingNonceAt(account any) (nonce uint64) {
	_account := Address(account)
	for {
		client := c.it.WaitNext()
		nonce, err := client.PendingNonceAt(c.ctx, _account)
		if err != nil {
			c.record(client.PendingNonceAt, client, err)
			continue
		}
		return nonce
	}
}

// FilterLogs executes a filter query.
func (c *Clientx) FilterLogs(q ethereum.FilterQuery) (logs []types.Log) {
	for {
		client := c.it.WaitNext()
		logs, err := client.FilterLogs(c.ctx, q)
		if err != nil {
			c.record(client.FilterLogs, client, err)
			continue
		}
		return logs
	}
}

// SuggestGasPrice retrieves the currently suggested gas price to allow a timely
// execution of a transaction.
func (c *Clientx) SuggestGasPrice() (gasPrice *big.Int) {
	for {
		client := c.it.WaitNext()
		gasPrice, err := client.SuggestGasPrice(c.ctx)
		if err != nil {
			c.record(client.SuggestGasPrice, client, err)
			continue
		}
		return gasPrice
	}
}

// SuggestGasTipCap retrieves the currently suggested gas tip cap after 1559 to
// allow a timely execution of a transaction.
func (c *Clientx) SuggestGasTipCap() (gasTipCap *big.Int) {
	for {
		client := c.it.WaitNext()
		gasTipCap, err := client.SuggestGasTipCap(c.ctx)
		if err != nil {
			c.record(client.SuggestGasTipCap, client, err)
			continue
		}
		return gasTipCap
	}
}

// FeeHistory retrieves the fee market history.
func (c *Clientx) FeeHistory(blockCount uint64, lastBlock any, rewardPercentiles []float64) (feeHistory *ethereum.FeeHistory) {
	for {
		client := c.it.WaitNext()
		feeHistory, err := client.FeeHistory(c.ctx, blockCount, BigInt(lastBlock), rewardPercentiles)
		if err != nil {
			c.record(client.FeeHistory, client, err)
			continue
		}
		return feeHistory
	}
}

// StorageAt returns the value of key in the contract storage of the given account.
// The block number can be nil, in which case the value is taken from the latest known block.
func (c *Clientx) StorageAt(account, keyHash any, blockNumber ...any) (storage []byte) {
	var _blockNumber *big.Int
	if len(blockNumber) > 0 {
		_blockNumber = BigInt(blockNumber[0])
	}
	_account, _keyHash := Address(account), Hash(keyHash)
	for {
		client := c.it.WaitNext()
		storage, err := client.StorageAt(c.ctx, _account, _keyHash, _blockNumber)
		if err != nil {
			c.record(client.StorageAt, client, err)
			continue
		}
		return storage
	}
}

// PendingStorageAt returns the value of key in the contract storage of the given account in the pending state.
func (c *Clientx) PendingStorageAt(account, keyHash any) (storage []byte) {
	_account, _keyHash := Address(account), Hash(keyHash)
	for {
		client := c.it.WaitNext()
		storage, err := client.PendingStorageAt(c.ctx, _account, _keyHash)
		if err != nil {
			c.record(client.PendingStorageAt, client, err)
			continue
		}
		return storage
	}
}

// CodeAt returns the contract code of the given account.
// The block number can be nil, in which case the code is taken from the latest known block.
func (c *Clientx) CodeAt(account any, blockNumber ...any) (code []byte) {
	var _blockNumber *big.Int
	if len(blockNumber) > 0 {
		_blockNumber = BigInt(blockNumber[0])
	}
	_account := Address(account)
	for {
		client := c.it.WaitNext()
		code, err := client.CodeAt(c.ctx, _account, _blockNumber)
		if err != nil {
			c.record(client.CodeAt, client, err)
			continue
		}
		return code
	}
}

// PendingCodeAt returns the contract code of the given account in the pending state.
func (c *Clientx) PendingCodeAt(account any) (code []byte) {
	_account := Address(account)
	for {
		client := c.it.WaitNext()
		code, err := client.PendingCodeAt(c.ctx, _account)
		if err != nil {
			c.record(client.PendingCodeAt, client, err)
			continue
		}
		return code
	}
}

// notFoundReturn returns the stop blockNumber for the notFound-error.
func (c *Clientx) notFoundReturn(notFoundBlocks ...uint64) uint64 {
	if len(notFoundBlocks) > 0 {
		return c.BlockNumber() + notFoundBlocks[0]
	}
	return c.BlockNumber() + c.notFoundBlocks
}

// BlockByHash returns the given full block.
//
// Note that loading full blocks requires two requests. Use HeaderByHash
// if you don't need all transactions or uncle headers.
func (c *Clientx) BlockByHash(hash any, notFoundBlocks ...uint64) (block *types.Block, err error) {
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()
	_hash := Hash(hash)
	_notFoundReturn := c.notFoundReturn(notFoundBlocks...)
	for {
		client := c.it.WaitNext()
		block, err = client.BlockByHash(c.ctx, _hash)
		if err != nil {
			c.record(client.BlockByHash, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if _notFoundReturn <= c.BlockNumber() {
					return nil, err
				}
			}
			<-queryTicker.C
			continue
		}
		return block, nil
	}
}

// BlockByNumber returns a block from the current canonical chain. If number is nil, the
// latest known block is returned.
//
// Note that loading full blocks requires two requests. Use HeaderByNumber
// if you don't need all transactions or uncle headers.
func (c *Clientx) BlockByNumber(blockNumber any, notFoundBlocks ...uint64) (block *types.Block, err error) {
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()
	_blockNumber := BigInt(blockNumber)
	_notFoundReturn := c.notFoundReturn(notFoundBlocks...)
	for {
		client := c.it.WaitNext()
		block, err = client.BlockByNumber(c.ctx, _blockNumber)
		if err != nil {
			c.record(client.BlockByNumber, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if _notFoundReturn <= c.BlockNumber() {
					return nil, err
				}
			}
			<-queryTicker.C
			continue
		}
		return block, nil
	}
}

// HeaderByHash returns the block header with the given hash.
func (c *Clientx) HeaderByHash(hash any, notFoundBlocks ...uint64) (header *types.Header, err error) {
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()
	_hash := Hash(hash)
	_notFoundReturn := c.notFoundReturn(notFoundBlocks...)
	for {
		client := c.it.WaitNext()
		header, err = client.HeaderByHash(c.ctx, _hash)
		if err != nil {
			c.record(client.HeaderByHash, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if _notFoundReturn <= c.BlockNumber() {
					return nil, err
				}
			}
			<-queryTicker.C
			continue
		}
		return header, nil
	}
}

// HeaderByNumber returns a block header from the current canonical chain. If number is
// nil, the latest known header is returned.
func (c *Clientx) HeaderByNumber(blockNumber any, notFoundBlocks ...uint64) (header *types.Header, err error) {
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()
	_blockNumber := BigInt(blockNumber)
	_notFoundReturn := c.notFoundReturn(notFoundBlocks...)
	for {
		client := c.it.WaitNext()
		header, err = client.HeaderByNumber(c.ctx, _blockNumber)
		if err != nil {
			c.record(client.HeaderByNumber, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if _notFoundReturn <= c.BlockNumber() {
					return nil, err
				}
			}
			<-queryTicker.C
			continue
		}
		return header, nil
	}
}

// TransactionByHash returns the transaction with the given hash.
func (c *Clientx) TransactionByHash(hash any, notFoundBlocks ...uint64) (tx *types.Transaction, isPending bool, err error) {
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()
	_hash := Hash(hash)
	_notFoundReturn := c.notFoundReturn(notFoundBlocks...)
	for {
		client := c.it.WaitNext()
		tx, isPending, err = client.TransactionByHash(c.ctx, _hash)
		if err != nil {
			c.record(client.TransactionByHash, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if _notFoundReturn <= c.BlockNumber() {
					return nil, isPending, err
				}
			}
			<-queryTicker.C
			continue
		}
		return tx, isPending, nil
	}
}

// TransactionSender returns the sender address of the given transaction. The transaction
// must be known to the remote node and included in the blockchain at the given block and
// index. The sender is the one derived by the protocol at the time of inclusion.
//
// There is a fast-path for transactions retrieved by TransactionByHash and
// TransactionInBlock. Getting their sender address can be done without an RPC interaction.
func (c *Clientx) TransactionSender(tx *types.Transaction, blockHash any, index uint, notFoundBlocks ...uint64) (sender common.Address, err error) {
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()
	_blockHash := Hash(blockHash)
	_notFoundReturn := c.notFoundReturn(notFoundBlocks...)
	for {
		client := c.it.WaitNext()
		sender, err = client.TransactionSender(c.ctx, tx, _blockHash, index)
		if err != nil {
			c.record(client.TransactionSender, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if _notFoundReturn <= c.BlockNumber() {
					return sender, err
				}
			}
			<-queryTicker.C
			continue
		}
		return sender, err
	}
}

// EstimateGas estimate tx gasUsed with maxTry.
func (c *Clientx) EstimateGas(msg ethereum.CallMsg, maxTry ...int) (gasLimit uint64, err error) {
	n := 1
	if len(maxTry) > 0 && maxTry[0] > n {
		n = maxTry[0]
	}
	for i := 0; i < n; i++ {
		client := c.it.WaitNext()
		gasLimit, err = client.EstimateGas(c.ctx, msg)
		if err != nil {
			c.record(client.EstimateGas, client, err)
			continue
		}
		break
	}
	return gasLimit, err
}

// SendTransaction send Transaction with maxTry.
func (c *Clientx) SendTransaction(tx *types.Transaction, maxTry ...int) (err error) {
	n := 1
	if len(maxTry) > 0 && maxTry[0] > n {
		n = maxTry[0]
	}
	for i := 0; i < n; i++ {
		client := c.it.WaitNext()
		err = client.SendTransaction(c.ctx, tx)
		if err != nil {
			c.record(client.SendTransaction, client, err)
			continue
		}
		break
	}
	return err
}

// TransactionCount returns the total number of transactions in the given block.
func (c *Clientx) TransactionCount(blockHash any, notFoundBlocks ...uint64) (count uint, err error) {
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()
	_blockHash := Hash(blockHash)
	_notFoundReturn := c.notFoundReturn(notFoundBlocks...)
	for {
		client := c.it.WaitNext()
		count, err = client.TransactionCount(c.ctx, _blockHash)
		if err != nil {
			c.record(client.TransactionCount, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if _notFoundReturn <= c.BlockNumber() {
					return count, err
				}
			}
			<-queryTicker.C
			continue
		}
		return count, nil
	}
}

// PendingTransactionCount returns the total number of transactions in the pending state.
func (c *Clientx) PendingTransactionCount() (count uint) {
	var err error
	for {
		client := c.it.WaitNext()
		count, err = client.PendingTransactionCount(c.ctx)
		if err != nil {
			c.record(client.PendingTransactionCount, client, err)
			continue
		}
		return
	}
}

// TransactionInBlock returns a single transaction at index in the given block.
func (c *Clientx) TransactionInBlock(blockHash any, index uint, notFoundBlocks ...uint64) (tx *types.Transaction, err error) {
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()
	_blockHash := Hash(blockHash)
	_notFoundReturn := c.notFoundReturn(notFoundBlocks...)
	for {
		client := c.it.WaitNext()
		tx, err = client.TransactionInBlock(c.ctx, _blockHash, index)
		if err != nil {
			c.record(client.TransactionInBlock, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if _notFoundReturn <= c.BlockNumber() {
					return tx, err
				}
			}
			<-queryTicker.C
			continue
		}
		return tx, err
	}
}

// TransactionReceipt returns the receipt of a transaction by transaction hash.
// Note that the receipt is not available for pending transactions.
func (c *Clientx) TransactionReceipt(txHash any, notFoundBlocks ...uint64) (receipt *types.Receipt, err error) {
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()
	_txHash := Hash(txHash)
	_notFoundReturn := c.notFoundReturn(notFoundBlocks...)
	for {
		client := c.it.WaitNext()
		receipt, err = client.TransactionReceipt(c.ctx, _txHash)
		if err != nil {
			c.record(client.TransactionReceipt, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if _notFoundReturn <= c.BlockNumber() {
					return receipt, err
				}
			}
			<-queryTicker.C
			continue
		}
		return receipt, err
	}
}

// WaitMined waits for tx to be mined on the blockchain.
// It stops waiting when the context is canceled.
// ethereum/go-ethereum@v1.11.6/accounts/abi/bind/util.go:32
func (c *Clientx) WaitMined(tx *types.Transaction, confirmBlocks uint64, notFoundBlocks ...uint64) (*types.Receipt, error) {
	queryTicker := time.NewTicker(c.miningInterval)
	defer queryTicker.Stop()
	txHash := tx.Hash()
	_confirmReturn := c.BlockNumber() + confirmBlocks
	_notFoundReturn := c.notFoundReturn(notFoundBlocks...)
	for {
		client := c.it.WaitNext()
		receipt, err := client.TransactionReceipt(c.ctx, txHash)
		if err != nil {
			c.record(c.WaitMined, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if _notFoundReturn <= c.BlockNumber() {
					return nil, err
				}
			}
		} else {
			_notFoundReturn = c.notFoundReturn(notFoundBlocks...)
			if _confirmReturn <= c.BlockNumber() {
				return receipt, nil
			}
		}
		<-queryTicker.C
	}
}

// WaitDeployed waits for a contract deployment transaction and returns the on-chain
// contract address when it is mined. It stops waiting when ctx is canceled.
func (c *Clientx) WaitDeployed(tx *types.Transaction, confirmBlocks uint64, notFoundBlocks ...uint64) (common.Address, error) {
	if tx.To() != nil {
		return common.Address{}, errors.New("tx is not contract creation")
	}
	receipt, err := c.WaitMined(tx, confirmBlocks, notFoundBlocks...)
	if err != nil {
		return common.Address{}, err
	}
	if receipt.ContractAddress == (common.Address{}) {
		return common.Address{}, errors.New("zero address")
	}
	// Check that code has indeed been deployed at the address.
	// This matters on pre-Homestead chains: OOG in the constructor
	// could leave an empty account behind.
	code := c.CodeAt(receipt.ContractAddress, nil)
	if err == nil && len(code) == 0 {
		err = errors.New("no contract code after deployment")
	}
	return receipt.ContractAddress, err
}

func (c *Clientx) Shuffle() {
	c.it.Shuffle()
}

type RawLogFilter struct {
	client    *Clientx
	addresses []common.Address
	topics    [][]common.Hash
	config    EventConfig
	txHashSet mapset.Set[string]
	mu        sync.Mutex
}

type EventConfig struct {
	IntervalBlocks, OverrideBlocks, DelayBlocks uint64
}

func (e *EventConfig) panicIfNotValid() {
	if e.IntervalBlocks == 0 {
		panic(errors.New("EventConfig::require IntervalBlocks > 0, eg: 800"))
	}
	if e.IntervalBlocks+e.OverrideBlocks > 2000 {
		panic(errors.New("EventConfig::require IntervalBlocks + OverrideBlocks <= 2000, eg: 800 + 800"))
	}
	if e.DelayBlocks == 0 {
		log.Printf("[WARN] EventConfig::If you are tracking the latest logs, DelayBlocks==0 is risky, recommended >= 3. see: https://github.com/ethereum/go-ethereum/blob/master/core/types/log.go#L53")
	}
}

// NewRawLogFilter returns the RawLogFilter
// EventConfig require IntervalBlocks + OverrideBlocks <= 2000, eg: 800,800
func (c *Clientx) NewRawLogFilter(topics [][]common.Hash, addresses []common.Address, config EventConfig) *RawLogFilter {
	config.panicIfNotValid()
	return &RawLogFilter{
		client:    c,
		addresses: addresses,
		topics:    topics,
		config:    config,
		txHashSet: mapset.NewThreadUnsafeSet[string](),
		mu:        sync.Mutex{},
	}
}

func (r *RawLogFilter) Run(from, to uint64) (logs []types.Log, addressTopicLogsMap map[common.Address]map[common.Hash][]*types.Log) {
	fc := func(_from, _to uint64) {
		// Attention!!!Repeat scanning _to prevent missing blocks
		var query = ethereum.FilterQuery{
			Addresses: r.addresses,
			Topics:    r.topics,
			FromBlock: new(big.Int).SetUint64(_from),
			ToBlock:   new(big.Int).SetUint64(_to),
		}
		nLogs := r.client.FilterLogs(query)
		log.Printf("Run: %v(%v)-%v Success!\n", _from, query.FromBlock, _to)
		if len(nLogs) > 0 {
			r.mu.Lock()
			var hashID string
			for _, nLog := range nLogs {
				hashID = fmt.Sprintf("%v%v", nLog.TxHash, nLog.Index)
				if !r.txHashSet.Contains(hashID) {
					r.txHashSet.Add(hashID)
					logs = append(logs, nLog)
				}
			}
			r.mu.Unlock()
		}
	}
	segmentCallback(from, to, r.config, fc)
	addressTopicLogsMap = make(map[common.Address]map[common.Hash][]*types.Log)
	for _, nLog := range logs {
		if addressTopicLogsMap[nLog.Address] == nil {
			addressTopicLogsMap[nLog.Address] = make(map[common.Hash][]*types.Log)
		}
		addressTopicLogsMap[nLog.Address][nLog.Topics[0]] = append(addressTopicLogsMap[nLog.Address][nLog.Topics[0]], &nLog)
	}

	return logs, addressTopicLogsMap
}

func segmentCallback(from, to uint64, config EventConfig, callback func(from, to uint64)) (newStart uint64) {
	if from+config.DelayBlocks <= to {
		// count
		to -= config.DelayBlocks
		count := (to - from) / config.IntervalBlocks
		wg := new(sync.WaitGroup)
		wg.Add(int(count))
		// loop
		arrest := from + config.OverrideBlocks
		for i := uint64(0); i < count; i++ {
			_from := from + i*config.IntervalBlocks
			_to := _from + config.IntervalBlocks - 1
			if _from >= arrest {
				_from -= config.OverrideBlocks
			}
			go func() {
				callback(_from, _to)
				wg.Done()
			}()
		}
		if last := from + count*config.IntervalBlocks; last <= to {
			callback(last, to)
		}
		wg.Wait()
		return to + 1
	}
	return from
}

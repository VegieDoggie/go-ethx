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
	"golang.org/x/time/rate"
	"log"
	"math/big"
	"strings"
	"sync"
	"time"
)

type ClientxConfig struct {
	Event            EventConfig // MustContract and RawLogger
	GasLimit         uint64      // Clientx
	NotFundBlocks    uint64      // Clientx
	GasTipAdditional *big.Int    // Clientx
	MaxMustErrNumR   int         // MustContract
	MaxMustErrNumW   int         // MustContract
}

// Clientx defines typed wrappers for the Ethereum RPC API of a set of the Ethereum Clients.
type Clientx struct {
	it             *Iterator[*ethclient.Client]
	ctx            context.Context
	rpcMap         map[*ethclient.Client]string
	rpcErrCountMap map[*ethclient.Client]uint
	notFoundBlocks uint64
	chainId        *big.Int
	latestHeader   *types.Header
	startedAt      time.Time
	config         *ClientxConfig
	non1559Gas     bool
	mu             sync.Mutex
}

func NewClientxConfig() *ClientxConfig {
	return &ClientxConfig{
		Event: EventConfig{
			IntervalBlocks: 200,
			OverrideBlocks: 1300,
			DelayBlocks:    3,
		},
		GasLimit:         8000000,
		NotFundBlocks:    9,
		MaxMustErrNumR:   999,
		MaxMustErrNumW:   3,
		GasTipAdditional: big.NewInt(0),
	}
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
	return NewClientx(rpcList, weights, NewClientxConfig())
}

// NewClientx connects clients to the given URLs, to provide a reliable Ethereum RPC API call, includes
// a timer to regularly update block height(AutoBlockNumber).
// If weight <= 1, the weight is always 1.
// Note: If len(weightList) == 0, then default weight = 1 will be active.
func NewClientx(rpcList []string, weights []int, config *ClientxConfig, limiter ...*rate.Limiter) *Clientx {
	rpcList = mapset.NewThreadUnsafeSet[string](rpcList...).ToSlice()
	iterator, rpcMap, chainId := buildIterator(rpcList, weights, limiter...)
	notFundBlocks := uint64(len(rpcList))
	if notFundBlocks < config.NotFundBlocks {
		notFundBlocks = config.NotFundBlocks
	}
	c := &Clientx{
		ctx:            context.Background(),
		it:             iterator,
		rpcMap:         rpcMap,
		chainId:        chainId,
		rpcErrCountMap: make(map[*ethclient.Client]uint),
		notFoundBlocks: notFundBlocks,
		latestHeader:   &types.Header{Number: BigInt(0)},
		startedAt:      time.Now(),
		config:         config,
	}
	c.init()
	defer time.Sleep(2 * time.Second)
	return c
}

// buildIterator creates a clientIterator with wights.
func buildIterator(rpcList []string, weightList []int, limiter ...*rate.Limiter) (clientIterator *Iterator[*ethclient.Client], rpcMap map[*ethclient.Client]string, latestChainId *big.Int) {
	if len(rpcList) != len(weightList) {
		tmp := weightList
		weightList = make([]int, len(rpcList))
		copy(weightList, tmp)
	}

	clientIterator = new(Iterator[*ethclient.Client])
	rpcMap = make(map[*ethclient.Client]string)

	var reliableClients []*ethclient.Client
	var mu sync.Mutex
	update := func(rpc string, client *ethclient.Client, weight int) {
		mu.Lock()
		defer mu.Unlock()
		rpcMap[client] = rpc
		reliableClients = append(reliableClients, client)
		for k := 1; k < weight; k++ {
			reliableClients = append(reliableClients, client)
		}
		*clientIterator = *NewIterator[*ethclient.Client](reliableClients, limiter...).Shuffle()
	}

	for i, _rpc := range rpcList {
		client, chainId, err := checkChainid(_rpc, 3)
		if err != nil {
			log.Printf("[WARN] buildIterator::%v\n", err)
			continue
		}
		latestChainId = chainId
		update(_rpc, client, weightList[i])
		go func() {
			for j := i + 1; j < len(rpcList); j++ {
				_rpc = rpcList[j]
				client, chainId, err := checkChainid(_rpc, 3)
				if err != nil {
					log.Printf("[WARN] buildIterator::%v\n", err)
					continue
				}
				if latestChainId.Cmp(chainId) != 0 {
					panic(fmt.Sprintf("[ERROR] [ABORT] buildIterator::previous chainID is %v,but rpc(%v) is chainId(%v)!\n", latestChainId, _rpc, chainId))
				}
				update(_rpc, client, weightList[j])
			}
		}()
		break
	}
	if latestChainId == nil {
		panic(fmt.Errorf("[ERROR] buildIterator::Unreliable rpc List: %v\n", rpcList))
	}
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

// log.Println(fmt.Sprintf("[WARN] UpdateRPCs::required chainID is %v,but rpc(%v) chainId is %v\n", c.chainId, rpc, chainId))
func checkChainid(rpc string, maxErr ...int) (*ethclient.Client, *big.Int, error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, nil, err
	}
	num := 1
	if len(maxErr) > 0 {
		num = maxErr[0]
	}
	for i := 0; i < num; i++ {
		chainId, err := client.ChainID(context.TODO())
		if err != nil {
			continue
		}
		return client, chainId, nil
	}
	return nil, nil, fmt.Errorf("unreliable RPC: %v\n", rpc)
}

func (c *Clientx) UpdateRPCs(newRPCs []string) {
	if time.Since(c.startedAt) < 60*time.Second {
		log.Printf("[WARN] UpdateRPCs::started less than 60s\n")
		return
	}
	if len(newRPCs) == 0 {
		return
	}
	oldRPCs := c.GetRPCs()
	newSet := mapset.NewThreadUnsafeSet[string](newRPCs...)
	oldSet := mapset.NewThreadUnsafeSet[string](oldRPCs...)
	updated := false
	// rpc in newRPCs but not in oldRPCs: ADD
	for _, rpc := range newRPCs {
		if !oldSet.Contains(rpc) {
			client, chainId, err := checkChainid(rpc, 3)
			if err != nil {
				log.Printf("[WARN] UpdateRPCs::Unreliable rpc: %v\n", rpc)
				continue
			}
			if c.chainId.Cmp(chainId) != 0 {
				log.Printf("[WARN] UpdateRPCs::required chainID is %v,but rpc(%v) chainId is %v\n", c.chainId, rpc, chainId)
				continue
			}
			updated = true
			c.rpcMap[client] = rpc
			c.it.Add(client)
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
			if client != nil && c.it.Len() > 1 {
				updated = true
				delete(c.rpcMap, client)
				c.it.Remove(client)
			}
		}
	}
	if updated {
		log.Printf("[SUCCESS] UpdateRPCs::from=%v, to=%v\n", newRPCs, c.GetRPCs())
	}
}

func (c *Clientx) errorCallback(f any, client *ethclient.Client, err error) {
	//c.mu.Lock()
	//defer c.mu.Unlock()
	//c.rpcErrCountMap[client]++
	//prefix := err.Error()
	//if len(prefix) > 2 {
	//	prefix = prefix[:3]
	//}
	//switch prefix {
	//case "429", "521":
	//	if c.rpcErrCountMap[client]%10 == 0 {
	//		log.Printf("[%v] %v [WARN] func=%v, rpc=%v #%v, err=%v\r\n", c.BlockNumber(), time.Now().Format(time.DateTime), getFuncName(f), c.rpcMap[client], c.rpcErrCountMap[client], "Too Many Requests!")
	//	}
	//default:
	//	log.Printf("[%v] %v [WARN] func=%v, rpc=%v #%v, err=%v\r\n", c.BlockNumber(), time.Now().Format(time.DateTime), getFuncName(f), c.rpcMap[client], c.rpcErrCountMap[client], err)
	//}
}

// TransactOpts create *bind.TransactOpts, and panic if privateKey err
// privateKeyLike eg: 0xf1...3, f1...3, []byte, *ecdsa.PrivateKey...
func (c *Clientx) TransactOpts(privateKeyLike any) *bind.TransactOpts {
	_privateKey := PrivateKey(privateKeyLike)
	opts, err := bind.NewKeyedTransactorWithChainID(_privateKey, c.chainId)
	if err != nil {
		panic(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		if c.latestHeader.BaseFee != nil {
			gasTipCap := c.SuggestGasTipCap()
			opts.GasFeeCap = Add(Mul(c.latestHeader.BaseFee, 2), gasTipCap)
			opts.GasTipCap = gasTipCap
		} else {
			opts.GasPrice = c.SuggestGasPrice()
		}
		opts.GasLimit = c.config.GasLimit
		wg.Done()
	}()
	go func() {
		opts.Nonce = BigInt(c.PendingNonceAt(opts.From))
		wg.Done()
	}()
	wg.Wait()
	return opts
}

type TransferOption struct {
	Data       []byte             // option
	AccessList types.AccessList   // option
	Opts       *bind.TransactOpts // option
}

func (c *Clientx) unsignedTx(privateKeyLike, to, amount any, options ...TransferOption) (tx *types.Transaction, opts *bind.TransactOpts) {
	var data []byte
	var accessList types.AccessList
	if len(options) > 0 {
		data = options[0].Data
		accessList = options[0].AccessList
		opts = options[0].Opts
	}
	if opts == nil {
		opts = c.TransactOpts(privateKeyLike)
	}
	switch {
	case opts.GasFeeCap != nil:
		tx = types.NewTx(&types.DynamicFeeTx{
			To:         AddressPtr(to),
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
			To:         AddressPtr(to),
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
			To:       AddressPtr(to),
			Nonce:    opts.Nonce.Uint64(),
			Gas:      opts.GasLimit,
			GasPrice: opts.GasPrice,
			Value:    BigInt(amount),
			Data:     data,
		})
	}
	return tx, opts
}

func (c *Clientx) TransferETH(privateKeyLike, to any, value ...any) (tx *types.Transaction, err error) {
	opts := c.TransactOpts(privateKeyLike)
	opts.GasLimit = 21001
	var gasCost *big.Int
	switch {
	case opts.GasFeeCap != nil:
		gasCost = Mul(opts.GasFeeCap, opts.GasLimit)
	default:
		gasCost = Mul(opts.GasPrice, opts.GasLimit)
	}
	var amount *big.Int
	if len(value) > 0 {
		amount = BigInt(value[0])
	} else {
		amount = c.BalanceAt(Address(privateKeyLike, true))
	}
	if Lte(amount, gasCost) {
		return nil, fmt.Errorf("Insufficient funds\n")
	}
	if len(value) > 0 {
		return c.Transfer(privateKeyLike, to, amount, TransferOption{Opts: opts})
	}
	return c.Transfer(privateKeyLike, to, Sub(amount, gasCost), TransferOption{Opts: opts})
}

// Transfer build transaction and send
// TransferOption is optional.
// see more: github.com/ethereum/go-ethereum/internal/ethapi/transaction_args.go:284
func (c *Clientx) Transfer(privateKeyLike, to, amount any, options ...TransferOption) (tx *types.Transaction, err error) {
	tx, opts := c.unsignedTx(privateKeyLike, to, amount, options...)
	if tx, err = c.send(tx, opts); err != nil {
		return nil, err
	}
	return tx, nil
}

func (c *Clientx) send(tx *types.Transaction, opts *bind.TransactOpts) (signedTx *types.Transaction, err error) {
	if tx, err = opts.Signer(opts.From, tx); err != nil {
		return nil, err
	}
	gasLimit, err := c.EstimateGas(CallMsg(opts.From, tx), len(c.rpcMap))
	if err != nil {
		return nil, err
	}
	if opts.GasLimit < gasLimit || gasLimit == 0 {
		return nil, errors.New(fmt.Sprintf("[ERROR] Transfer::Gas required %v, but %v.\n", gasLimit, opts.GasLimit))
	}
	if err = c.SendTransaction(tx, len(c.rpcMap)); err != nil {
		return nil, err
	}
	return tx, nil
}

func (c *Clientx) init() {
	go func() {
		doneTicker := time.NewTicker(time.Second)
		defer doneTicker.Stop()

		failTicker := time.NewTicker(time.Second / 2)
		defer failTicker.Stop()

		for {
			header, err := c.NextClient().HeaderByNumber(c.ctx, nil)
			if err == nil {
				if header.Number.Cmp(c.latestHeader.Number) > 0 {
					c.latestHeader = header
				}

				<-doneTicker.C
				continue
			}

			<-failTicker.C
		}
	}()

	for c.latestHeader.Number.Uint64() == 0 {
		time.Sleep(100 * time.Millisecond)
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
			c.errorCallback(client.NetworkID, client, err)
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
			c.errorCallback(client.BalanceAt, client, err)
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
			c.errorCallback(client.PendingBalanceAt, client, err)
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
			c.errorCallback(client.NonceAt, client, err)
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
			c.errorCallback(client.PendingNonceAt, client, err)
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
			c.errorCallback(client.FilterLogs, client, err)
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
			c.errorCallback(client.SuggestGasPrice, client, err)
			continue
		}
		return Add(gasPrice, c.config.GasTipAdditional)
	}
}

// SuggestGasTipCap retrieves the currently suggested gas tip cap after 1559 to
// allow a timely execution of a transaction.
func (c *Clientx) SuggestGasTipCap() (gasTipCap *big.Int) {
	for {
		client := c.it.WaitNext()
		if c.non1559Gas {
			return c.SuggestGasPrice()
		}
		gasTipCap, err := client.SuggestGasTipCap(c.ctx)
		if err != nil {
			c.errorCallback(client.SuggestGasTipCap, client, err)
			if strings.LastIndex(err.Error(), "not found") != -1 {
				c.non1559Gas = true
			}
			continue
		}
		return Add(gasTipCap, c.config.GasTipAdditional)
	}
}

// FeeHistory retrieves the fee market history.
func (c *Clientx) FeeHistory(blockCount uint64, lastBlock any, rewardPercentiles []float64) (feeHistory *ethereum.FeeHistory) {
	for {
		client := c.it.WaitNext()
		feeHistory, err := client.FeeHistory(c.ctx, blockCount, BigInt(lastBlock), rewardPercentiles)
		if err != nil {
			c.errorCallback(client.FeeHistory, client, err)
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
			c.errorCallback(client.StorageAt, client, err)
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
			c.errorCallback(client.PendingStorageAt, client, err)
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
			c.errorCallback(client.CodeAt, client, err)
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
			c.errorCallback(client.PendingCodeAt, client, err)
			continue
		}
		return code
	}
}

// notFoundReturn returns the stop blockNumber for the notFound-error.
func (c *Clientx) notFoundReturn(notFoundBlocks []uint64) uint64 {
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
	notFoundEnd := c.notFoundReturn(notFoundBlocks)
	for {
		client := c.it.WaitNext()
		block, err = client.BlockByHash(c.ctx, _hash)
		if err != nil {
			c.errorCallback(client.BlockByHash, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if notFoundEnd <= c.BlockNumber() {
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
	notFoundEnd := c.notFoundReturn(notFoundBlocks)
	for {
		client := c.it.WaitNext()
		block, err = client.BlockByNumber(c.ctx, _blockNumber)
		if err != nil {
			c.errorCallback(client.BlockByNumber, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if notFoundEnd <= c.BlockNumber() {
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
	notFoundEnd := c.notFoundReturn(notFoundBlocks)
	for {
		client := c.it.WaitNext()
		header, err = client.HeaderByHash(c.ctx, _hash)
		if err != nil {
			c.errorCallback(client.HeaderByHash, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if notFoundEnd <= c.BlockNumber() {
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
	notFoundEnd := c.notFoundReturn(notFoundBlocks)
	for {
		client := c.it.WaitNext()
		header, err = client.HeaderByNumber(c.ctx, _blockNumber)
		if err != nil {
			c.errorCallback(client.HeaderByNumber, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if notFoundEnd <= c.BlockNumber() {
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
	notFoundEnd := c.notFoundReturn(notFoundBlocks)
	for {
		client := c.it.WaitNext()
		tx, isPending, err = client.TransactionByHash(c.ctx, _hash)
		if err != nil {
			c.errorCallback(client.TransactionByHash, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if notFoundEnd <= c.BlockNumber() {
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
	failTicker := time.NewTicker(time.Second)
	defer failTicker.Stop()

	_blockHash := Hash(blockHash)
	notFoundEnd := c.notFoundReturn(notFoundBlocks)

	for {
		client := c.it.WaitNext()
		sender, err = client.TransactionSender(c.ctx, tx, _blockHash, index)
		if err != nil {
			c.errorCallback(client.TransactionSender, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if notFoundEnd <= c.BlockNumber() {
					return sender, err
				}
			}
			<-failTicker.C
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
			c.errorCallback(client.EstimateGas, client, err)
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
			c.errorCallback(client.SendTransaction, client, err)
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
	notFoundEnd := c.notFoundReturn(notFoundBlocks)
	for {
		client := c.it.WaitNext()
		count, err = client.TransactionCount(c.ctx, _blockHash)
		if err != nil {
			c.errorCallback(client.TransactionCount, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if notFoundEnd <= c.BlockNumber() {
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
			c.errorCallback(client.PendingTransactionCount, client, err)
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
	notFoundEnd := c.notFoundReturn(notFoundBlocks)
	for {
		client := c.it.WaitNext()
		tx, err = client.TransactionInBlock(c.ctx, _blockHash, index)
		if err != nil {
			c.errorCallback(client.TransactionInBlock, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if notFoundEnd <= c.BlockNumber() {
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
	notFoundEnd := c.notFoundReturn(notFoundBlocks)
	for {
		client := c.it.WaitNext()
		receipt, err = client.TransactionReceipt(c.ctx, _txHash)
		if err != nil {
			c.errorCallback(client.TransactionReceipt, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if notFoundEnd <= c.BlockNumber() {
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
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()

	confirmEnd := c.BlockNumber() + confirmBlocks
	notFoundEnd := c.notFoundReturn(notFoundBlocks)
	for txHash := tx.Hash(); ; {
		client := c.it.WaitNext()
		receipt, err := client.TransactionReceipt(c.ctx, txHash)
		if err != nil {
			c.errorCallback(c.WaitMined, client, fmt.Errorf("%v (hash=%v)", err, txHash))
			if errors.Is(err, ethereum.NotFound) {
				if notFoundEnd <= c.BlockNumber() {
					return nil, err
				}
			}
		} else {
			notFoundEnd = c.notFoundReturn(notFoundBlocks)
			if confirmEnd <= c.BlockNumber() {
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

func (c *Clientx) Shuffle() *Clientx {
	c.it.Shuffle()
	return c
}

func (c *Clientx) segmentCallback(from, to uint64, config *EventConfig, callback func(from, to uint64)) (newStart uint64) {
	if from+config.DelayBlocks <= to {
		to -= config.DelayBlocks
		count := (to - from) / config.IntervalBlocks

		wg := new(sync.WaitGroup)
		arrestFrom := from + config.OverrideBlocks
		for i := uint64(0); i < count; {
			wg.Add(1)
			segFrom := from + i*config.IntervalBlocks
			segTo := segFrom + config.IntervalBlocks - 1
			if segFrom >= arrestFrom {
				segFrom -= config.OverrideBlocks
			}
			go func() {
				callback(segFrom, segTo)
				wg.Done()
			}()
			i++
			if (i+1)%128 == 0 {
				wg.Wait()
			}
		}

		if segFrom := from + count*config.IntervalBlocks; segFrom <= to {
			if segFrom >= arrestFrom {
				segFrom -= config.OverrideBlocks
			}
			if c.it.Len() > 3 {
				for i := 0; i < 3; i++ {
					wg.Add(1)
					go func() {
						callback(segFrom, to)
						wg.Done()
					}()
				}
			} else {
				callback(segFrom, to)
				time.Sleep(1 * time.Second)
				callback(segFrom, to)
				time.Sleep(1 * time.Second)
				callback(segFrom, to)
			}
		}

		wg.Wait()
		return to + 1
	}
	return from
}

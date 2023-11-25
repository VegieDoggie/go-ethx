package ethx

import (
	"context"
	"errors"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/time/rate"
	"log"
	"math"
	"math/big"
	"sync"
	"time"
)

// Clientx defines typed wrappers for the Ethereum RPC API of a set of the Ethereum Clients.
type Clientx struct {
	*Iterator[*ethclient.Client]
	ctx            context.Context
	RpcMap         map[*ethclient.Client]string
	rpcErrCountMap map[*ethclient.Client]uint
	NotFoundBlocks uint64
	chainId        *big.Int
	latestHeader   *types.Header
}

// NewClientx connects clients to the given URLs, to provide a reliable Ethereum RPC API call, includes
// a timer to regularly update block height(AutoBlockNumber).
// If weight <= 1, the weight is always 1.
// Note: If len(weightList) == 0, then default weight = 1 will be active.
func NewClientx(rpcList []string, weights []int, defaultNotFoundBlocks uint64, limiter ...*rate.Limiter) *Clientx {
	iterator, rpcMap, chainId := newClientIteratorWithWeight(rpcList, weights, limiter...)
	c := &Clientx{
		ctx:            context.Background(),
		Iterator:       iterator,
		RpcMap:         rpcMap,
		chainId:        chainId,
		rpcErrCountMap: make(map[*ethclient.Client]uint),
		NotFoundBlocks: defaultNotFoundBlocks,
		latestHeader:   &types.Header{Number: BigInt(0)},
	}
	go func() {
		queryTicker := time.NewTicker(time.Second)
		defer queryTicker.Stop()
		for {
			header, err := c.HeaderByNumber(nil, 0)
			if err == nil && header.Number.Cmp(c.latestHeader.Number) > 1 {
				c.latestHeader = header
			}
			<-queryTicker.C
		}
	}()
	queryTicker := time.NewTicker(100 * time.Millisecond)
	defer queryTicker.Stop()
	for c.latestHeader.Number.Uint64() == 0 {
		<-queryTicker.C
	}
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
				if latestChainId != nil {
					latestChainId = chainId
				} else if latestChainId.Cmp(chainId) != 0 {
					panic(errors.New(fmt.Sprintf("[ERROR] rpc(%v) chainID is %v,but rpc(%v) chainId is %v\n", rpcList[i-1], latestChainId, rpc, chainId)))
				}
				for j := 1; j < weightList[i]; j++ {
					reliableClients = append(reliableClients, client)
				}
				reliableClients = append(reliableClients, client)
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
	return
}

func (c *Clientx) logWarn(f any, client *ethclient.Client, err error) {
	c.rpcErrCountMap[client]++
	log.Printf("%v [WARN] func=%v, rpc=%v #%v, err=%v\r\n", time.Now().Format("2006-01-02 15:04:05"), getFuncName(f), c.RpcMap[client], c.rpcErrCountMap[client], err)
}

func (c *Clientx) NewMust(constructor any, addressLike any, maxErrNum ...int) func(f any, args ...any) any {
	n := math.MaxInt
	if len(maxErrNum) > 0 {
		n = maxErrNum[0]
	}
	address := Address(addressLike)
	return func(f any, args ...any) any {
		for i := 0; i < n; i++ {
			client := c.WaitNext()
			ret := callFunc(constructor, address, client)
			if err, ok := ret[len(ret)-1].(error); ok {
				panic(err)
			}
			ret = callStructFunc(ret[0], f, args...)
			if err, ok := ret[len(ret)-1].(error); ok {
				c.logWarn(f, client, err)
				continue
			}
			return ret[:len(ret)-1][0]
		}
		panic(errors.New(fmt.Sprintf("%v [%v]: exceed maxErrNum(%v)", address, getFuncName(f), n)))
	}
}

// Close all clients connections.
func (c *Clientx) Close() {
	clients := c.Iterator.All()
	for i := range clients {
		clients[i].Close()
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
		client := c.WaitNext()
		networkID, err = client.NetworkID(c.ctx)
		if err != nil {
			c.logWarn(client.NetworkID, client, err)
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
		client := c.WaitNext()
		balance, err := client.BalanceAt(c.ctx, _account, _blockNumber)
		if err != nil {
			c.logWarn(client.BalanceAt, client, err)
			continue
		}
		return balance
	}
}

// PendingBalanceAt returns the wei balance of the given account in the pending state.
func (c *Clientx) PendingBalanceAt(account any) (balance *big.Int) {
	_account := Address(account)
	for {
		client := c.WaitNext()
		balance, err := client.PendingBalanceAt(c.ctx, _account)
		if err != nil {
			c.logWarn(client.PendingBalanceAt, client, err)
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
		client := c.WaitNext()
		nonce, err := client.NonceAt(c.ctx, _account, _blockNumber)
		if err != nil {
			c.logWarn(client.NonceAt, client, err)
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
		client := c.WaitNext()
		nonce, err := client.PendingNonceAt(c.ctx, _account)
		if err != nil {
			c.logWarn(client.PendingNonceAt, client, err)
			continue
		}
		return nonce
	}
}

// FilterLogs executes a filter query.
func (c *Clientx) FilterLogs(q ethereum.FilterQuery) (logs []types.Log) {
	for {
		client := c.WaitNext()
		logs, err := client.FilterLogs(c.ctx, q)
		if err != nil {
			c.logWarn(client.FilterLogs, client, err)
			continue
		}
		return logs
	}
}

// SuggestGasPrice retrieves the currently suggested gas price to allow a timely
// execution of a transaction.
func (c *Clientx) SuggestGasPrice() (gasPrice *big.Int) {
	for {
		client := c.WaitNext()
		gasPrice, err := client.SuggestGasPrice(c.ctx)
		if err != nil {
			c.logWarn(client.SuggestGasPrice, client, err)
			continue
		}
		return gasPrice
	}
}

// SuggestGasTipCap retrieves the currently suggested gas tip cap after 1559 to
// allow a timely execution of a transaction.
func (c *Clientx) SuggestGasTipCap() (gasTipCap *big.Int) {
	for {
		client := c.WaitNext()
		gasTipCap, err := client.SuggestGasTipCap(c.ctx)
		if err != nil {
			c.logWarn(client.SuggestGasTipCap, client, err)
			continue
		}
		return gasTipCap
	}
}

// FeeHistory retrieves the fee market history.
func (c *Clientx) FeeHistory(blockCount uint64, lastBlock any, rewardPercentiles []float64) (feeHistory *ethereum.FeeHistory) {
	for {
		client := c.WaitNext()
		feeHistory, err := client.FeeHistory(c.ctx, blockCount, BigInt(lastBlock), rewardPercentiles)
		if err != nil {
			c.logWarn(client.FeeHistory, client, err)
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
		client := c.WaitNext()
		storage, err := client.StorageAt(c.ctx, _account, _keyHash, _blockNumber)
		if err != nil {
			c.logWarn(client.StorageAt, client, err)
			continue
		}
		return storage
	}
}

// PendingStorageAt returns the value of key in the contract storage of the given account in the pending state.
func (c *Clientx) PendingStorageAt(account, keyHash any) (storage []byte) {
	_account, _keyHash := Address(account), Hash(keyHash)
	for {
		client := c.WaitNext()
		storage, err := client.PendingStorageAt(c.ctx, _account, _keyHash)
		if err != nil {
			c.logWarn(client.PendingStorageAt, client, err)
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
		client := c.WaitNext()
		code, err := client.CodeAt(c.ctx, _account, _blockNumber)
		if err != nil {
			c.logWarn(client.CodeAt, client, err)
			continue
		}
		return code
	}
}

// PendingCodeAt returns the contract code of the given account in the pending state.
func (c *Clientx) PendingCodeAt(account any) (code []byte) {
	_account := Address(account)
	for {
		client := c.WaitNext()
		code, err := client.PendingCodeAt(c.ctx, _account)
		if err != nil {
			c.logWarn(client.PendingCodeAt, client, err)
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
	return c.BlockNumber() + c.NotFoundBlocks
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
		client := c.WaitNext()
		block, err = client.BlockByHash(c.ctx, _hash)
		if err != nil {
			c.logWarn(client.BlockByHash, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if _notFoundReturn >= c.BlockNumber() {
					return nil, err
				}
				<-queryTicker.C
			}
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
		client := c.WaitNext()
		block, err = client.BlockByNumber(c.ctx, _blockNumber)
		if err != nil {
			c.logWarn(client.BlockByNumber, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if _notFoundReturn >= c.BlockNumber() {
					return nil, err
				}
				<-queryTicker.C
			}
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
		client := c.WaitNext()
		header, err = client.HeaderByHash(c.ctx, _hash)
		if err != nil {
			c.logWarn(client.HeaderByHash, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if _notFoundReturn >= c.BlockNumber() {
					return nil, err
				}
				<-queryTicker.C
			}
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
		client := c.WaitNext()
		header, err = client.HeaderByNumber(c.ctx, _blockNumber)
		if err != nil {
			c.logWarn(client.HeaderByNumber, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if _notFoundReturn >= c.BlockNumber() {
					return nil, err
				}
				<-queryTicker.C
			}
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
		client := c.WaitNext()
		tx, isPending, err = client.TransactionByHash(c.ctx, _hash)
		if err != nil {
			c.logWarn(client.TransactionByHash, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if _notFoundReturn >= c.BlockNumber() {
					return nil, isPending, err
				}
				<-queryTicker.C
			}
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
		client := c.WaitNext()
		sender, err = client.TransactionSender(c.ctx, tx, _blockHash, index)
		if err != nil {
			c.logWarn(client.TransactionSender, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if _notFoundReturn >= c.BlockNumber() {
					return sender, err
				}
				<-queryTicker.C
			}
			continue
		}
		return sender, err
	}
}

// TransactionCount returns the total number of transactions in the given block.
func (c *Clientx) TransactionCount(blockHash any, notFoundBlocks ...uint64) (count uint, err error) {
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()
	_blockHash := Hash(blockHash)
	_notFoundReturn := c.notFoundReturn(notFoundBlocks...)
	for {
		client := c.WaitNext()
		count, err = client.TransactionCount(c.ctx, _blockHash)
		if err != nil {
			c.logWarn(client.TransactionCount, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if _notFoundReturn >= c.BlockNumber() {
					return count, err
				}
				<-queryTicker.C
			}
			continue
		}
		return count, nil
	}
}

// PendingTransactionCount returns the total number of transactions in the pending state.
func (c *Clientx) PendingTransactionCount() (count uint) {
	var err error
	for {
		client := c.WaitNext()
		count, err = client.PendingTransactionCount(c.ctx)
		if err != nil {
			c.logWarn(client.PendingTransactionCount, client, err)
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
		client := c.WaitNext()
		tx, err = client.TransactionInBlock(c.ctx, _blockHash, index)
		if err != nil {
			c.logWarn(client.TransactionInBlock, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if _notFoundReturn >= c.BlockNumber() {
					return tx, err
				}
				<-queryTicker.C
			}
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
		client := c.WaitNext()
		receipt, err = client.TransactionReceipt(c.ctx, _txHash)
		if err != nil {
			c.logWarn(client.TransactionReceipt, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if _notFoundReturn >= c.BlockNumber() {
					return receipt, err
				}
				<-queryTicker.C
			}
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
	txHash := tx.Hash()
	_confirmReturn := c.BlockNumber() + confirmBlocks
	_notFoundReturn := c.notFoundReturn(notFoundBlocks...)
	for {
		client := c.WaitNext()
		receipt, err := client.TransactionReceipt(c.ctx, txHash)
		if err != nil {
			c.logWarn(c.WaitMined, client, err)
			if errors.Is(err, ethereum.NotFound) {
				if _notFoundReturn >= c.BlockNumber() {
					return receipt, err
				}
			}
		} else {
			if _confirmReturn >= c.BlockNumber() {
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

type Scanner struct {
	*Clientx
	Addresses      []common.Address
	Topics         [][]common.Hash
	OverrideBlocks uint64
	IntervalBlocks uint64
	DelayBlocks    uint64
	TxHashSet      mapset.Set[string]
	Mu             sync.Mutex
}

type AddressTopicLogsMap = map[common.Address]map[common.Hash][]*types.Log

// NewScanner returns the next Ethereum Client.
func (c *Clientx) NewScanner(topics [][]common.Hash, addresses []common.Address, intervalBlocks, overrideBlocks, delayBlocks uint64) *Scanner {
	return &Scanner{
		Clientx:        c,
		Addresses:      addresses,
		Topics:         topics,
		IntervalBlocks: intervalBlocks,
		OverrideBlocks: overrideBlocks,
		DelayBlocks:    delayBlocks,
		TxHashSet:      mapset.NewThreadUnsafeSet[string](),
		Mu:             sync.Mutex{},
	}
}

func (s *Scanner) Scan(from, to uint64) (logs []types.Log, addressTopicLogsMap AddressTopicLogsMap) {
	if from+s.DelayBlocks > to {
		return
	}
	to -= s.DelayBlocks
	s.Iterator.Shuffle()
	fetch := func(_from, _to uint64) {
		// Attention!!!Repeat scanning _to prevent missing blocks
		var query = ethereum.FilterQuery{
			ToBlock:   new(big.Int).SetUint64(_to),
			Addresses: s.Addresses,
			Topics:    s.Topics,
		}
		if _from > from+s.OverrideBlocks {
			query.FromBlock = new(big.Int).SetUint64(_from - s.OverrideBlocks)
		} else {
			query.FromBlock = new(big.Int).SetUint64(_from)
		}
		nLogs := s.FilterLogs(query)
		log.Printf("Scan: %v(%v)-%v Success!\n", _from, query.FromBlock, _to)
		if len(nLogs) > 0 {
			s.Mu.Lock()
			var hashID string
			for _, nLog := range nLogs {
				hashID = fmt.Sprintf("%v%v", nLog.TxHash, nLog.Index)
				if !s.TxHashSet.Contains(hashID) {
					s.TxHashSet.Add(hashID)
					logs = append(logs, nLog)
				}
			}
			s.Mu.Unlock()
		}
	}

	s.execute(from, to, s.IntervalBlocks, fetch)

	addressTopicLogsMap = make(AddressTopicLogsMap)
	for _, nLog := range logs {
		if addressTopicLogsMap[nLog.Address] == nil {
			addressTopicLogsMap[nLog.Address] = make(map[common.Hash][]*types.Log)
		}
		addressTopicLogsMap[nLog.Address][nLog.Topics[0]] = append(addressTopicLogsMap[nLog.Address][nLog.Topics[0]], &nLog)
	}

	return logs, addressTopicLogsMap
}

func (s *Scanner) execute(from, to, interval uint64, fc func(from, to uint64)) {
	ranges := to - from
	count := ranges / interval
	wg := sync.WaitGroup{}
	wg.Add(int(count))
	for i := uint64(0); i < count; i++ {
		go func(i uint64) {
			fc(from+i*interval, from+(i+1)*interval-1)
			wg.Done()
		}(i)
	}
	if from+count*interval < to {
		fc(from+count*interval, to)
	}
	wg.Wait()
}

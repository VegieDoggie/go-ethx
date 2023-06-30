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
	"math/big"
	"sync"
	"time"
)

// Clientx defines typed wrappers for the Ethereum RPC API of a set of the Ethereum Clients.
type Clientx struct {
	*Iterator[*ethclient.Client]
	Ctx             context.Context
	RpcMap          map[*ethclient.Client]string
	NotFoundBlocks  uint64
	AutoBlockNumber uint64
}

// NewClientx connects clients to the given URLs, to provide a reliable Ethereum RPC API call, includes
// a timer to regularly update block height(AutoBlockNumber).
// If weight <= 1, the weight is always 1.
// Note: If len(weightList) == 0, then default weight = 1 will be active.
func NewClientx(rpcList []string, weights []int, notFoundBlocks uint64, limiter ...*rate.Limiter) *Clientx {
	iterator, rpcMap := newClientIteratorWithWeight(rpcList, weights, limiter...)
	c := &Clientx{
		Ctx:             context.Background(),
		Iterator:        iterator,
		RpcMap:          rpcMap,
		NotFoundBlocks:  notFoundBlocks,
		AutoBlockNumber: 0,
	}
	go func() {
		queryTicker := time.NewTicker(time.Second)
		defer queryTicker.Stop()
		var (
			err         error
			blockNumber uint64
			client      *ethclient.Client
		)
		for {
			for {
				client = c.Next()
				blockNumber, err = client.BlockNumber(c.Ctx)
				if err != nil {
					log.Println(fmt.Sprintf("[ERROR] %v BlockNumber %v %v", c.RpcMap[client], blockNumber, err))
					<-queryTicker.C
					continue
				}
				break
			}
			if blockNumber > c.AutoBlockNumber {
				c.AutoBlockNumber = blockNumber
			}
			<-queryTicker.C
		}
	}()
	queryTicker := time.NewTicker(100 * time.Millisecond)
	defer queryTicker.Stop()
	for c.AutoBlockNumber == 0 {
		<-queryTicker.C
	}
	return c
}

// newClientIteratorWithWeight creates a clientIterator with wights.
func newClientIteratorWithWeight(rpcList []string, weightList []int, limiter ...*rate.Limiter) (clientIterator *Iterator[*ethclient.Client], rpcMap map[*ethclient.Client]string) {
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
			blockNumber, err := client.BlockNumber(context.TODO())
			if err == nil && blockNumber > 0 {
				for j := 1; j < weightList[i]; j++ {
					reliableClients = append(reliableClients, client)
				}
				reliableClients = append(reliableClients, client)
				rpcMap[client] = rpc
				continue
			}
		}
		log.Println(fmt.Sprintf("[ClientX] Bad Rpc: %v", rpc))
	}
	if len(reliableClients) == 0 {
		panic("No reliable rpc nodes connection!")
	}
	clientIterator = NewIterator[*ethclient.Client](reliableClients, limiter...).Shuffle()
	return
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
	var err error
	for i := 0; ; i++ {
		client := c.WaitNext()
		blockNumber, err = client.BlockNumber(c.Ctx)
		if err != nil || blockNumber == 0 {
			log.Println(fmt.Sprintf("[ERROR] %v BlockNumber %v %v", c.RpcMap[client], blockNumber, err))
			continue
		}
		return
	}
}

// ChainID retrieves the current chain ID for transaction replay protection.
func (c *Clientx) ChainID() (chainID *big.Int) {
	var err error
	for i := 0; ; i++ {
		client := c.WaitNext()
		chainID, err = client.ChainID(c.Ctx)
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] %v ChainID %v", c.RpcMap[client], err))
			continue
		}
		return
	}
}

// NetworkID returns the network ID.
func (c *Clientx) NetworkID() (networkID *big.Int) {
	var err error
	for i := 0; ; i++ {
		client := c.WaitNext()
		networkID, err = client.NetworkID(c.Ctx)
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] %v NetworkID %v", c.RpcMap[client], err))
			continue
		}
		return
	}
}

// BalanceAt returns the wei balance of the given account.
// The block number can be nil, in which case the balance is taken from the latest known block.
func (c *Clientx) BalanceAt(account, blockNumber any) (balance *big.Int) {
	var err error
	for i := 0; ; i++ {
		client := c.WaitNext()
		balance, err = client.BalanceAt(c.Ctx, Address(account), BigInt(blockNumber))
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] %v BalanceAt %v", c.RpcMap[client], err))
			continue
		}
		return
	}
}

// PendingBalanceAt returns the wei balance of the given account in the pending state.
func (c *Clientx) PendingBalanceAt(account any) (balance *big.Int) {
	var err error
	for i := 0; ; i++ {
		client := c.WaitNext()
		balance, err = client.PendingBalanceAt(c.Ctx, Address(account))
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] %v PendingBalanceAt %v", c.RpcMap[client], err))
			continue
		}
		return
	}
}

// NonceAt returns the account nonce of the given account.
// The block number can be nil, in which case the nonce is taken from the latest known block.
func (c *Clientx) NonceAt(account, blockNumber any) (nonce uint64) {
	var err error
	for i := 0; ; i++ {
		client := c.WaitNext()
		nonce, err = client.NonceAt(c.Ctx, Address(account), BigInt(blockNumber))
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] %v NonceAt %v", c.RpcMap[client], err))
			continue
		}
		return
	}
}

// PendingNonceAt returns the account nonce of the given account in the pending state.
// This is the nonce that should be used for the next transaction.
func (c *Clientx) PendingNonceAt(account any) (nonce uint64) {
	var err error
	for i := 0; ; i++ {
		client := c.WaitNext()
		nonce, err = client.PendingNonceAt(c.Ctx, Address(account))
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] %v PendingNonceAt %v", c.RpcMap[client], err))
			continue
		}
		return
	}
}

// FilterLogs executes a filter query.
func (c *Clientx) FilterLogs(q ethereum.FilterQuery) (logs []types.Log) {
	var err error
	for i := 0; ; i++ {
		client := c.WaitNext()
		logs, err = client.FilterLogs(c.Ctx, q)
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] %v FilterLogs %v", c.RpcMap[client], err))
			continue
		}
		return
	}
}

// SuggestGasPrice retrieves the currently suggested gas price to allow a timely
// execution of a transaction.
func (c *Clientx) SuggestGasPrice() (gasPrice *big.Int) {
	var err error
	for i := 0; ; i++ {
		client := c.WaitNext()
		gasPrice, err = client.SuggestGasPrice(c.Ctx)
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] %v SuggestGasPrice %v", c.RpcMap[client], err))
			continue
		}
		return
	}
}

// SuggestGasTipCap retrieves the currently suggested gas tip cap after 1559 to
// allow a timely execution of a transaction.
func (c *Clientx) SuggestGasTipCap() (gasTipCap *big.Int) {
	var err error
	for i := 0; ; i++ {
		client := c.WaitNext()
		gasTipCap, err = client.SuggestGasTipCap(c.Ctx)
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] %v SuggestGasTipCap %v", c.RpcMap[client], err))
			continue
		}
		return
	}
}

// FeeHistory retrieves the fee market history.
func (c *Clientx) FeeHistory(blockCount uint64, lastBlock any, rewardPercentiles []float64) (feeHistory *ethereum.FeeHistory) {
	var err error
	for i := 0; ; i++ {
		client := c.WaitNext()
		feeHistory, err = client.FeeHistory(c.Ctx, blockCount, BigInt(lastBlock), rewardPercentiles)
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] %v FeeHistory %v", c.RpcMap[client], err))
			continue
		}
		return
	}
}

// StorageAt returns the value of key in the contract storage of the given account.
// The block number can be nil, in which case the value is taken from the latest known block.
func (c *Clientx) StorageAt(account, keyHash, blockNumber any) (storage []byte) {
	var err error
	for i := 0; ; i++ {
		client := c.WaitNext()
		storage, err = client.StorageAt(c.Ctx, Address(account), Hash(keyHash), BigInt(blockNumber))
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] %v StorageAt %v", c.RpcMap[client], err))
			continue
		}
		return
	}
}

// PendingStorageAt returns the value of key in the contract storage of the given account in the pending state.
func (c *Clientx) PendingStorageAt(account, keyHash any) (storage []byte) {
	var err error
	for i := 0; ; i++ {
		client := c.WaitNext()
		storage, err = client.PendingStorageAt(c.Ctx, Address(account), Hash(keyHash))
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] %v PendingStorageAt %v", c.RpcMap[client], err))
			continue
		}
		return
	}
}

// CodeAt returns the contract code of the given account.
// The block number can be nil, in which case the code is taken from the latest known block.
func (c *Clientx) CodeAt(account, blockNumber any) (code []byte) {
	var err error
	for i := 0; ; i++ {
		client := c.WaitNext()
		code, err = client.CodeAt(c.Ctx, Address(account), BigInt(blockNumber))
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] %v CodeAt %v", c.RpcMap[client], err))
			continue
		}
		return
	}
}

// PendingCodeAt returns the contract code of the given account in the pending state.
func (c *Clientx) PendingCodeAt(account any) (code []byte) {
	var err error
	for i := 0; ; i++ {
		client := c.WaitNext()
		code, err = client.PendingCodeAt(c.Ctx, Address(account))
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] %v PendingCodeAt %v", c.RpcMap[client], err))
			continue
		}
		return
	}
}

// notFoundStopBlockNumber returns the stop blockNumber for the notFound-error.
func (c *Clientx) notFoundStopBlockNumber(notFoundBlocks ...uint64) uint64 {
	if len(notFoundBlocks) > 0 {
		return c.AutoBlockNumber + notFoundBlocks[0]
	}
	return c.AutoBlockNumber + c.NotFoundBlocks
}

// BlockByHash returns the given full block.
//
// Note that loading full blocks requires two requests. Use HeaderByHash
// if you don't need all transactions or uncle headers.
func (c *Clientx) BlockByHash(hash any, notFoundBlocks ...uint64) (block *types.Block, err error) {
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()
	var notFoundStopBlockNumber uint64
	for i := 0; ; i++ {
		client := c.WaitNext()
		block, err = client.BlockByHash(c.Ctx, Hash(hash))
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] %v BlockByHash %v", c.RpcMap[client], err))
			if errors.Is(err, ethereum.NotFound) {
				if notFoundStopBlockNumber == 0 {
					notFoundStopBlockNumber = c.notFoundStopBlockNumber(notFoundBlocks...)
				}
				if notFoundStopBlockNumber >= c.AutoBlockNumber {
					return
				}
				<-queryTicker.C
			}
			continue
		}
		return
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
	var notFoundStopBlockNumber uint64
	for i := 0; ; i++ {
		client := c.WaitNext()
		block, err = client.BlockByNumber(c.Ctx, BigInt(blockNumber))
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] %v BlockByNumber %v", c.RpcMap[client], err))
			if errors.Is(err, ethereum.NotFound) {
				if notFoundStopBlockNumber == 0 {
					notFoundStopBlockNumber = c.notFoundStopBlockNumber(notFoundBlocks...)
				}
				if notFoundStopBlockNumber >= c.AutoBlockNumber {
					return
				}
				<-queryTicker.C
			}
			continue
		}
		return
	}
}

// HeaderByHash returns the block header with the given hash.
func (c *Clientx) HeaderByHash(hash any, notFoundBlocks ...uint64) (header *types.Header, err error) {
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()
	var notFoundStopBlockNumber uint64
	for i := 0; ; i++ {
		client := c.WaitNext()
		header, err = client.HeaderByHash(c.Ctx, Hash(hash))
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] %v HeaderByHash %v", c.RpcMap[client], err))
			if errors.Is(err, ethereum.NotFound) {
				if notFoundStopBlockNumber == 0 {
					notFoundStopBlockNumber = c.notFoundStopBlockNumber(notFoundBlocks...)
				}
				if notFoundStopBlockNumber >= c.AutoBlockNumber {
					return
				}
				<-queryTicker.C
			}
			continue
		}
		return
	}
}

// HeaderByNumber returns a block header from the current canonical chain. If number is
// nil, the latest known header is returned.
func (c *Clientx) HeaderByNumber(blockNumber any, notFoundBlocks ...uint64) (header *types.Header, err error) {
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()
	var notFoundStopBlockNumber uint64
	for i := 0; ; i++ {
		client := c.WaitNext()
		header, err = client.HeaderByNumber(c.Ctx, BigInt(blockNumber))
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] %v HeaderByNumber %v", c.RpcMap[client], err))
			if errors.Is(err, ethereum.NotFound) {
				if notFoundStopBlockNumber == 0 {
					notFoundStopBlockNumber = c.notFoundStopBlockNumber(notFoundBlocks...)
				}
				if notFoundStopBlockNumber >= c.AutoBlockNumber {
					return
				}
				<-queryTicker.C
			}
			continue
		}
		return
	}
}

// TransactionByHash returns the transaction with the given hash.
func (c *Clientx) TransactionByHash(hash any, notFoundBlocks ...uint64) (tx *types.Transaction, isPending bool, err error) {
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()
	var notFoundStopBlockNumber uint64
	for i := 0; ; i++ {
		client := c.WaitNext()
		tx, isPending, err = client.TransactionByHash(c.Ctx, Hash(hash))
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] %v TransactionByHash %v", c.RpcMap[client], err))
			if errors.Is(err, ethereum.NotFound) {
				if notFoundStopBlockNumber == 0 {
					notFoundStopBlockNumber = c.notFoundStopBlockNumber(notFoundBlocks...)
				}
				if notFoundStopBlockNumber >= c.AutoBlockNumber {
					return
				}
				<-queryTicker.C
			}
			continue
		}
		return
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
	var notFoundStopBlockNumber uint64
	for i := 0; ; i++ {
		client := c.WaitNext()
		sender, err = client.TransactionSender(c.Ctx, tx, Hash(blockHash), index)
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] %v TransactionSender %v", c.RpcMap[client], err))
			if errors.Is(err, ethereum.NotFound) {
				if notFoundStopBlockNumber == 0 {
					notFoundStopBlockNumber = c.notFoundStopBlockNumber(notFoundBlocks...)
				}
				if notFoundStopBlockNumber >= c.AutoBlockNumber {
					return
				}
				<-queryTicker.C
			}
			continue
		}
		return
	}
}

// TransactionCount returns the total number of transactions in the given block.
func (c *Clientx) TransactionCount(blockHash any, notFoundBlocks ...uint64) (count uint, err error) {
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()
	var notFoundStopBlockNumber uint64
	for i := 0; ; i++ {
		client := c.WaitNext()
		count, err = client.TransactionCount(c.Ctx, Hash(blockHash))
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] %v TransactionCount %v", c.RpcMap[client], err))
			if errors.Is(err, ethereum.NotFound) {
				if notFoundStopBlockNumber == 0 {
					notFoundStopBlockNumber = c.notFoundStopBlockNumber(notFoundBlocks...)
				}
				if notFoundStopBlockNumber >= c.AutoBlockNumber {
					return
				}
				<-queryTicker.C
			}
			continue
		}
		return
	}
}

// PendingTransactionCount returns the total number of transactions in the pending state.
func (c *Clientx) PendingTransactionCount() (count uint) {
	var err error
	for i := 0; ; i++ {
		client := c.WaitNext()
		count, err = client.PendingTransactionCount(c.Ctx)
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] %v PendingTransactionCount %v", c.RpcMap[client], err))
			continue
		}
		return
	}
}

// TransactionInBlock returns a single transaction at index in the given block.
func (c *Clientx) TransactionInBlock(blockHash any, index uint, notFoundBlocks ...uint64) (tx *types.Transaction, err error) {
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()
	var notFoundStopBlockNumber uint64
	for i := 0; ; i++ {
		client := c.WaitNext()
		tx, err = client.TransactionInBlock(c.Ctx, Hash(blockHash), index)
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] %v TransactionInBlock %v", c.RpcMap[client], err))
			if errors.Is(err, ethereum.NotFound) {
				if notFoundStopBlockNumber == 0 {
					notFoundStopBlockNumber = c.notFoundStopBlockNumber(notFoundBlocks...)
				}
				if notFoundStopBlockNumber >= c.AutoBlockNumber {
					return
				}
				<-queryTicker.C
			}
			continue
		}
		return
	}
}

// TransactionReceipt returns the receipt of a transaction by transaction hash.
// Note that the receipt is not available for pending transactions.
func (c *Clientx) TransactionReceipt(txHash any, notFoundBlocks ...uint64) (receipt *types.Receipt, err error) {
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()
	var notFoundStopBlockNumber uint64
	for i := 0; ; i++ {
		client := c.WaitNext()
		receipt, err = client.TransactionReceipt(c.Ctx, Hash(txHash))
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] %v TransactionReceipt %v", c.RpcMap[client], err))
			if errors.Is(err, ethereum.NotFound) {
				if notFoundStopBlockNumber == 0 {
					notFoundStopBlockNumber = c.notFoundStopBlockNumber(notFoundBlocks...)
				}
				if notFoundStopBlockNumber >= c.AutoBlockNumber {
					return
				}
				<-queryTicker.C
			}
			continue
		}
		return
	}
}

// WaitMined waits for tx to be mined on the blockchain.
// It stops waiting when the context is canceled.
// ethereum/go-ethereum@v1.11.6/accounts/abi/bind/util.go:32
func (c *Clientx) WaitMined(tx *types.Transaction, confirmBlocks uint64, notFoundBlocks ...uint64) (*types.Receipt, error) {
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()
	var confirmStopBlockNumber, notFoundStopBlockNumber uint64
	for i := 0; ; i++ {
		client := c.WaitNext()
		receipt, err := client.TransactionReceipt(c.Ctx, tx.Hash())
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] %v WaitMined %v", c.RpcMap[client], err))
			if errors.Is(err, ethereum.NotFound) {
				if notFoundStopBlockNumber == 0 {
					notFoundStopBlockNumber = c.notFoundStopBlockNumber(notFoundBlocks...)
				}
				if notFoundStopBlockNumber >= c.AutoBlockNumber {
					return nil, err
				}
			}
		} else {
			if confirmStopBlockNumber == 0 {
				confirmStopBlockNumber = c.AutoBlockNumber + confirmBlocks
			}
			if confirmStopBlockNumber >= c.AutoBlockNumber {
				return receipt, nil
			}
		}
		<-queryTicker.C
	}
}

// WaitDeployed waits for a contract deployment transaction and returns the on-chain
// contract address when it is mined. It stops waiting when Ctx is canceled.
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
	fetch := func(from, to uint64) {
		// Attention!!!Repeat scanning to prevent missing blocks
		var query = ethereum.FilterQuery{
			ToBlock:   new(big.Int).SetUint64(to),
			Addresses: s.Addresses,
			Topics:    s.Topics,
		}
		if from >= s.OverrideBlocks {
			query.FromBlock = new(big.Int).SetUint64(from - s.OverrideBlocks)
		} else {
			query.FromBlock = new(big.Int).SetUint64(from)
		}
		nLogs := s.FilterLogs(query)
		log.Println(fmt.Sprintf("Scan start for: %v-%v Success!", from, to))
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
	fc(from+count*interval, to)
	wg.Wait()
}

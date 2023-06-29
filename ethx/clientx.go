package ethx

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"time"
)

// Clientx defines typed wrappers for the Ethereum RPC API of a set of the Ethereum Clients.
type Clientx struct {
	ctx             context.Context
	iterator        *Iterator[*ethclient.Client]
	logErrThreshold int
	notFoundBlocks  uint64
	autoBlockNumber uint64
	confirmBlocks   uint64
}

// NewLightClientx connects clients to the given URLs, to provide a reliable Ethereum RPC API call.
// If weight <= 1, the weight is always 1.
func NewLightClientx(clientIterator *Iterator[*ethclient.Client], confirmBlocks, notFoundBlocks uint64) *Clientx {
	return &Clientx{
		iterator:        clientIterator,
		ctx:             context.Background(),
		notFoundBlocks:  notFoundBlocks,
		logErrThreshold: 5,
		autoBlockNumber: 0,
		confirmBlocks:   confirmBlocks,
	}
}

// NewClientx connects clients to the given URLs, to provide a reliable Ethereum RPC API call, includes
// a timer to regularly update block height(autoBlockNumber).
// If weight <= 1, the weight is always 1.
func NewClientx(clientIterator *Iterator[*ethclient.Client], confirmBlocks, notFoundBlocks uint64) *Clientx {
	c := NewLightClientx(clientIterator, confirmBlocks, notFoundBlocks)
	go func() {
		queryTicker := time.NewTicker(time.Second)
		var err error
		var blockNumber uint64
		for {
			for i := 0; ; i++ {
				blockNumber, err = c.Next().BlockNumber(c.ctx)
				if err != nil {
					if i > c.logErrThreshold {
						log.Printf("[ERROR] autoBlockNumber: %v", err)
					}
					continue
				}
				break
			}
			if blockNumber > c.autoBlockNumber {
				c.autoBlockNumber = blockNumber
			}
			<-queryTicker.C
		}
	}()
	queryTicker := time.NewTicker(100 * time.Millisecond)
	defer queryTicker.Stop()
	for c.autoBlockNumber == 0 {
		<-queryTicker.C
	}
	return c
}

// Next returns the next Ethereum Client.
func (c *Clientx) Next() *ethclient.Client {
	return c.iterator.Next()
}

// WaitNext returns the next Ethereum Client with a wait time.
func (c *Clientx) WaitNext() *ethclient.Client {
	return c.iterator.WaitNext()
}

// Close all clients connections.
func (c *Clientx) Close() {
	clients := c.iterator.All()
	for i := range clients {
		clients[i].Close()
	}
}

// LocalBlockNumber returns the Current BlockNumber which is recorded in Clientx.
func (c *Clientx) LocalBlockNumber() (blockNumber uint64) {
	return c.autoBlockNumber
}

// BlockNumber returns the most recent block number
func (c *Clientx) BlockNumber() (blockNumber uint64) {
	var err error
	for i := 0; ; i++ {
		blockNumber, err = c.WaitNext().BlockNumber(c.ctx)
		if err != nil || blockNumber == 0 {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] BlockNumber: %v", err)
			}
			continue
		}
		return
	}
}

// ChainID retrieves the current chain ID for transaction replay protection.
func (c *Clientx) ChainID() (chainID *big.Int) {
	var err error
	for i := 0; ; i++ {
		chainID, err = c.WaitNext().ChainID(c.ctx)
		if err != nil {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] ChainID: %v", err)
			}
			continue
		}
		return
	}
}

// NetworkID returns the network ID.
func (c *Clientx) NetworkID() (networkID *big.Int) {
	var err error
	for i := 0; ; i++ {
		networkID, err = c.WaitNext().NetworkID(c.ctx)
		if err != nil {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] NetworkID: %v", err)
			}
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
		balance, err = c.WaitNext().BalanceAt(c.ctx, Address(account), BigInt(blockNumber))
		if err != nil {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] BalanceAt: %v", err)
			}
			continue
		}
		return
	}
}

// PendingBalanceAt returns the wei balance of the given account in the pending state.
func (c *Clientx) PendingBalanceAt(account any) (balance *big.Int) {
	var err error
	for i := 0; ; i++ {
		balance, err = c.WaitNext().PendingBalanceAt(c.ctx, Address(account))
		if err != nil {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] PendingBalanceAt: %v", err)
			}
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
		nonce, err = c.WaitNext().NonceAt(c.ctx, Address(account), BigInt(blockNumber))
		if err != nil {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] NonceAt: %v", err)
			}
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
		nonce, err = c.WaitNext().PendingNonceAt(c.ctx, Address(account))
		if err != nil {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] PendingNonceAt: %v", err)
			}
			continue
		}
		return
	}
}

// FilterLogs executes a filter query.
func (c *Clientx) FilterLogs(q ethereum.FilterQuery) (logs []types.Log) {
	var err error
	for i := 0; ; i++ {
		logs, err = c.WaitNext().FilterLogs(c.ctx, q)
		if err != nil {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] FilterLogs: %v", err)
			}
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
		gasPrice, err = c.WaitNext().SuggestGasPrice(c.ctx)
		if err != nil {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] SuggestGasPrice: %v", err)
			}
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
		gasTipCap, err = c.WaitNext().SuggestGasTipCap(c.ctx)
		if err != nil {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] SuggestGasTipCap: %v", err)
			}
			continue
		}
		return
	}
}

// FeeHistory retrieves the fee market history.
func (c *Clientx) FeeHistory(blockCount uint64, lastBlock any, rewardPercentiles []float64) (feeHistory *ethereum.FeeHistory) {
	var err error
	for i := 0; ; i++ {
		feeHistory, err = c.WaitNext().FeeHistory(c.ctx, blockCount, BigInt(lastBlock), rewardPercentiles)
		if err != nil {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] FeeHistory: %v", err)
			}
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
		storage, err = c.WaitNext().StorageAt(c.ctx, Address(account), Hash(keyHash), BigInt(blockNumber))
		if err != nil {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] StorageAt: %v", err)
			}
			continue
		}
		return
	}
}

// PendingStorageAt returns the value of key in the contract storage of the given account in the pending state.
func (c *Clientx) PendingStorageAt(account, keyHash any) (storage []byte) {
	var err error
	for i := 0; ; i++ {
		storage, err = c.WaitNext().PendingStorageAt(c.ctx, Address(account), Hash(keyHash))
		if err != nil {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] PendingStorageAt: %v", err)
			}
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
		code, err = c.WaitNext().CodeAt(c.ctx, Address(account), BigInt(blockNumber))
		if err != nil {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] CodeAt: %v", err)
			}
			continue
		}
		return
	}
}

// PendingCodeAt returns the contract code of the given account in the pending state.
func (c *Clientx) PendingCodeAt(account any) (code []byte) {
	var err error
	for i := 0; ; i++ {
		code, err = c.WaitNext().PendingCodeAt(c.ctx, Address(account))
		if err != nil {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] PendingCodeAt: %v", err)
			}
			continue
		}
		return
	}
}

// notFoundStopBlockNumber returns the stop blockNumber for the notFound-error.
func (c *Clientx) notFoundStopBlockNumber(notFoundBlocks ...uint64) uint64 {
	if len(notFoundBlocks) > 0 {
		return c.autoBlockNumber + notFoundBlocks[0]
	}
	return c.autoBlockNumber + c.notFoundBlocks
}

// BlockByHash returns the given full block.
//
// Note that loading full blocks requires two requests. Use HeaderByHash
// if you don't need all transactions or uncle headers.
func (c *Clientx) BlockByHash(hash any, notFoundBlocks ...uint64) (block *types.Block, err error) {
	notFoundStopBlockNumber, queryTicker := c.notFoundStopBlockNumber(notFoundBlocks...), time.NewTicker(time.Second)
	defer queryTicker.Stop()
	for i := 0; ; i++ {
		block, err = c.WaitNext().BlockByHash(c.ctx, Hash(hash))
		if err != nil {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] BlockByHash: %v", err)
			}
			if errors.Is(err, ethereum.NotFound) {
				if notFoundStopBlockNumber >= c.autoBlockNumber {
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
	notFoundStopBlockNumber, queryTicker := c.notFoundStopBlockNumber(notFoundBlocks...), time.NewTicker(time.Second)
	defer queryTicker.Stop()
	for i := 0; ; i++ {
		block, err = c.WaitNext().BlockByNumber(c.ctx, BigInt(blockNumber))
		if err != nil {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] BlockByNumber: %v", err)
			}
			if errors.Is(err, ethereum.NotFound) {
				if notFoundStopBlockNumber >= c.autoBlockNumber {
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
	notFoundStopBlockNumber, queryTicker := c.notFoundStopBlockNumber(notFoundBlocks...), time.NewTicker(time.Second)
	defer queryTicker.Stop()
	for i := 0; ; i++ {
		header, err = c.WaitNext().HeaderByHash(c.ctx, Hash(hash))
		if err != nil {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] HeaderByHash: %v", err)
			}
			if errors.Is(err, ethereum.NotFound) {
				if notFoundStopBlockNumber >= c.autoBlockNumber {
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
	notFoundStopBlockNumber, queryTicker := c.notFoundStopBlockNumber(notFoundBlocks...), time.NewTicker(time.Second)
	defer queryTicker.Stop()
	for i := 0; ; i++ {
		header, err = c.WaitNext().HeaderByNumber(c.ctx, BigInt(blockNumber))
		if err != nil {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] HeaderByNumber: %v", err)
			}
			if errors.Is(err, ethereum.NotFound) {
				if notFoundStopBlockNumber >= c.autoBlockNumber {
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
	notFoundStopBlockNumber, queryTicker := c.notFoundStopBlockNumber(notFoundBlocks...), time.NewTicker(time.Second)
	defer queryTicker.Stop()
	for i := 0; ; i++ {
		tx, isPending, err = c.WaitNext().TransactionByHash(c.ctx, Hash(hash))
		if err != nil {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] TransactionByHash: %v", err)
			}
			if errors.Is(err, ethereum.NotFound) {
				if notFoundStopBlockNumber >= c.autoBlockNumber {
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
	notFoundStopBlockNumber, queryTicker := c.notFoundStopBlockNumber(notFoundBlocks...), time.NewTicker(time.Second)
	defer queryTicker.Stop()
	for i := 0; ; i++ {
		sender, err = c.WaitNext().TransactionSender(c.ctx, tx, Hash(blockHash), index)
		if err != nil {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] TransactionSender: %v", err)
			}
			if errors.Is(err, ethereum.NotFound) {
				if notFoundStopBlockNumber >= c.autoBlockNumber {
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
	notFoundStopBlockNumber, queryTicker := c.notFoundStopBlockNumber(notFoundBlocks...), time.NewTicker(time.Second)
	defer queryTicker.Stop()
	for i := 0; ; i++ {
		count, err = c.WaitNext().TransactionCount(c.ctx, Hash(blockHash))
		if err != nil {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] TransactionCount: %v", err)
			}
			if errors.Is(err, ethereum.NotFound) {
				if notFoundStopBlockNumber >= c.autoBlockNumber {
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
		count, err = c.WaitNext().PendingTransactionCount(c.ctx)
		if err != nil {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] PendingTransactionCount: %v", err)
			}
			continue
		}
		return
	}
}

// TransactionInBlock returns a single transaction at index in the given block.
func (c *Clientx) TransactionInBlock(blockHash any, index uint, notFoundBlocks ...uint64) (tx *types.Transaction, err error) {
	notFoundStopBlockNumber, queryTicker := c.notFoundStopBlockNumber(notFoundBlocks...), time.NewTicker(time.Second)
	defer queryTicker.Stop()
	for i := 0; ; i++ {
		tx, err = c.WaitNext().TransactionInBlock(c.ctx, Hash(blockHash), index)
		if err != nil {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] TransactionInBlock: %v", err)
			}
			if errors.Is(err, ethereum.NotFound) {
				if notFoundStopBlockNumber >= c.autoBlockNumber {
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
	notFoundStopBlockNumber, queryTicker := c.notFoundStopBlockNumber(notFoundBlocks...), time.NewTicker(time.Second)
	defer queryTicker.Stop()
	for i := 0; ; i++ {
		receipt, err = c.WaitNext().TransactionReceipt(c.ctx, Hash(txHash))
		if err != nil {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] TransactionReceipt: %v", err)
			}
			if errors.Is(err, ethereum.NotFound) {
				if notFoundStopBlockNumber >= c.autoBlockNumber {
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
func (c *Clientx) WaitMined(tx *types.Transaction, notFoundBlocks ...uint64) (*types.Receipt, error) {
	confirmStopBlockNumber := c.autoBlockNumber + c.confirmBlocks
	notFoundStopBlockNumber, queryTicker := c.notFoundStopBlockNumber(notFoundBlocks...), time.NewTicker(time.Second)
	defer queryTicker.Stop()
	for i := 0; ; i++ {
		receipt, err := c.WaitNext().TransactionReceipt(c.ctx, tx.Hash())
		if err != nil {
			if i > c.logErrThreshold {
				log.Printf("[ERROR] WaitMined: %v", err)
			}
			if errors.Is(err, ethereum.NotFound) {
				if i <= c.logErrThreshold {
					log.Printf("[ERROR] WaitMined: %v", err)
				}
				if notFoundStopBlockNumber >= c.autoBlockNumber {
					return nil, err
				}
			}
		} else {
			if confirmStopBlockNumber >= c.autoBlockNumber {
				return receipt, nil
			}
		}
		<-queryTicker.C
	}
}

// WaitDeployed waits for a contract deployment transaction and returns the on-chain
// contract address when it is mined. It stops waiting when ctx is canceled.
func (c *Clientx) WaitDeployed(tx *types.Transaction, notFoundBlocks ...uint64) (common.Address, error) {
	if tx.To() != nil {
		return common.Address{}, errors.New("tx is not contract creation")
	}
	receipt, err := c.WaitMined(tx, notFoundBlocks...)
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

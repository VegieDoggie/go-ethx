package ethx

import (
	"context"
	mapset "github.com/deckarep/golang-set"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"sync"
)

type key int

const key0 = key(0)

type AddressTopic0LogsMap = map[common.Address]map[common.Hash][]*types.Log

func ScannerCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, key0, make(AddressTopic0LogsMap))
}

func ScannerCtxValue(ctx context.Context) (addressTopic0LogsMap AddressTopic0LogsMap) {
	if t := ctx.Value(key0); t != nil {
		return t.(AddressTopic0LogsMap)
	}
	return nil
}

func ScannerCtxLogs(ctx context.Context, contractAddress common.Address, topic0 common.Hash) (logs []*types.Log) {
	if addressTopic0LogsMap := ScannerCtxValue(ctx); addressTopic0LogsMap != nil {
		if addressTopic0LogsMap[contractAddress] != nil {
			return addressTopic0LogsMap[contractAddress][topic0]
		}
	}
	return nil
}

type Scanner struct {
	Addresses       []common.Address
	Topics          [][]common.Hash
	Iterator        *Iterator[*ethclient.Client]
	OverridePerScan uint64
	IntervalPerScan uint64
	DelayBlocks     uint64
	txHashSet       mapset.Set
	mu0             sync.Mutex
	mu1             sync.Mutex
}

func NewScanner(iterator *Iterator[*ethclient.Client], topics [][]common.Hash, addresses []common.Address) *Scanner {
	return &Scanner{
		Addresses:       addresses,
		Topics:          topics,
		OverridePerScan: uint64(500),
		IntervalPerScan: uint64(500),
		DelayBlocks:     uint64(3),
		Iterator:        iterator,
		txHashSet:       mapset.NewThreadUnsafeSet(),
		mu0:             sync.Mutex{},
		mu1:             sync.Mutex{},
	}
}

func (s *Scanner) Scan(ctx context.Context, from, to uint64) (logs []types.Log, addressTopic0LogsMap AddressTopic0LogsMap) {
	if b := s.mu0.TryLock(); !b {
		return
	}
	defer s.mu0.Unlock()

	fetch := func(from, to uint64) {
		var (
			err   error
			nLogs []types.Log
			query = ethereum.FilterQuery{
				ToBlock:   new(big.Int).SetUint64(to),
				Addresses: s.Addresses,
				Topics:    s.Topics,
			}
		)
		// Attention!!!Repeat scanning to prevent missing blocks
		if from > s.OverridePerScan {
			query.FromBlock = new(big.Int).SetUint64(from - s.OverridePerScan)
		} else {
			query.FromBlock = new(big.Int).SetUint64(from)
		}
		for {
			if nLogs, err = s.Iterator.WaitNext().FilterLogs(ctx, query); err == nil {
				break
			}
			log.Printf("Scan start for: %v-%v Failed!", from, to)
		}
		if len(nLogs) > 0 {
			s.mu1.Lock()
			for _, nLog := range nLogs {
				if !s.txHashSet.Contains(nLog.TxHash) {
					s.txHashSet.Add(nLog.TxHash)
					logs = append(logs, nLog)
				}
			}
			s.mu1.Unlock()
		}
		log.Printf("Scan start for: %v-%v Success!", from, to)
	}

	s.execute(from, to, s.IntervalPerScan, fetch)

	if t := ctx.Value(key0); t != nil {
		addressTopic0LogsMap = t.(AddressTopic0LogsMap)
		for k := range addressTopic0LogsMap {
			delete(addressTopic0LogsMap, k)
		}
	} else {
		addressTopic0LogsMap = make(AddressTopic0LogsMap)
	}
	for _, o := range logs {
		if addressTopic0LogsMap[o.Address] == nil {
			addressTopic0LogsMap[o.Address] = make(map[common.Hash][]*types.Log)
		}
		addressTopic0LogsMap[o.Address][o.Topics[0]] = append(addressTopic0LogsMap[o.Address][o.Topics[0]], &o)
	}

	return logs, addressTopic0LogsMap
}

func (s *Scanner) Reset() {
	s.txHashSet = mapset.NewThreadUnsafeSet()
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

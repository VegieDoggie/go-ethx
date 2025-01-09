package ethx

import (
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"math/big"
	"sync"
	"time"
)

type RawLogger struct {
	client    *Clientx
	addresses []common.Address
	topics    [][]common.Hash
	config    *ClientxConfig
	txHashSet mapset.Set[string]
	mu        sync.Mutex
}

// NewRawLogger returns the RawLogger
// EventConfig require IntervalBlocks + OverrideBlocks <= 2000, eg: 800,800
// addresses []common.Address : optional, could be nil
// topics topics [][]common.Hash : optional, could be nil
func (c *Clientx) NewRawLogger(addresses []common.Address, topics [][]common.Hash, eventConfig ...*ClientxConfig) *RawLogger {
	return &RawLogger{
		client:    c,
		addresses: addresses,
		topics:    topics,
		config:    c.mustClientxConfig(eventConfig),
		txHashSet: mapset.NewThreadUnsafeSet[string](),
		mu:        sync.Mutex{},
	}
}

// Filter get logs from any blocks range, eg: (0, 10000000)
// Example Usage :
//
//	for l := range chLogs {
//		// DO 1: log handle
//	}
//
// // DO 2: get next turn new start/from
// chNewStart := <-chNewStart
func (r *RawLogger) Filter(from, to uint64) (chLogs chan types.Log, chNewStart chan uint64) {
	chLogs = make(chan types.Log)
	chNewStart = make(chan uint64, 1)
	go func() {
		defer close(chLogs)
		hasEvent := false
		fc := func(_from, _to uint64) {
			// Attention!!!Repeat scanning _to prevent missing blocks
			var query = ethereum.FilterQuery{
				Addresses: r.addresses,
				Topics:    r.topics,
				FromBlock: new(big.Int).SetUint64(_from),
				ToBlock:   new(big.Int).SetUint64(_to),
			}
			nLogs := r.client.FilterLogs(query)
			if len(nLogs) > 0 {
				r.mu.Lock()
				var hashID string
				for _, nLog := range nLogs {
					hashID = fmt.Sprintf("%v%v", nLog.TxHash, nLog.Index)
					if !r.txHashSet.Contains(hashID) {
						if !hasEvent {
							hasEvent = true
						}
						r.txHashSet.Add(hashID)
						chLogs <- nLog
					}
				}
				r.mu.Unlock()
			}
		}
		log.Printf("Filter: [from=%v,to=%v] start... | %v\n", from, to, r.addresses)
		beforeTime := time.Now()
		newStart := segmentCallback(from, to, r.config.Event.Clone(), fc)
		if !hasEvent {
			if to > from && (to-from)/r.config.Event.IntervalBlocks > 0 {
				chNewStart <- to - r.config.Event.IntervalBlocks
			} else {
				chNewStart <- from
			}
		} else {
			chNewStart <- newStart
		}
		log.Printf("Filter(%v): [from=%v,to=%v] Success! | %v\n", time.Since(beforeTime), from, to, r.addresses)
	}()
	return chLogs, chNewStart
}

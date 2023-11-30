package ethx

import (
	"errors"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"log"
	"reflect"
	"time"
)

// NewMustContract is safe contract caller
func (c *Clientx) NewMustContract(constructor any, addressLike any, eventConfig ...EventConfig) *MustContract {
	return &MustContract{
		client:          c,
		constructor:     constructor,
		contractName:    getFuncName(constructor)[3:],
		contractAddress: Address(addressLike),
		maxErrNum:       999,
		eventConfig:     c.newEventConfig(eventConfig),
	}
}

type MustContract struct {
	client          *Clientx
	constructor     any
	contractName    string
	contractAddress common.Address
	maxErrNum       int
	eventConfig     EventConfig
}

// Call0 contract func safely
// 注意: 如果缺失*bind.CallOpts首参数，则会自动以nil补全 (方便查询)
func (m *MustContract) Call0(f any, args ...any) any {
	return m.Call(f, args...)[0]
}

// Call contract func safely
// 注意: 如果缺失*bind.CallOpts首参数，则会自动以nil补全 (方便查询)
func (m *MustContract) Call(f any, args ...any) []any {
	for i := 0; i < m.maxErrNum; i++ {
		ret, err := m.callContract(f, args...)
		if err != nil {
			continue
		}
		return ret
	}
	panic(errors.New(fmt.Sprintf("Call::%v exceed maxErrNum(%v), contract=%v", getFuncName(f), m.maxErrNum, m.contractAddress)))
}

// CallWithMaxErrNum fit unsafe action, eg: maybe write failed
// 注意: 如果缺失*bind.CallOpts首参数，则会自动以nil补全 (方便查询)
func (m *MustContract) CallWithMaxErrNum(maxErrNum int, f any, args ...any) []any {
	for i := 0; i < maxErrNum; i++ {
		ret, err := m.callContract(f, args...)
		if err != nil {
			continue
		}
		return ret
	}
	panic(errors.New(fmt.Sprintf("Call::%v exceed maxErrNum(%v), contract=%v", getFuncName(f), m.maxErrNum, m.contractAddress)))
}

// Subscribe contract event
// eventName is from ch, so just pass ch!
func (m *MustContract) Subscribe(ch any, from any, index ...any) (sub event.Subscription) {
	chType := reflect.TypeOf(ch)
	if chType.Kind() != reflect.Chan {
		panic(fmt.Errorf("Subscribe::`ch` param not `chan`!\n"))
	}
	chPtrType := chType.Elem()
	if chPtrType.Kind() != reflect.Ptr {
		panic(fmt.Errorf("Subscribe::`ch` param not `chan pointer`!\n"))
	}
	chPtrEnumType := chPtrType.Elem()
	if chPtrEnumType.Kind() != reflect.Struct {
		panic(fmt.Errorf("Subscribe::`ch` param not `chan pointer struct`!\n"))
	}
	eventName := chPtrEnumType.Name()[len(m.contractName):]
	chEvent, stop := m.subscribe(BigInt(from).Uint64(), eventName, index...)
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer m.ignoreCloseChannelPanic()
		defer close(chEvent)
		chValue := reflect.ValueOf(ch)
		for {
			select {
			case l := <-chEvent:
				chValue.Send(reflect.ValueOf(l))
			case <-quit:
				log.Printf("[WARN] Unsubscribe Actively: %v=%v [%v]\n", m.contractName, m.contractAddress, eventName)
				stop <- true
				return nil
			}
		}
	})
}

func (m *MustContract) callContract(f any, args ...any) ([]any, error) {
	client := m.client.it.WaitNext()
	instanceRet := callFunc(m.constructor, m.contractAddress, client)
	if err, ok := instanceRet[len(instanceRet)-1].(error); ok {
		panic(err)
	}
	ret := callStructMethod(instanceRet[0], f, args...)
	if i := len(ret) - 1; i >= 0 {
		if err, ok := ret[i].(error); ok {
			m.client.errorCallback(f, client, err)
			return nil, err
		}
		return ret[:i], nil
	}
	return nil, nil
}

func (m *MustContract) subscribe(from uint64, eventName string, index ...any) (chEvent chan any, stop chan bool) {
	chEvent, stop = make(chan any, 128), make(chan bool, 1)
	go func() {
		tick := time.NewTicker(m.client.miningInterval)
		defer tick.Stop()
		filterFcName := "Filter" + eventName
		// reflect.TypeOf(m.constructor) : func(common.Address, bind.ContractBackend) (*TestLog.TestLog, error)
		// reflect.TypeOf(m.constructor).Out(0) : *TestLog.TestLog
		// reflect.New(reflect.TypeOf(m.constructor).Out(0).Elem()) : new(TestLog.TestLog)
		// reflect.New(reflect.TypeOf(m.constructor).Out(0).Elem()).MethodByName(filterFcName).Type() : func(*bind.FilterOpts, []common.Address) (*TestLog.TestLogIndex1Iterator, error)
		paramNum := reflect.New(reflect.TypeOf(m.constructor).Out(0).Elem()).MethodByName(filterFcName).Type().NumIn()
		if diff := paramNum - len(index); diff > 1 {
			for i := 1; i < diff; i++ {
				index = append(index, nil)
			}
		}
		txHashSet := mapset.NewThreadUnsafeSet[string]()
		filterFc := func(from, to uint64) {
			defer m.ignoreCloseChannelPanic()
			log.Println("filterFc:", from, to)
			opts := &bind.FilterOpts{
				Start:   from,
				End:     &to,
				Context: m.client.ctx,
			}
			// reflect call c.FilterXxx()
			var iterator any
			switch len(index) {
			case 0:
				iterator = m.Call(filterFcName, opts)[0]
			case 1:
				iterator = m.Call(filterFcName, opts, index[0])[0]
			case 2:
				iterator = m.Call(filterFcName, opts, index[0], index[1])[0]
			case 3:
				iterator = m.Call(filterFcName, opts, index[0], index[1], index[2])[0]
			default:
				panic(errors.New("Subscribe::the number of event-indexes don't exceed 3.\n"))
			}
			// reflect *XxxIterator
			itValue := reflect.ValueOf(iterator)
			// reflect next func
			itNext := itValue.MethodByName("Next")
			// reflect current *event
			itEvent := itValue.Elem().FieldByName("Event")
			for {
				if itNext.Call(nil)[0].Interface().(bool) {
					// reflect the raw log of *event
					itRaw := itEvent.Elem().FieldByName("Raw").Interface().(types.Log)
					hashID := fmt.Sprintf("%v%v", itRaw.TxHash, itRaw.Index)
					if !txHashSet.Contains(hashID) {
						txHashSet.Add(hashID)
						chEvent <- itEvent.Interface()
					}
				} else {
					return
				}
			}
		}
		_from, _to := from, m.client.BlockNumber()
		for {
			_from = segmentCallback(_from, _to, m.eventConfig, filterFc)
			_to = m.client.BlockNumber()
			log.Println("from:", _from, "_to:", _to)
			select {
			case <-tick.C:
			case <-stop:
				close(stop)
				return
			}
		}
	}()

	return chEvent, stop
}

func (m *MustContract) ignoreCloseChannelPanic() {
	if err, ok := recover().(error); ok {
		if err.Error() != "send on closed channel" {
			panic(err)
		}
	}
}

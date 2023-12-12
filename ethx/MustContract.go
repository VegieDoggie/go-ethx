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
func (c *Clientx) NewMustContract(constructor any, addressLike any, config ...*ClientxConfig) *MustContract {
	return &MustContract{
		client:          c,
		constructor:     constructor,
		contractName:    getFuncName(constructor)[3:],
		contractAddress: Address(addressLike),
		config:          c.mustClientxConfig(config),
	}
}

type MustContract struct {
	client          *Clientx
	constructor     any
	contractName    string
	contractAddress common.Address
	config          *ClientxConfig
}

// Read0 read from contract safely, return first
// Attention: missing the first param/*bind.CallOpts is legal
func (m *MustContract) Read0(f any, args ...any) any {
	return m.Read(f, args...)[0]
}

// Read from contract safely, return all(not include last error)
// Attention: missing the first param/*bind.CallOpts is legal
func (m *MustContract) Read(f any, args ...any) []any {
	ret, err := m.Call(m.config.MaxMustErrNumR, f, args...)
	if err != nil {
		panic(err)
	}
	return ret
}

// Write to contract
// Attention: missing the first *bind.CallOpts or PrivateKey is illegal
func (m *MustContract) Write(f any, args ...any) (*types.Transaction, error) {
	rets, err := m.Call(m.config.MaxMustErrNumW, f, args...)
	if err != nil {
		return nil, err
	}
	return rets[0].(*types.Transaction), nil
}

// Call fit unsafe action, eg: maybe write failed
// If READ: missing the first *bind.CallOpts is legal
// If WRITE: missing the first *bind.CallOpts or PrivateKey is illegal
func (m *MustContract) Call(maxErrNum int, f any, args ...any) (ret []any, err error) {
	funcType := reflect.ValueOf(f).Type()
	paramNum := funcType.NumIn()
	if paramNum > 0 && callOptsPtrType.ConvertibleTo(funcType.In(0)) {
		missNum := paramNum - len(args)
		switch missNum {
		case 0:
			if args[0] != nil && !callOptsPtrType.ConvertibleTo(reflect.TypeOf(args[0])) {
				args[0] = m.client.TransactOpts(args[0])
			}
		case 1:
			if len(args) == 0 || (args[0] != nil && !callOptsPtrType.ConvertibleTo(reflect.TypeOf(args[0]))) {
				args = append([]any{nil}, args...)
			}
		}
	}
	for i := 0; i < maxErrNum; i++ {
		ret, err = m.callContract(f, args...)
		if err != nil {
			continue
		}
		return ret, nil
	}
	return nil, fmt.Errorf("Call::%v exceed maxErrNumR(%v), contract=%v, err=%v\n", getFuncName(f), maxErrNum, m.contractAddress, err)
}

// Subscribe contract event
// eventName is from ch, so just pass ch!
func (m *MustContract) Subscribe(ch any, from any, index ...any) (sub event.Subscription, blockNumber chan uint64) {
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
	chEvent, blockNumber, stop := m.subscribe(BigInt(from).Uint64(), eventName, index...)
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
	}), blockNumber
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

func (m *MustContract) subscribe(from uint64, eventName string, index ...any) (chEvent chan any, blockNumber chan uint64, stop chan bool) {
	chEvent, blockNumber, stop = make(chan any, 512), make(chan uint64), make(chan bool, 1)
	go func() {
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
			opts := &bind.FilterOpts{
				Start:   from,
				End:     &to,
				Context: m.client.ctx,
			}
			// reflect call c.FilterXxx()
			var iterator any
			switch len(index) {
			case 0:
				iterator = m.Read(filterFcName, opts)[0]
			case 1:
				iterator = m.Read(filterFcName, opts, index[0])[0]
			case 2:
				iterator = m.Read(filterFcName, opts, index[0], index[1])[0]
			case 3:
				iterator = m.Read(filterFcName, opts, index[0], index[1], index[2])[0]
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

		subFrom, subTo := from, m.client.BlockNumber()
		subTicker := time.NewTicker(time.Second)
		defer subTicker.Stop()
		for {
			subFrom = segmentCallback(subFrom, subTo, m.config.Event, filterFc)
			go func() {
				blockNumber <- subFrom
			}()
			subTo = m.client.BlockNumber()
			select {
			case <-subTicker.C:
			case <-stop:
				close(stop)
				return
			}
		}
	}()

	return chEvent, blockNumber, stop
}

func (m *MustContract) ignoreCloseChannelPanic() {
	if err, ok := recover().(error); ok {
		if err.Error() != "send on closed channel" {
			panic(err)
		}
	}
}

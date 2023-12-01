// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package TestLog

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// TestLogMetaData contains all meta data concerning the TestLog contract.
var TestLogMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"u0\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"u1\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"u2\",\"type\":\"address\"}],\"name\":\"Index0\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"u0\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"u1\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"u2\",\"type\":\"address\"}],\"name\":\"Index1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"u0\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"u1\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"u2\",\"type\":\"address\"}],\"name\":\"Index2\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"u0\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"u1\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"u2\",\"type\":\"address\"}],\"name\":\"Index3\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"u0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"u1\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"u2\",\"type\":\"address\"}],\"name\":\"emitLog0\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"u0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"u1\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"u2\",\"type\":\"address\"}],\"name\":\"emitLog1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"u0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"u1\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"u2\",\"type\":\"address\"}],\"name\":\"emitLog2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"u0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"u1\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"u2\",\"type\":\"address\"}],\"name\":\"emitLog3\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getData0\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"data\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"d0\",\"type\":\"uint256\"}],\"name\":\"getData1\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"data\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"d0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"d1\",\"type\":\"uint256\"}],\"name\":\"getData2\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"data\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// TestLogABI is the input ABI used to generate the binding from.
// Deprecated: Use TestLogMetaData.ABI instead.
var TestLogABI = TestLogMetaData.ABI

// TestLog is an auto generated Go binding around an Ethereum contract.
type TestLog struct {
	TestLogCaller     // Read-only binding to the contract
	TestLogTransactor // Write-only binding to the contract
	TestLogFilterer   // Log filterer for contract events
}

// TestLogCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestLogCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestLogTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestLogTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestLogFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestLogFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestLogSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestLogSession struct {
	Contract     *TestLog          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Read options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestLogCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestLogCallerSession struct {
	Contract *TestLogCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Read options to use throughout this session
}

// TestLogTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestLogTransactorSession struct {
	Contract     *TestLogTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// TestLogRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestLogRaw struct {
	Contract *TestLog // Generic contract binding to access the raw methods on
}

// TestLogCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestLogCallerRaw struct {
	Contract *TestLogCaller // Generic read-only contract binding to access the raw methods on
}

// TestLogTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestLogTransactorRaw struct {
	Contract *TestLogTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestLog creates a new instance of TestLog, bound to a specific deployed contract.
func NewTestLog(address common.Address, backend bind.ContractBackend) (*TestLog, error) {
	contract, err := bindTestLog(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestLog{TestLogCaller: TestLogCaller{contract: contract}, TestLogTransactor: TestLogTransactor{contract: contract}, TestLogFilterer: TestLogFilterer{contract: contract}}, nil
}

// NewTestLogCaller creates a new read-only instance of TestLog, bound to a specific deployed contract.
func NewTestLogCaller(address common.Address, caller bind.ContractCaller) (*TestLogCaller, error) {
	contract, err := bindTestLog(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestLogCaller{contract: contract}, nil
}

// NewTestLogTransactor creates a new write-only instance of TestLog, bound to a specific deployed contract.
func NewTestLogTransactor(address common.Address, transactor bind.ContractTransactor) (*TestLogTransactor, error) {
	contract, err := bindTestLog(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestLogTransactor{contract: contract}, nil
}

// NewTestLogFilterer creates a new log filterer instance of TestLog, bound to a specific deployed contract.
func NewTestLogFilterer(address common.Address, filterer bind.ContractFilterer) (*TestLogFilterer, error) {
	contract, err := bindTestLog(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestLogFilterer{contract: contract}, nil
}

// bindTestLog binds a generic wrapper to an already deployed contract.
func bindTestLog(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestLogMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestLog *TestLogRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestLog.Contract.TestLogCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestLog *TestLogRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestLog.Contract.TestLogTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestLog *TestLogRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestLog.Contract.TestLogTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestLog *TestLogCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestLog.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestLog *TestLogTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestLog.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestLog *TestLogTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestLog.Contract.contract.Transact(opts, method, params...)
}

// GetData0 is a free data retrieval call binding the contract method 0xd42d2a9a.
//
// Solidity: function getData0() pure returns(uint256 data)
func (_TestLog *TestLogCaller) GetData0(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestLog.contract.Call(opts, &out, "getData0")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetData0 is a free data retrieval call binding the contract method 0xd42d2a9a.
//
// Solidity: function getData0() pure returns(uint256 data)
func (_TestLog *TestLogSession) GetData0() (*big.Int, error) {
	return _TestLog.Contract.GetData0(&_TestLog.CallOpts)
}

// GetData0 is a free data retrieval call binding the contract method 0xd42d2a9a.
//
// Solidity: function getData0() pure returns(uint256 data)
func (_TestLog *TestLogCallerSession) GetData0() (*big.Int, error) {
	return _TestLog.Contract.GetData0(&_TestLog.CallOpts)
}

// GetData1 is a free data retrieval call binding the contract method 0xf25c8077.
//
// Solidity: function getData1(uint256 d0) pure returns(uint256 data)
func (_TestLog *TestLogCaller) GetData1(opts *bind.CallOpts, d0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _TestLog.contract.Call(opts, &out, "getData1", d0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetData1 is a free data retrieval call binding the contract method 0xf25c8077.
//
// Solidity: function getData1(uint256 d0) pure returns(uint256 data)
func (_TestLog *TestLogSession) GetData1(d0 *big.Int) (*big.Int, error) {
	return _TestLog.Contract.GetData1(&_TestLog.CallOpts, d0)
}

// GetData1 is a free data retrieval call binding the contract method 0xf25c8077.
//
// Solidity: function getData1(uint256 d0) pure returns(uint256 data)
func (_TestLog *TestLogCallerSession) GetData1(d0 *big.Int) (*big.Int, error) {
	return _TestLog.Contract.GetData1(&_TestLog.CallOpts, d0)
}

// GetData2 is a free data retrieval call binding the contract method 0x8657472e.
//
// Solidity: function getData2(uint256 d0, uint256 d1) pure returns(uint256 data)
func (_TestLog *TestLogCaller) GetData2(opts *bind.CallOpts, d0 *big.Int, d1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _TestLog.contract.Call(opts, &out, "getData2", d0, d1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetData2 is a free data retrieval call binding the contract method 0x8657472e.
//
// Solidity: function getData2(uint256 d0, uint256 d1) pure returns(uint256 data)
func (_TestLog *TestLogSession) GetData2(d0 *big.Int, d1 *big.Int) (*big.Int, error) {
	return _TestLog.Contract.GetData2(&_TestLog.CallOpts, d0, d1)
}

// GetData2 is a free data retrieval call binding the contract method 0x8657472e.
//
// Solidity: function getData2(uint256 d0, uint256 d1) pure returns(uint256 data)
func (_TestLog *TestLogCallerSession) GetData2(d0 *big.Int, d1 *big.Int) (*big.Int, error) {
	return _TestLog.Contract.GetData2(&_TestLog.CallOpts, d0, d1)
}

// EmitLog0 is a paid mutator transaction binding the contract method 0xad770894.
//
// Solidity: function emitLog0(address u0, address u1, address u2) returns()
func (_TestLog *TestLogTransactor) EmitLog0(opts *bind.TransactOpts, u0 common.Address, u1 common.Address, u2 common.Address) (*types.Transaction, error) {
	return _TestLog.contract.Transact(opts, "emitLog0", u0, u1, u2)
}

// EmitLog0 is a paid mutator transaction binding the contract method 0xad770894.
//
// Solidity: function emitLog0(address u0, address u1, address u2) returns()
func (_TestLog *TestLogSession) EmitLog0(u0 common.Address, u1 common.Address, u2 common.Address) (*types.Transaction, error) {
	return _TestLog.Contract.EmitLog0(&_TestLog.TransactOpts, u0, u1, u2)
}

// EmitLog0 is a paid mutator transaction binding the contract method 0xad770894.
//
// Solidity: function emitLog0(address u0, address u1, address u2) returns()
func (_TestLog *TestLogTransactorSession) EmitLog0(u0 common.Address, u1 common.Address, u2 common.Address) (*types.Transaction, error) {
	return _TestLog.Contract.EmitLog0(&_TestLog.TransactOpts, u0, u1, u2)
}

// EmitLog1 is a paid mutator transaction binding the contract method 0x6ffa8941.
//
// Solidity: function emitLog1(address u0, address u1, address u2) returns()
func (_TestLog *TestLogTransactor) EmitLog1(opts *bind.TransactOpts, u0 common.Address, u1 common.Address, u2 common.Address) (*types.Transaction, error) {
	return _TestLog.contract.Transact(opts, "emitLog1", u0, u1, u2)
}

// EmitLog1 is a paid mutator transaction binding the contract method 0x6ffa8941.
//
// Solidity: function emitLog1(address u0, address u1, address u2) returns()
func (_TestLog *TestLogSession) EmitLog1(u0 common.Address, u1 common.Address, u2 common.Address) (*types.Transaction, error) {
	return _TestLog.Contract.EmitLog1(&_TestLog.TransactOpts, u0, u1, u2)
}

// EmitLog1 is a paid mutator transaction binding the contract method 0x6ffa8941.
//
// Solidity: function emitLog1(address u0, address u1, address u2) returns()
func (_TestLog *TestLogTransactorSession) EmitLog1(u0 common.Address, u1 common.Address, u2 common.Address) (*types.Transaction, error) {
	return _TestLog.Contract.EmitLog1(&_TestLog.TransactOpts, u0, u1, u2)
}

// EmitLog2 is a paid mutator transaction binding the contract method 0xf6401221.
//
// Solidity: function emitLog2(address u0, address u1, address u2) returns()
func (_TestLog *TestLogTransactor) EmitLog2(opts *bind.TransactOpts, u0 common.Address, u1 common.Address, u2 common.Address) (*types.Transaction, error) {
	return _TestLog.contract.Transact(opts, "emitLog2", u0, u1, u2)
}

// EmitLog2 is a paid mutator transaction binding the contract method 0xf6401221.
//
// Solidity: function emitLog2(address u0, address u1, address u2) returns()
func (_TestLog *TestLogSession) EmitLog2(u0 common.Address, u1 common.Address, u2 common.Address) (*types.Transaction, error) {
	return _TestLog.Contract.EmitLog2(&_TestLog.TransactOpts, u0, u1, u2)
}

// EmitLog2 is a paid mutator transaction binding the contract method 0xf6401221.
//
// Solidity: function emitLog2(address u0, address u1, address u2) returns()
func (_TestLog *TestLogTransactorSession) EmitLog2(u0 common.Address, u1 common.Address, u2 common.Address) (*types.Transaction, error) {
	return _TestLog.Contract.EmitLog2(&_TestLog.TransactOpts, u0, u1, u2)
}

// EmitLog3 is a paid mutator transaction binding the contract method 0xf71853e6.
//
// Solidity: function emitLog3(address u0, address u1, address u2) returns()
func (_TestLog *TestLogTransactor) EmitLog3(opts *bind.TransactOpts, u0 common.Address, u1 common.Address, u2 common.Address) (*types.Transaction, error) {
	return _TestLog.contract.Transact(opts, "emitLog3", u0, u1, u2)
}

// EmitLog3 is a paid mutator transaction binding the contract method 0xf71853e6.
//
// Solidity: function emitLog3(address u0, address u1, address u2) returns()
func (_TestLog *TestLogSession) EmitLog3(u0 common.Address, u1 common.Address, u2 common.Address) (*types.Transaction, error) {
	return _TestLog.Contract.EmitLog3(&_TestLog.TransactOpts, u0, u1, u2)
}

// EmitLog3 is a paid mutator transaction binding the contract method 0xf71853e6.
//
// Solidity: function emitLog3(address u0, address u1, address u2) returns()
func (_TestLog *TestLogTransactorSession) EmitLog3(u0 common.Address, u1 common.Address, u2 common.Address) (*types.Transaction, error) {
	return _TestLog.Contract.EmitLog3(&_TestLog.TransactOpts, u0, u1, u2)
}

// TestLogIndex0Iterator is returned from FilterIndex0 and is used to iterate over the raw logs and unpacked data for Index0 events raised by the TestLog contract.
type TestLogIndex0Iterator struct {
	Event *TestLogIndex0 // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TestLogIndex0Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestLogIndex0)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TestLogIndex0)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TestLogIndex0Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestLogIndex0Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestLogIndex0 represents a Index0 event raised by the TestLog contract.
type TestLogIndex0 struct {
	U0  common.Address
	U1  common.Address
	U2  common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterIndex0 is a free log retrieval operation binding the contract event 0xa81bd2a9a42530a0660fd7c21ebe71f482579c5f07a70f990066c233109132d8.
//
// Solidity: event Index0(address u0, address u1, address u2)
func (_TestLog *TestLogFilterer) FilterIndex0(opts *bind.FilterOpts) (*TestLogIndex0Iterator, error) {

	logs, sub, err := _TestLog.contract.FilterLogs(opts, "Index0")
	if err != nil {
		return nil, err
	}
	return &TestLogIndex0Iterator{contract: _TestLog.contract, event: "Index0", logs: logs, sub: sub}, nil
}

// WatchIndex0 is a free log subscription operation binding the contract event 0xa81bd2a9a42530a0660fd7c21ebe71f482579c5f07a70f990066c233109132d8.
//
// Solidity: event Index0(address u0, address u1, address u2)
func (_TestLog *TestLogFilterer) WatchIndex0(opts *bind.WatchOpts, sink chan<- *TestLogIndex0) (event.Subscription, error) {

	logs, sub, err := _TestLog.contract.WatchLogs(opts, "Index0")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestLogIndex0)
				if err := _TestLog.contract.UnpackLog(event, "Index0", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseIndex0 is a log parse operation binding the contract event 0xa81bd2a9a42530a0660fd7c21ebe71f482579c5f07a70f990066c233109132d8.
//
// Solidity: event Index0(address u0, address u1, address u2)
func (_TestLog *TestLogFilterer) ParseIndex0(log types.Log) (*TestLogIndex0, error) {
	event := new(TestLogIndex0)
	if err := _TestLog.contract.UnpackLog(event, "Index0", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestLogIndex1Iterator is returned from FilterIndex1 and is used to iterate over the raw logs and unpacked data for Index1 events raised by the TestLog contract.
type TestLogIndex1Iterator struct {
	Event *TestLogIndex1 // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TestLogIndex1Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestLogIndex1)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TestLogIndex1)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TestLogIndex1Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestLogIndex1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestLogIndex1 represents a Index1 event raised by the TestLog contract.
type TestLogIndex1 struct {
	U0  common.Address
	U1  common.Address
	U2  common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterIndex1 is a free log retrieval operation binding the contract event 0xab82c07a250ac9261c7e58da7b01a893f71e96da4db4cad5809564f1851b8983.
//
// Solidity: event Index1(address indexed u0, address u1, address u2)
func (_TestLog *TestLogFilterer) FilterIndex1(opts *bind.FilterOpts, u0 []common.Address) (*TestLogIndex1Iterator, error) {

	var u0Rule []interface{}
	for _, u0Item := range u0 {
		u0Rule = append(u0Rule, u0Item)
	}

	logs, sub, err := _TestLog.contract.FilterLogs(opts, "Index1", u0Rule)
	if err != nil {
		return nil, err
	}
	return &TestLogIndex1Iterator{contract: _TestLog.contract, event: "Index1", logs: logs, sub: sub}, nil
}

// WatchIndex1 is a free log subscription operation binding the contract event 0xab82c07a250ac9261c7e58da7b01a893f71e96da4db4cad5809564f1851b8983.
//
// Solidity: event Index1(address indexed u0, address u1, address u2)
func (_TestLog *TestLogFilterer) WatchIndex1(opts *bind.WatchOpts, sink chan<- *TestLogIndex1, u0 []common.Address) (event.Subscription, error) {

	var u0Rule []interface{}
	for _, u0Item := range u0 {
		u0Rule = append(u0Rule, u0Item)
	}

	logs, sub, err := _TestLog.contract.WatchLogs(opts, "Index1", u0Rule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestLogIndex1)
				if err := _TestLog.contract.UnpackLog(event, "Index1", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseIndex1 is a log parse operation binding the contract event 0xab82c07a250ac9261c7e58da7b01a893f71e96da4db4cad5809564f1851b8983.
//
// Solidity: event Index1(address indexed u0, address u1, address u2)
func (_TestLog *TestLogFilterer) ParseIndex1(log types.Log) (*TestLogIndex1, error) {
	event := new(TestLogIndex1)
	if err := _TestLog.contract.UnpackLog(event, "Index1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestLogIndex2Iterator is returned from FilterIndex2 and is used to iterate over the raw logs and unpacked data for Index2 events raised by the TestLog contract.
type TestLogIndex2Iterator struct {
	Event *TestLogIndex2 // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TestLogIndex2Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestLogIndex2)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TestLogIndex2)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TestLogIndex2Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestLogIndex2Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestLogIndex2 represents a Index2 event raised by the TestLog contract.
type TestLogIndex2 struct {
	U0  common.Address
	U1  common.Address
	U2  common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterIndex2 is a free log retrieval operation binding the contract event 0x861831cf4cdb2c9aafaff6c8ea3a04b507fd8ab94687eb078b0666aa7aa7832e.
//
// Solidity: event Index2(address indexed u0, address indexed u1, address u2)
func (_TestLog *TestLogFilterer) FilterIndex2(opts *bind.FilterOpts, u0 []common.Address, u1 []common.Address) (*TestLogIndex2Iterator, error) {

	var u0Rule []interface{}
	for _, u0Item := range u0 {
		u0Rule = append(u0Rule, u0Item)
	}
	var u1Rule []interface{}
	for _, u1Item := range u1 {
		u1Rule = append(u1Rule, u1Item)
	}

	logs, sub, err := _TestLog.contract.FilterLogs(opts, "Index2", u0Rule, u1Rule)
	if err != nil {
		return nil, err
	}
	return &TestLogIndex2Iterator{contract: _TestLog.contract, event: "Index2", logs: logs, sub: sub}, nil
}

// WatchIndex2 is a free log subscription operation binding the contract event 0x861831cf4cdb2c9aafaff6c8ea3a04b507fd8ab94687eb078b0666aa7aa7832e.
//
// Solidity: event Index2(address indexed u0, address indexed u1, address u2)
func (_TestLog *TestLogFilterer) WatchIndex2(opts *bind.WatchOpts, sink chan<- *TestLogIndex2, u0 []common.Address, u1 []common.Address) (event.Subscription, error) {

	var u0Rule []interface{}
	for _, u0Item := range u0 {
		u0Rule = append(u0Rule, u0Item)
	}
	var u1Rule []interface{}
	for _, u1Item := range u1 {
		u1Rule = append(u1Rule, u1Item)
	}

	logs, sub, err := _TestLog.contract.WatchLogs(opts, "Index2", u0Rule, u1Rule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestLogIndex2)
				if err := _TestLog.contract.UnpackLog(event, "Index2", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseIndex2 is a log parse operation binding the contract event 0x861831cf4cdb2c9aafaff6c8ea3a04b507fd8ab94687eb078b0666aa7aa7832e.
//
// Solidity: event Index2(address indexed u0, address indexed u1, address u2)
func (_TestLog *TestLogFilterer) ParseIndex2(log types.Log) (*TestLogIndex2, error) {
	event := new(TestLogIndex2)
	if err := _TestLog.contract.UnpackLog(event, "Index2", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestLogIndex3Iterator is returned from FilterIndex3 and is used to iterate over the raw logs and unpacked data for Index3 events raised by the TestLog contract.
type TestLogIndex3Iterator struct {
	Event *TestLogIndex3 // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TestLogIndex3Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestLogIndex3)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TestLogIndex3)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TestLogIndex3Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestLogIndex3Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestLogIndex3 represents a Index3 event raised by the TestLog contract.
type TestLogIndex3 struct {
	U0  common.Address
	U1  common.Address
	U2  common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterIndex3 is a free log retrieval operation binding the contract event 0xd24059a7b65a58c0f1a52da5b1782177ddbf131f39250fc9f6a09a54f039d84d.
//
// Solidity: event Index3(address indexed u0, address indexed u1, address indexed u2)
func (_TestLog *TestLogFilterer) FilterIndex3(opts *bind.FilterOpts, u0 []common.Address, u1 []common.Address, u2 []common.Address) (*TestLogIndex3Iterator, error) {

	var u0Rule []interface{}
	for _, u0Item := range u0 {
		u0Rule = append(u0Rule, u0Item)
	}
	var u1Rule []interface{}
	for _, u1Item := range u1 {
		u1Rule = append(u1Rule, u1Item)
	}
	var u2Rule []interface{}
	for _, u2Item := range u2 {
		u2Rule = append(u2Rule, u2Item)
	}

	logs, sub, err := _TestLog.contract.FilterLogs(opts, "Index3", u0Rule, u1Rule, u2Rule)
	if err != nil {
		return nil, err
	}
	return &TestLogIndex3Iterator{contract: _TestLog.contract, event: "Index3", logs: logs, sub: sub}, nil
}

// WatchIndex3 is a free log subscription operation binding the contract event 0xd24059a7b65a58c0f1a52da5b1782177ddbf131f39250fc9f6a09a54f039d84d.
//
// Solidity: event Index3(address indexed u0, address indexed u1, address indexed u2)
func (_TestLog *TestLogFilterer) WatchIndex3(opts *bind.WatchOpts, sink chan<- *TestLogIndex3, u0 []common.Address, u1 []common.Address, u2 []common.Address) (event.Subscription, error) {

	var u0Rule []interface{}
	for _, u0Item := range u0 {
		u0Rule = append(u0Rule, u0Item)
	}
	var u1Rule []interface{}
	for _, u1Item := range u1 {
		u1Rule = append(u1Rule, u1Item)
	}
	var u2Rule []interface{}
	for _, u2Item := range u2 {
		u2Rule = append(u2Rule, u2Item)
	}

	logs, sub, err := _TestLog.contract.WatchLogs(opts, "Index3", u0Rule, u1Rule, u2Rule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestLogIndex3)
				if err := _TestLog.contract.UnpackLog(event, "Index3", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseIndex3 is a log parse operation binding the contract event 0xd24059a7b65a58c0f1a52da5b1782177ddbf131f39250fc9f6a09a54f039d84d.
//
// Solidity: event Index3(address indexed u0, address indexed u1, address indexed u2)
func (_TestLog *TestLogFilterer) ParseIndex3(log types.Log) (*TestLogIndex3, error) {
	event := new(TestLogIndex3)
	if err := _TestLog.contract.UnpackLog(event, "Index3", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

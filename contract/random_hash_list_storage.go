// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

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
)

// RandomHashListStorageMetaData contains all meta data concerning the RandomHashListStorage contract.
var RandomHashListStorageMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"reqID\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"collection\",\"type\":\"string[]\"}],\"name\":\"CollectionSaved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"reqID\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"hashList\",\"type\":\"string[]\"}],\"name\":\"HashListStored\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"reqID\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"collection\",\"type\":\"string[]\"}],\"name\":\"generatorSaveCollection\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"reqID\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getCollectionByReqIDAndAddress\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"reqID\",\"type\":\"string\"}],\"name\":\"getHashList\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"reqID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"hash\",\"type\":\"string\"}],\"name\":\"hashExistsInList\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"reqID\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"_hashList\",\"type\":\"string[]\"}],\"name\":\"storeHashList\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// RandomHashListStorageABI is the input ABI used to generate the binding from.
// Deprecated: Use RandomHashListStorageMetaData.ABI instead.
var RandomHashListStorageABI = RandomHashListStorageMetaData.ABI

// RandomHashListStorage is an auto generated Go binding around an Ethereum contract.
type RandomHashListStorage struct {
	RandomHashListStorageCaller     // Read-only binding to the contract
	RandomHashListStorageTransactor // Write-only binding to the contract
	RandomHashListStorageFilterer   // Log filterer for contract events
}

// RandomHashListStorageCaller is an auto generated read-only Go binding around an Ethereum contract.
type RandomHashListStorageCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RandomHashListStorageTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RandomHashListStorageTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RandomHashListStorageFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RandomHashListStorageFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RandomHashListStorageSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RandomHashListStorageSession struct {
	Contract     *RandomHashListStorage // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// RandomHashListStorageCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RandomHashListStorageCallerSession struct {
	Contract *RandomHashListStorageCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// RandomHashListStorageTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RandomHashListStorageTransactorSession struct {
	Contract     *RandomHashListStorageTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// RandomHashListStorageRaw is an auto generated low-level Go binding around an Ethereum contract.
type RandomHashListStorageRaw struct {
	Contract *RandomHashListStorage // Generic contract binding to access the raw methods on
}

// RandomHashListStorageCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RandomHashListStorageCallerRaw struct {
	Contract *RandomHashListStorageCaller // Generic read-only contract binding to access the raw methods on
}

// RandomHashListStorageTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RandomHashListStorageTransactorRaw struct {
	Contract *RandomHashListStorageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRandomHashListStorage creates a new instance of RandomHashListStorage, bound to a specific deployed contract.
func NewRandomHashListStorage(address common.Address, backend bind.ContractBackend) (*RandomHashListStorage, error) {
	contract, err := bindRandomHashListStorage(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RandomHashListStorage{RandomHashListStorageCaller: RandomHashListStorageCaller{contract: contract}, RandomHashListStorageTransactor: RandomHashListStorageTransactor{contract: contract}, RandomHashListStorageFilterer: RandomHashListStorageFilterer{contract: contract}}, nil
}

// NewRandomHashListStorageCaller creates a new read-only instance of RandomHashListStorage, bound to a specific deployed contract.
func NewRandomHashListStorageCaller(address common.Address, caller bind.ContractCaller) (*RandomHashListStorageCaller, error) {
	contract, err := bindRandomHashListStorage(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RandomHashListStorageCaller{contract: contract}, nil
}

// NewRandomHashListStorageTransactor creates a new write-only instance of RandomHashListStorage, bound to a specific deployed contract.
func NewRandomHashListStorageTransactor(address common.Address, transactor bind.ContractTransactor) (*RandomHashListStorageTransactor, error) {
	contract, err := bindRandomHashListStorage(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RandomHashListStorageTransactor{contract: contract}, nil
}

// NewRandomHashListStorageFilterer creates a new log filterer instance of RandomHashListStorage, bound to a specific deployed contract.
func NewRandomHashListStorageFilterer(address common.Address, filterer bind.ContractFilterer) (*RandomHashListStorageFilterer, error) {
	contract, err := bindRandomHashListStorage(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RandomHashListStorageFilterer{contract: contract}, nil
}

// bindRandomHashListStorage binds a generic wrapper to an already deployed contract.
func bindRandomHashListStorage(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RandomHashListStorageABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RandomHashListStorage *RandomHashListStorageRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RandomHashListStorage.Contract.RandomHashListStorageCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RandomHashListStorage *RandomHashListStorageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RandomHashListStorage.Contract.RandomHashListStorageTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RandomHashListStorage *RandomHashListStorageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RandomHashListStorage.Contract.RandomHashListStorageTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RandomHashListStorage *RandomHashListStorageCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RandomHashListStorage.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RandomHashListStorage *RandomHashListStorageTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RandomHashListStorage.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RandomHashListStorage *RandomHashListStorageTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RandomHashListStorage.Contract.contract.Transact(opts, method, params...)
}

// GetCollectionByReqIDAndAddress is a free data retrieval call binding the contract method 0xf641f4a7.
//
// Solidity: function getCollectionByReqIDAndAddress(string reqID, address addr) view returns(string[])
func (_RandomHashListStorage *RandomHashListStorageCaller) GetCollectionByReqIDAndAddress(opts *bind.CallOpts, reqID string, addr common.Address) ([]string, error) {
	var out []interface{}
	err := _RandomHashListStorage.contract.Call(opts, &out, "getCollectionByReqIDAndAddress", reqID, addr)

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetCollectionByReqIDAndAddress is a free data retrieval call binding the contract method 0xf641f4a7.
//
// Solidity: function getCollectionByReqIDAndAddress(string reqID, address addr) view returns(string[])
func (_RandomHashListStorage *RandomHashListStorageSession) GetCollectionByReqIDAndAddress(reqID string, addr common.Address) ([]string, error) {
	return _RandomHashListStorage.Contract.GetCollectionByReqIDAndAddress(&_RandomHashListStorage.CallOpts, reqID, addr)
}

// GetCollectionByReqIDAndAddress is a free data retrieval call binding the contract method 0xf641f4a7.
//
// Solidity: function getCollectionByReqIDAndAddress(string reqID, address addr) view returns(string[])
func (_RandomHashListStorage *RandomHashListStorageCallerSession) GetCollectionByReqIDAndAddress(reqID string, addr common.Address) ([]string, error) {
	return _RandomHashListStorage.Contract.GetCollectionByReqIDAndAddress(&_RandomHashListStorage.CallOpts, reqID, addr)
}

// GetHashList is a free data retrieval call binding the contract method 0x9fd9d965.
//
// Solidity: function getHashList(string reqID) view returns(string[])
func (_RandomHashListStorage *RandomHashListStorageCaller) GetHashList(opts *bind.CallOpts, reqID string) ([]string, error) {
	var out []interface{}
	err := _RandomHashListStorage.contract.Call(opts, &out, "getHashList", reqID)

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetHashList is a free data retrieval call binding the contract method 0x9fd9d965.
//
// Solidity: function getHashList(string reqID) view returns(string[])
func (_RandomHashListStorage *RandomHashListStorageSession) GetHashList(reqID string) ([]string, error) {
	return _RandomHashListStorage.Contract.GetHashList(&_RandomHashListStorage.CallOpts, reqID)
}

// GetHashList is a free data retrieval call binding the contract method 0x9fd9d965.
//
// Solidity: function getHashList(string reqID) view returns(string[])
func (_RandomHashListStorage *RandomHashListStorageCallerSession) GetHashList(reqID string) ([]string, error) {
	return _RandomHashListStorage.Contract.GetHashList(&_RandomHashListStorage.CallOpts, reqID)
}

// HashExistsInList is a free data retrieval call binding the contract method 0xa9f6e826.
//
// Solidity: function hashExistsInList(string reqID, string hash) view returns(bool)
func (_RandomHashListStorage *RandomHashListStorageCaller) HashExistsInList(opts *bind.CallOpts, reqID string, hash string) (bool, error) {
	var out []interface{}
	err := _RandomHashListStorage.contract.Call(opts, &out, "hashExistsInList", reqID, hash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HashExistsInList is a free data retrieval call binding the contract method 0xa9f6e826.
//
// Solidity: function hashExistsInList(string reqID, string hash) view returns(bool)
func (_RandomHashListStorage *RandomHashListStorageSession) HashExistsInList(reqID string, hash string) (bool, error) {
	return _RandomHashListStorage.Contract.HashExistsInList(&_RandomHashListStorage.CallOpts, reqID, hash)
}

// HashExistsInList is a free data retrieval call binding the contract method 0xa9f6e826.
//
// Solidity: function hashExistsInList(string reqID, string hash) view returns(bool)
func (_RandomHashListStorage *RandomHashListStorageCallerSession) HashExistsInList(reqID string, hash string) (bool, error) {
	return _RandomHashListStorage.Contract.HashExistsInList(&_RandomHashListStorage.CallOpts, reqID, hash)
}

// GeneratorSaveCollection is a paid mutator transaction binding the contract method 0x01157cfe.
//
// Solidity: function generatorSaveCollection(string reqID, string[] collection) returns()
func (_RandomHashListStorage *RandomHashListStorageTransactor) GeneratorSaveCollection(opts *bind.TransactOpts, reqID string, collection []string) (*types.Transaction, error) {
	return _RandomHashListStorage.contract.Transact(opts, "generatorSaveCollection", reqID, collection)
}

// GeneratorSaveCollection is a paid mutator transaction binding the contract method 0x01157cfe.
//
// Solidity: function generatorSaveCollection(string reqID, string[] collection) returns()
func (_RandomHashListStorage *RandomHashListStorageSession) GeneratorSaveCollection(reqID string, collection []string) (*types.Transaction, error) {
	return _RandomHashListStorage.Contract.GeneratorSaveCollection(&_RandomHashListStorage.TransactOpts, reqID, collection)
}

// GeneratorSaveCollection is a paid mutator transaction binding the contract method 0x01157cfe.
//
// Solidity: function generatorSaveCollection(string reqID, string[] collection) returns()
func (_RandomHashListStorage *RandomHashListStorageTransactorSession) GeneratorSaveCollection(reqID string, collection []string) (*types.Transaction, error) {
	return _RandomHashListStorage.Contract.GeneratorSaveCollection(&_RandomHashListStorage.TransactOpts, reqID, collection)
}

// StoreHashList is a paid mutator transaction binding the contract method 0xa04500af.
//
// Solidity: function storeHashList(string reqID, string[] _hashList) returns()
func (_RandomHashListStorage *RandomHashListStorageTransactor) StoreHashList(opts *bind.TransactOpts, reqID string, _hashList []string) (*types.Transaction, error) {
	return _RandomHashListStorage.contract.Transact(opts, "storeHashList", reqID, _hashList)
}

// StoreHashList is a paid mutator transaction binding the contract method 0xa04500af.
//
// Solidity: function storeHashList(string reqID, string[] _hashList) returns()
func (_RandomHashListStorage *RandomHashListStorageSession) StoreHashList(reqID string, _hashList []string) (*types.Transaction, error) {
	return _RandomHashListStorage.Contract.StoreHashList(&_RandomHashListStorage.TransactOpts, reqID, _hashList)
}

// StoreHashList is a paid mutator transaction binding the contract method 0xa04500af.
//
// Solidity: function storeHashList(string reqID, string[] _hashList) returns()
func (_RandomHashListStorage *RandomHashListStorageTransactorSession) StoreHashList(reqID string, _hashList []string) (*types.Transaction, error) {
	return _RandomHashListStorage.Contract.StoreHashList(&_RandomHashListStorage.TransactOpts, reqID, _hashList)
}

// RandomHashListStorageCollectionSavedIterator is returned from FilterCollectionSaved and is used to iterate over the raw logs and unpacked data for CollectionSaved events raised by the RandomHashListStorage contract.
type RandomHashListStorageCollectionSavedIterator struct {
	Event *RandomHashListStorageCollectionSaved // Event containing the contract specifics and raw log

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
func (it *RandomHashListStorageCollectionSavedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RandomHashListStorageCollectionSaved)
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
		it.Event = new(RandomHashListStorageCollectionSaved)
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
func (it *RandomHashListStorageCollectionSavedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RandomHashListStorageCollectionSavedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RandomHashListStorageCollectionSaved represents a CollectionSaved event raised by the RandomHashListStorage contract.
type RandomHashListStorageCollectionSaved struct {
	ReqID      common.Hash
	Sender     common.Address
	Collection []string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCollectionSaved is a free log retrieval operation binding the contract event 0x5b2ef0047cedb4e7f134dc0e4c09bd0d9f38099655c02cf0ca57e58823fb21c0.
//
// Solidity: event CollectionSaved(string indexed reqID, address indexed sender, string[] collection)
func (_RandomHashListStorage *RandomHashListStorageFilterer) FilterCollectionSaved(opts *bind.FilterOpts, reqID []string, sender []common.Address) (*RandomHashListStorageCollectionSavedIterator, error) {

	var reqIDRule []interface{}
	for _, reqIDItem := range reqID {
		reqIDRule = append(reqIDRule, reqIDItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _RandomHashListStorage.contract.FilterLogs(opts, "CollectionSaved", reqIDRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &RandomHashListStorageCollectionSavedIterator{contract: _RandomHashListStorage.contract, event: "CollectionSaved", logs: logs, sub: sub}, nil
}

// WatchCollectionSaved is a free log subscription operation binding the contract event 0x5b2ef0047cedb4e7f134dc0e4c09bd0d9f38099655c02cf0ca57e58823fb21c0.
//
// Solidity: event CollectionSaved(string indexed reqID, address indexed sender, string[] collection)
func (_RandomHashListStorage *RandomHashListStorageFilterer) WatchCollectionSaved(opts *bind.WatchOpts, sink chan<- *RandomHashListStorageCollectionSaved, reqID []string, sender []common.Address) (event.Subscription, error) {

	var reqIDRule []interface{}
	for _, reqIDItem := range reqID {
		reqIDRule = append(reqIDRule, reqIDItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _RandomHashListStorage.contract.WatchLogs(opts, "CollectionSaved", reqIDRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RandomHashListStorageCollectionSaved)
				if err := _RandomHashListStorage.contract.UnpackLog(event, "CollectionSaved", log); err != nil {
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

// ParseCollectionSaved is a log parse operation binding the contract event 0x5b2ef0047cedb4e7f134dc0e4c09bd0d9f38099655c02cf0ca57e58823fb21c0.
//
// Solidity: event CollectionSaved(string indexed reqID, address indexed sender, string[] collection)
func (_RandomHashListStorage *RandomHashListStorageFilterer) ParseCollectionSaved(log types.Log) (*RandomHashListStorageCollectionSaved, error) {
	event := new(RandomHashListStorageCollectionSaved)
	if err := _RandomHashListStorage.contract.UnpackLog(event, "CollectionSaved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RandomHashListStorageHashListStoredIterator is returned from FilterHashListStored and is used to iterate over the raw logs and unpacked data for HashListStored events raised by the RandomHashListStorage contract.
type RandomHashListStorageHashListStoredIterator struct {
	Event *RandomHashListStorageHashListStored // Event containing the contract specifics and raw log

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
func (it *RandomHashListStorageHashListStoredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RandomHashListStorageHashListStored)
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
		it.Event = new(RandomHashListStorageHashListStored)
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
func (it *RandomHashListStorageHashListStoredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RandomHashListStorageHashListStoredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RandomHashListStorageHashListStored represents a HashListStored event raised by the RandomHashListStorage contract.
type RandomHashListStorageHashListStored struct {
	ReqID    common.Hash
	HashList []string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterHashListStored is a free log retrieval operation binding the contract event 0x697c9abc51f1b29ebfddcfe0acbf7961bd0f849f3bffffe1e26964879571d8b3.
//
// Solidity: event HashListStored(string indexed reqID, string[] hashList)
func (_RandomHashListStorage *RandomHashListStorageFilterer) FilterHashListStored(opts *bind.FilterOpts, reqID []string) (*RandomHashListStorageHashListStoredIterator, error) {

	var reqIDRule []interface{}
	for _, reqIDItem := range reqID {
		reqIDRule = append(reqIDRule, reqIDItem)
	}

	logs, sub, err := _RandomHashListStorage.contract.FilterLogs(opts, "HashListStored", reqIDRule)
	if err != nil {
		return nil, err
	}
	return &RandomHashListStorageHashListStoredIterator{contract: _RandomHashListStorage.contract, event: "HashListStored", logs: logs, sub: sub}, nil
}

// WatchHashListStored is a free log subscription operation binding the contract event 0x697c9abc51f1b29ebfddcfe0acbf7961bd0f849f3bffffe1e26964879571d8b3.
//
// Solidity: event HashListStored(string indexed reqID, string[] hashList)
func (_RandomHashListStorage *RandomHashListStorageFilterer) WatchHashListStored(opts *bind.WatchOpts, sink chan<- *RandomHashListStorageHashListStored, reqID []string) (event.Subscription, error) {

	var reqIDRule []interface{}
	for _, reqIDItem := range reqID {
		reqIDRule = append(reqIDRule, reqIDItem)
	}

	logs, sub, err := _RandomHashListStorage.contract.WatchLogs(opts, "HashListStored", reqIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RandomHashListStorageHashListStored)
				if err := _RandomHashListStorage.contract.UnpackLog(event, "HashListStored", log); err != nil {
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

// ParseHashListStored is a log parse operation binding the contract event 0x697c9abc51f1b29ebfddcfe0acbf7961bd0f849f3bffffe1e26964879571d8b3.
//
// Solidity: event HashListStored(string indexed reqID, string[] hashList)
func (_RandomHashListStorage *RandomHashListStorageFilterer) ParseHashListStored(log types.Log) (*RandomHashListStorageHashListStored, error) {
	event := new(RandomHashListStorageHashListStored)
	if err := _RandomHashListStorage.contract.UnpackLog(event, "HashListStored", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

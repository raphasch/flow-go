// Code generated by mockery v1.0.0. DO NOT EDIT.

package mempool

import (
	flow "github.com/onflow/flow-go/model/flow"
	mempool "github.com/onflow/flow-go/module/mempool"

	mock "github.com/stretchr/testify/mock"
)

// ReceiptsForest is an autogenerated mock type for the ReceiptsForest type
type ReceiptsForest struct {
	mock.Mock
}

// Add provides a mock function with given fields: receipt, block
func (_m *ReceiptsForest) Add(receipt *flow.ExecutionReceipt, block *flow.Header) (bool, error) {
	ret := _m.Called(receipt, block)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*flow.ExecutionReceipt, *flow.Header) bool); ok {
		r0 = rf(receipt, block)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*flow.ExecutionReceipt, *flow.Header) error); ok {
		r1 = rf(receipt, block)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LowestHeight provides a mock function with given fields:
func (_m *ReceiptsForest) LowestHeight() uint64 {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// PruneUpToHeight provides a mock function with given fields: newLowestHeight
func (_m *ReceiptsForest) PruneUpToHeight(newLowestHeight uint64) error {
	ret := _m.Called(newLowestHeight)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint64) error); ok {
		r0 = rf(newLowestHeight)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReachableReceipts provides a mock function with given fields: resultID, blockFilter, receiptFilter
func (_m *ReceiptsForest) ReachableReceipts(resultID flow.Identifier, blockFilter mempool.BlockFilter, receiptFilter mempool.ReceiptFilter) ([]*flow.ExecutionReceipt, error) {
	ret := _m.Called(resultID, blockFilter, receiptFilter)

	var r0 []*flow.ExecutionReceipt
	if rf, ok := ret.Get(0).(func(flow.Identifier, mempool.BlockFilter, mempool.ReceiptFilter) []*flow.ExecutionReceipt); ok {
		r0 = rf(resultID, blockFilter, receiptFilter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*flow.ExecutionReceipt)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Identifier, mempool.BlockFilter, mempool.ReceiptFilter) error); ok {
		r1 = rf(resultID, blockFilter, receiptFilter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Size provides a mock function with given fields:
func (_m *ReceiptsForest) Size() uint {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}
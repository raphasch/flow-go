// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	protocol "github.com/dapperlabs/flow-go/state/protocol"
)

// ClusterRootQCVoter is an autogenerated mock type for the ClusterRootQCVoter type
type ClusterRootQCVoter struct {
	mock.Mock
}

// Vote provides a mock function with given fields: _a0, _a1
func (_m *ClusterRootQCVoter) Vote(_a0 context.Context, _a1 protocol.Epoch) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, protocol.Epoch) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

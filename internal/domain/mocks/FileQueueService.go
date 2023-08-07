// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"
	domain "pubsub-assignment/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// FileQueueService is an autogenerated mock type for the FileQueueService type
type FileQueueService struct {
	mock.Mock
}

// ReadFile provides a mock function with given fields: ctx
func (_m *FileQueueService) ReadFile(ctx context.Context) (domain.File, error) {
	ret := _m.Called(ctx)

	var r0 domain.File
	if rf, ok := ret.Get(0).(func(context.Context) domain.File); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(domain.File)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WriteFile provides a mock function with given fields: ctx, f
func (_m *FileQueueService) WriteFile(ctx context.Context, f domain.File) error {
	ret := _m.Called(ctx, f)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.File) error); ok {
		r0 = rf(ctx, f)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewFileQueueService interface {
	mock.TestingT
	Cleanup(func())
}

// NewFileQueueService creates a new instance of FileQueueService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFileQueueService(t mockConstructorTestingTNewFileQueueService) *FileQueueService {
	mock := &FileQueueService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
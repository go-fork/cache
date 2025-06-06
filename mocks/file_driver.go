// Code generated by mockery. DO NOT EDIT.

package cache_mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// MockFileDriver is an autogenerated mock type for the FileDriver type
type MockFileDriver struct {
	mock.Mock
}

type MockFileDriver_Expecter struct {
	mock *mock.Mock
}

func (_m *MockFileDriver) EXPECT() *MockFileDriver_Expecter {
	return &MockFileDriver_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with no fields
func (_m *MockFileDriver) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockFileDriver_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type MockFileDriver_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *MockFileDriver_Expecter) Close() *MockFileDriver_Close_Call {
	return &MockFileDriver_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *MockFileDriver_Close_Call) Run(run func()) *MockFileDriver_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockFileDriver_Close_Call) Return(_a0 error) *MockFileDriver_Close_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFileDriver_Close_Call) RunAndReturn(run func() error) *MockFileDriver_Close_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, key
func (_m *MockFileDriver) Delete(ctx context.Context, key string) error {
	ret := _m.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockFileDriver_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockFileDriver_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
func (_e *MockFileDriver_Expecter) Delete(ctx interface{}, key interface{}) *MockFileDriver_Delete_Call {
	return &MockFileDriver_Delete_Call{Call: _e.mock.On("Delete", ctx, key)}
}

func (_c *MockFileDriver_Delete_Call) Run(run func(ctx context.Context, key string)) *MockFileDriver_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockFileDriver_Delete_Call) Return(_a0 error) *MockFileDriver_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFileDriver_Delete_Call) RunAndReturn(run func(context.Context, string) error) *MockFileDriver_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteMultiple provides a mock function with given fields: ctx, keys
func (_m *MockFileDriver) DeleteMultiple(ctx context.Context, keys []string) error {
	ret := _m.Called(ctx, keys)

	if len(ret) == 0 {
		panic("no return value specified for DeleteMultiple")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []string) error); ok {
		r0 = rf(ctx, keys)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockFileDriver_DeleteMultiple_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteMultiple'
type MockFileDriver_DeleteMultiple_Call struct {
	*mock.Call
}

// DeleteMultiple is a helper method to define mock.On call
//   - ctx context.Context
//   - keys []string
func (_e *MockFileDriver_Expecter) DeleteMultiple(ctx interface{}, keys interface{}) *MockFileDriver_DeleteMultiple_Call {
	return &MockFileDriver_DeleteMultiple_Call{Call: _e.mock.On("DeleteMultiple", ctx, keys)}
}

func (_c *MockFileDriver_DeleteMultiple_Call) Run(run func(ctx context.Context, keys []string)) *MockFileDriver_DeleteMultiple_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]string))
	})
	return _c
}

func (_c *MockFileDriver_DeleteMultiple_Call) Return(_a0 error) *MockFileDriver_DeleteMultiple_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFileDriver_DeleteMultiple_Call) RunAndReturn(run func(context.Context, []string) error) *MockFileDriver_DeleteMultiple_Call {
	_c.Call.Return(run)
	return _c
}

// Flush provides a mock function with given fields: ctx
func (_m *MockFileDriver) Flush(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Flush")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockFileDriver_Flush_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Flush'
type MockFileDriver_Flush_Call struct {
	*mock.Call
}

// Flush is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockFileDriver_Expecter) Flush(ctx interface{}) *MockFileDriver_Flush_Call {
	return &MockFileDriver_Flush_Call{Call: _e.mock.On("Flush", ctx)}
}

func (_c *MockFileDriver_Flush_Call) Run(run func(ctx context.Context)) *MockFileDriver_Flush_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockFileDriver_Flush_Call) Return(_a0 error) *MockFileDriver_Flush_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFileDriver_Flush_Call) RunAndReturn(run func(context.Context) error) *MockFileDriver_Flush_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, key
func (_m *MockFileDriver) Get(ctx context.Context, key string) (interface{}, bool) {
	ret := _m.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 interface{}
	var r1 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) (interface{}, bool)); ok {
		return rf(ctx, key)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) interface{}); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) bool); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// MockFileDriver_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockFileDriver_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
func (_e *MockFileDriver_Expecter) Get(ctx interface{}, key interface{}) *MockFileDriver_Get_Call {
	return &MockFileDriver_Get_Call{Call: _e.mock.On("Get", ctx, key)}
}

func (_c *MockFileDriver_Get_Call) Run(run func(ctx context.Context, key string)) *MockFileDriver_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockFileDriver_Get_Call) Return(_a0 interface{}, _a1 bool) *MockFileDriver_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockFileDriver_Get_Call) RunAndReturn(run func(context.Context, string) (interface{}, bool)) *MockFileDriver_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetMultiple provides a mock function with given fields: ctx, keys
func (_m *MockFileDriver) GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, []string) {
	ret := _m.Called(ctx, keys)

	if len(ret) == 0 {
		panic("no return value specified for GetMultiple")
	}

	var r0 map[string]interface{}
	var r1 []string
	if rf, ok := ret.Get(0).(func(context.Context, []string) (map[string]interface{}, []string)); ok {
		return rf(ctx, keys)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []string) map[string]interface{}); ok {
		r0 = rf(ctx, keys)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []string) []string); ok {
		r1 = rf(ctx, keys)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]string)
		}
	}

	return r0, r1
}

// MockFileDriver_GetMultiple_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMultiple'
type MockFileDriver_GetMultiple_Call struct {
	*mock.Call
}

// GetMultiple is a helper method to define mock.On call
//   - ctx context.Context
//   - keys []string
func (_e *MockFileDriver_Expecter) GetMultiple(ctx interface{}, keys interface{}) *MockFileDriver_GetMultiple_Call {
	return &MockFileDriver_GetMultiple_Call{Call: _e.mock.On("GetMultiple", ctx, keys)}
}

func (_c *MockFileDriver_GetMultiple_Call) Run(run func(ctx context.Context, keys []string)) *MockFileDriver_GetMultiple_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]string))
	})
	return _c
}

func (_c *MockFileDriver_GetMultiple_Call) Return(_a0 map[string]interface{}, _a1 []string) *MockFileDriver_GetMultiple_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockFileDriver_GetMultiple_Call) RunAndReturn(run func(context.Context, []string) (map[string]interface{}, []string)) *MockFileDriver_GetMultiple_Call {
	_c.Call.Return(run)
	return _c
}

// Has provides a mock function with given fields: ctx, key
func (_m *MockFileDriver) Has(ctx context.Context, key string) bool {
	ret := _m.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for Has")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockFileDriver_Has_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Has'
type MockFileDriver_Has_Call struct {
	*mock.Call
}

// Has is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
func (_e *MockFileDriver_Expecter) Has(ctx interface{}, key interface{}) *MockFileDriver_Has_Call {
	return &MockFileDriver_Has_Call{Call: _e.mock.On("Has", ctx, key)}
}

func (_c *MockFileDriver_Has_Call) Run(run func(ctx context.Context, key string)) *MockFileDriver_Has_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockFileDriver_Has_Call) Return(_a0 bool) *MockFileDriver_Has_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFileDriver_Has_Call) RunAndReturn(run func(context.Context, string) bool) *MockFileDriver_Has_Call {
	_c.Call.Return(run)
	return _c
}

// Remember provides a mock function with given fields: ctx, key, ttl, callback
func (_m *MockFileDriver) Remember(ctx context.Context, key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	ret := _m.Called(ctx, key, ttl, callback)

	if len(ret) == 0 {
		panic("no return value specified for Remember")
	}

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Duration, func() (interface{}, error)) (interface{}, error)); ok {
		return rf(ctx, key, ttl, callback)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Duration, func() (interface{}, error)) interface{}); ok {
		r0 = rf(ctx, key, ttl, callback)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, time.Duration, func() (interface{}, error)) error); ok {
		r1 = rf(ctx, key, ttl, callback)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockFileDriver_Remember_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Remember'
type MockFileDriver_Remember_Call struct {
	*mock.Call
}

// Remember is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
//   - ttl time.Duration
//   - callback func()(interface{} , error)
func (_e *MockFileDriver_Expecter) Remember(ctx interface{}, key interface{}, ttl interface{}, callback interface{}) *MockFileDriver_Remember_Call {
	return &MockFileDriver_Remember_Call{Call: _e.mock.On("Remember", ctx, key, ttl, callback)}
}

func (_c *MockFileDriver_Remember_Call) Run(run func(ctx context.Context, key string, ttl time.Duration, callback func() (interface{}, error))) *MockFileDriver_Remember_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(time.Duration), args[3].(func() (interface{}, error)))
	})
	return _c
}

func (_c *MockFileDriver_Remember_Call) Return(_a0 interface{}, _a1 error) *MockFileDriver_Remember_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockFileDriver_Remember_Call) RunAndReturn(run func(context.Context, string, time.Duration, func() (interface{}, error)) (interface{}, error)) *MockFileDriver_Remember_Call {
	_c.Call.Return(run)
	return _c
}

// Set provides a mock function with given fields: ctx, key, value, ttl
func (_m *MockFileDriver) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	ret := _m.Called(ctx, key, value, ttl)

	if len(ret) == 0 {
		panic("no return value specified for Set")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}, time.Duration) error); ok {
		r0 = rf(ctx, key, value, ttl)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockFileDriver_Set_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Set'
type MockFileDriver_Set_Call struct {
	*mock.Call
}

// Set is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
//   - value interface{}
//   - ttl time.Duration
func (_e *MockFileDriver_Expecter) Set(ctx interface{}, key interface{}, value interface{}, ttl interface{}) *MockFileDriver_Set_Call {
	return &MockFileDriver_Set_Call{Call: _e.mock.On("Set", ctx, key, value, ttl)}
}

func (_c *MockFileDriver_Set_Call) Run(run func(ctx context.Context, key string, value interface{}, ttl time.Duration)) *MockFileDriver_Set_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(interface{}), args[3].(time.Duration))
	})
	return _c
}

func (_c *MockFileDriver_Set_Call) Return(_a0 error) *MockFileDriver_Set_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFileDriver_Set_Call) RunAndReturn(run func(context.Context, string, interface{}, time.Duration) error) *MockFileDriver_Set_Call {
	_c.Call.Return(run)
	return _c
}

// SetMultiple provides a mock function with given fields: ctx, values, ttl
func (_m *MockFileDriver) SetMultiple(ctx context.Context, values map[string]interface{}, ttl time.Duration) error {
	ret := _m.Called(ctx, values, ttl)

	if len(ret) == 0 {
		panic("no return value specified for SetMultiple")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}, time.Duration) error); ok {
		r0 = rf(ctx, values, ttl)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockFileDriver_SetMultiple_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetMultiple'
type MockFileDriver_SetMultiple_Call struct {
	*mock.Call
}

// SetMultiple is a helper method to define mock.On call
//   - ctx context.Context
//   - values map[string]interface{}
//   - ttl time.Duration
func (_e *MockFileDriver_Expecter) SetMultiple(ctx interface{}, values interface{}, ttl interface{}) *MockFileDriver_SetMultiple_Call {
	return &MockFileDriver_SetMultiple_Call{Call: _e.mock.On("SetMultiple", ctx, values, ttl)}
}

func (_c *MockFileDriver_SetMultiple_Call) Run(run func(ctx context.Context, values map[string]interface{}, ttl time.Duration)) *MockFileDriver_SetMultiple_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(map[string]interface{}), args[2].(time.Duration))
	})
	return _c
}

func (_c *MockFileDriver_SetMultiple_Call) Return(_a0 error) *MockFileDriver_SetMultiple_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFileDriver_SetMultiple_Call) RunAndReturn(run func(context.Context, map[string]interface{}, time.Duration) error) *MockFileDriver_SetMultiple_Call {
	_c.Call.Return(run)
	return _c
}

// Stats provides a mock function with given fields: ctx
func (_m *MockFileDriver) Stats(ctx context.Context) map[string]interface{} {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Stats")
	}

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(context.Context) map[string]interface{}); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	return r0
}

// MockFileDriver_Stats_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stats'
type MockFileDriver_Stats_Call struct {
	*mock.Call
}

// Stats is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockFileDriver_Expecter) Stats(ctx interface{}) *MockFileDriver_Stats_Call {
	return &MockFileDriver_Stats_Call{Call: _e.mock.On("Stats", ctx)}
}

func (_c *MockFileDriver_Stats_Call) Run(run func(ctx context.Context)) *MockFileDriver_Stats_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockFileDriver_Stats_Call) Return(_a0 map[string]interface{}) *MockFileDriver_Stats_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFileDriver_Stats_Call) RunAndReturn(run func(context.Context) map[string]interface{}) *MockFileDriver_Stats_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockFileDriver creates a new instance of MockFileDriver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockFileDriver(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockFileDriver {
	mock := &MockFileDriver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

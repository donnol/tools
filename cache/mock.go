package cache

import "time"

type CacheMock struct {
	GetFunc func(key string) interface{}

	LookupFunc func(key string) (interface{}, bool)

	SetFunc func(key string, value interface{}) bool

	SetNXFunc func(key string, value interface{}, expire time.Duration) bool
}

var _ Cache = &CacheMock{}

func (*CacheMock) Get(key string) interface{} {
	panic("Need to be implement!")
}

func (*CacheMock) Lookup(key string) (interface{}, bool) {
	panic("Need to be implement!")
}

func (*CacheMock) Set(key string, value interface{}) bool {
	panic("Need to be implement!")
}

func (*CacheMock) SetNX(key string, value interface{}, expire time.Duration) bool {
	panic("Need to be implement!")
}

type ImemImplMock struct {
	GetFunc func(key string) interface{}

	LookupFunc func(key string) (interface{}, bool)

	SetFunc func(key string, value interface{}) bool

	SetNXFunc func(key string, value interface{}, expire time.Duration) bool
}

var _ ImemImpl = &ImemImplMock{}

func (*ImemImplMock) Get(key string) interface{} {
	panic("Need to be implement!")
}

func (*ImemImplMock) Lookup(key string) (interface{}, bool) {
	panic("Need to be implement!")
}

func (*ImemImplMock) Set(key string, value interface{}) bool {
	panic("Need to be implement!")
}

func (*ImemImplMock) SetNX(key string, value interface{}, expire time.Duration) bool {
	panic("Need to be implement!")
}

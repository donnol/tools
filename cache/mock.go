package cache

import "time"

type memImplMock struct {
	GetFunc func(key string) any

	LookupFunc func(key string) (any, bool)

	SetFunc func(key string, value any) bool

	SetNXFunc func(key string, value any, expire time.Duration) bool
}

var _ ImemImpl = &memImplMock{}

func (mockRecv *memImplMock) Get(key string) any {
	return mockRecv.GetFunc(key)
}

func (mockRecv *memImplMock) Lookup(key string) (any, bool) {
	return mockRecv.LookupFunc(key)
}

func (mockRecv *memImplMock) Set(key string, value any) bool {
	return mockRecv.SetFunc(key, value)
}

func (mockRecv *memImplMock) SetNX(key string, value any, expire time.Duration) bool {
	return mockRecv.SetNXFunc(key, value, expire)
}

type CacheMockMock struct {
	GetFunc func(key string) any

	LookupFunc func(key string) (any, bool)

	SetFunc func(key string, value any) bool

	SetNXFunc func(key string, value any, expire time.Duration) bool
}

var _ ICacheMock = &CacheMockMock{}

func (mockRecv *CacheMockMock) Get(key string) any {
	return mockRecv.GetFunc(key)
}

func (mockRecv *CacheMockMock) Lookup(key string) (any, bool) {
	return mockRecv.LookupFunc(key)
}

func (mockRecv *CacheMockMock) Set(key string, value any) bool {
	return mockRecv.SetFunc(key, value)
}

func (mockRecv *CacheMockMock) SetNX(key string, value any, expire time.Duration) bool {
	return mockRecv.SetNXFunc(key, value, expire)
}

type memImplMockMock struct {
	GetFunc func(key string) any

	LookupFunc func(key string) (any, bool)

	SetFunc func(key string, value any) bool

	SetNXFunc func(key string, value any, expire time.Duration) bool
}

var _ ImemImplMock = &memImplMockMock{}

func (mockRecv *memImplMockMock) Get(key string) any {
	return mockRecv.GetFunc(key)
}

func (mockRecv *memImplMockMock) Lookup(key string) (any, bool) {
	return mockRecv.LookupFunc(key)
}

func (mockRecv *memImplMockMock) Set(key string, value any) bool {
	return mockRecv.SetFunc(key, value)
}

func (mockRecv *memImplMockMock) SetNX(key string, value any, expire time.Duration) bool {
	return mockRecv.SetNXFunc(key, value, expire)
}

type CacheMock struct {
	GetFunc func(key string) any

	LookupFunc func(key string) (any, bool)

	SetFunc func(key string, value any) bool

	SetNXFunc func(key string, value any, expire time.Duration) bool
}

var _ Cache = &CacheMock{}

func (mockRecv *CacheMock) Get(key string) any {
	return mockRecv.GetFunc(key)
}

func (mockRecv *CacheMock) Lookup(key string) (any, bool) {
	return mockRecv.LookupFunc(key)
}

func (mockRecv *CacheMock) Set(key string, value any) bool {
	return mockRecv.SetFunc(key, value)
}

func (mockRecv *CacheMock) SetNX(key string, value any, expire time.Duration) bool {
	return mockRecv.SetNXFunc(key, value, expire)
}

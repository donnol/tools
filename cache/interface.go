package cache

import "time"

type ImemImpl interface {
	Get(key string) any
	Lookup(key string) (any, bool)
	Set(key string, value any) bool
	SetNX(key string, value any, expire time.Duration) bool
}

type ICacheMock interface {
	Get(key string) any
	Lookup(key string) (any, bool)
	Set(key string, value any) bool
	SetNX(key string, value any, expire time.Duration) bool
}

type ImemImplMock interface {
	Get(key string) any
	Lookup(key string) (any, bool)
	Set(key string, value any) bool
	SetNX(key string, value any, expire time.Duration) bool
}

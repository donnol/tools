package cache

import "time"

type ImemImpl interface {
	Get(key string) interface{}
	Lookup(key string) (interface{}, bool)
	Set(key string, value interface{}) bool
	SetNX(key string, value interface{}, expire time.Duration) bool
}

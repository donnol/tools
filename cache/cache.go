package cache

import (
	"fmt"
	"time"
)

// Cache 缓存
type Cache interface {
	Get(key string) any
	Lookup(key string) (any, bool)
	Set(key string, value any) bool
	SetNX(key string, value any, expire time.Duration) bool
}

// Option 选项
type Option struct {
	Type Type // 类型：1 内存；2 redis；
}

// Type 缓存类型
type Type int

// 缓存类型枚举
const (
	TypeMem Type = iota + 1
	TypeRedis
)

// New 新建
func New(opt Option) Cache {
	switch opt.Type {
	case TypeMem:
		return newMemImpl()
	case TypeRedis:
		// TODO:
	default:
		panic(fmt.Sprintf("Now support type: %d", opt.Type))
	}
	return nil
}

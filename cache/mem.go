package cache

import (
	"log"
	"sync"
	"time"
)

// memImpl 实现
type memImpl struct {
	mutex sync.RWMutex

	// TODO:
	// From https://www.jianshu.com/p/9352d20fb2e0
	// 键删除策略
	// 1. 立即删除 -- 在键设置时绑定回调；占用CPU
	// 2. 惰性删除 -- 键查询的时候查看过期时间，如已过期则删除；占用内存
	// 3. 定时删除 -- 每隔一段时间执行一次删除操作，并通过限制删除操作执行的时长和频率，来减少删除操作对cpu的影响。另一方面定时删除也有效的减少了因惰性删除带来的内存浪费。
	m map[string]memValue

	delDuration time.Duration
}

type memValue struct {
	value       any
	haveExpired bool // 是否设置了过期时间
	expiredAt   time.Time
}

func newMemImpl() *memImpl {
	mi := &memImpl{
		m:           make(map[string]memValue),
		delDuration: time.Second * 1,
	}

	// 开启定时删除
	go mi.beginDelTask()

	return mi
}

var _ Cache = &memImpl{}

// Get 获取
func (i *memImpl) Get(key string) any {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	now := time.Now()
	mi, ok := i.lookup(key, now)
	if ok {
		return mi.value
	}

	return nil
}

// Lookup 寻找
func (i *memImpl) Lookup(key string) (any, bool) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	now := time.Now()
	mi, ok := i.lookup(key, now)
	if ok {
		return mi.value, true
	}

	return nil, false
}

// Set 设置
func (i *memImpl) Set(key string, value any) bool {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	i.m[key] = memValue{
		value: value,
	}

	return true
}

// SetNX 带超时，并且校验是否存在已有key，没有才设置
func (i *memImpl) SetNX(key string, value any, expire time.Duration) bool {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	now := time.Now()
	expiredAt := now.Add(expire)
	i.m[key] = memValue{
		value:       value,
		haveExpired: true,
		expiredAt:   expiredAt,
	}

	return true
}

func (i *memImpl) lookup(key string, now time.Time) (memValue, bool) {
	v, ok := i.m[key]
	// 存在且未过期
	if ok && (!v.haveExpired || now.Before(v.expiredAt)) {
		return v, true
	}
	if ok {
		return v, false
	}

	return memValue{}, false
}

func (i *memImpl) beginDelTask() {
	tickChan := time.Tick(i.delDuration)
	for {
		select {
		case <-tickChan:
			now := time.Now()
			for key, value := range i.m {
				if value.haveExpired && !now.Before(value.expiredAt) {
					delete(i.m, key)
					log.Printf("delete %s from cache\n", key)
				}
			}
		}
	}
}

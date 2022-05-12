package cache

import (
	"fmt"
	"log"
	"time"

	"github.com/donnol/tools/inject"
)

var (
	_gen_customCtxMap = make(map[string]inject.CtxFunc)
)

func RegisterProxyMethod(pctx inject.ProxyContext, cf inject.CtxFunc) {
	_gen_customCtxMap[pctx.Uniq()] = cf
}

type CacheMock struct {
	GetFunc func(key string) any

	LookupFunc func(key string) (any, bool)

	SetFunc func(key string, value any) bool

	SetNXFunc func(key string, value any, expire time.Duration) bool
}

var (
	_ Cache = &CacheMock{}

	cacheMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "github.com/donnol/tools/cache",
		InterfaceName: "Cache",
	}
	CacheMockGetProxyContext = func() (pctx inject.ProxyContext) {
		pctx = cacheMockCommonProxyContext
		pctx.MethodName = "Get"
		return
	}()
	CacheMockLookupProxyContext = func() (pctx inject.ProxyContext) {
		pctx = cacheMockCommonProxyContext
		pctx.MethodName = "Lookup"
		return
	}()
	CacheMockSetProxyContext = func() (pctx inject.ProxyContext) {
		pctx = cacheMockCommonProxyContext
		pctx.MethodName = "Set"
		return
	}()
	CacheMockSetNXProxyContext = func() (pctx inject.ProxyContext) {
		pctx = cacheMockCommonProxyContext
		pctx.MethodName = "SetNX"
		return
	}()

	_ = getCacheProxy
)

func getCacheProxy(base Cache) *CacheMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &CacheMock{
		GetFunc: func(key string) any {
			_gen_begin := time.Now()

			var _gen_r0 any

			_gen_ctx := CacheMockGetProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, key)

				_gen_res := _gen_cf(_gen_ctx, base.Get, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(any)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.Get(key)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		LookupFunc: func(key string) (any, bool) {
			_gen_begin := time.Now()

			var _gen_r0 any

			var _gen_r1 bool

			_gen_ctx := CacheMockLookupProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, key)

				_gen_res := _gen_cf(_gen_ctx, base.Lookup, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(any)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

				_gen_tmpr1, _gen_exist := _gen_res[1].(bool)
				if _gen_exist {
					_gen_r1 = _gen_tmpr1
				}

			} else {
				_gen_r0, _gen_r1 = base.Lookup(key)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0, _gen_r1
		},

		SetFunc: func(key string, value any) bool {
			_gen_begin := time.Now()

			var _gen_r0 bool

			_gen_ctx := CacheMockSetProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, key)

				_gen_params = append(_gen_params, value)

				_gen_res := _gen_cf(_gen_ctx, base.Set, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(bool)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.Set(key, value)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		SetNXFunc: func(key string, value any, expire time.Duration) bool {
			_gen_begin := time.Now()

			var _gen_r0 bool

			_gen_ctx := CacheMockSetNXProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, key)

				_gen_params = append(_gen_params, value)

				_gen_params = append(_gen_params, expire)

				_gen_res := _gen_cf(_gen_ctx, base.SetNX, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(bool)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.SetNX(key, value, expire)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},
	}
}

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

type CacheMockMock struct {
	GetFunc func(key string) any

	LookupFunc func(key string) (any, bool)

	SetFunc func(key string, value any) bool

	SetNXFunc func(key string, value any, expire time.Duration) bool
}

var (
	_ ICacheMock = &CacheMockMock{}

	cacheMockMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "github.com/donnol/tools/cache",
		InterfaceName: "ICacheMock",
	}
	CacheMockMockGetProxyContext = func() (pctx inject.ProxyContext) {
		pctx = cacheMockMockCommonProxyContext
		pctx.MethodName = "Get"
		return
	}()
	CacheMockMockLookupProxyContext = func() (pctx inject.ProxyContext) {
		pctx = cacheMockMockCommonProxyContext
		pctx.MethodName = "Lookup"
		return
	}()
	CacheMockMockSetProxyContext = func() (pctx inject.ProxyContext) {
		pctx = cacheMockMockCommonProxyContext
		pctx.MethodName = "Set"
		return
	}()
	CacheMockMockSetNXProxyContext = func() (pctx inject.ProxyContext) {
		pctx = cacheMockMockCommonProxyContext
		pctx.MethodName = "SetNX"
		return
	}()

	_ = getICacheMockProxy
)

func getICacheMockProxy(base ICacheMock) *CacheMockMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &CacheMockMock{
		GetFunc: func(key string) any {
			_gen_begin := time.Now()

			var _gen_r0 any

			_gen_ctx := CacheMockMockGetProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, key)

				_gen_res := _gen_cf(_gen_ctx, base.Get, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(any)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.Get(key)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		LookupFunc: func(key string) (any, bool) {
			_gen_begin := time.Now()

			var _gen_r0 any

			var _gen_r1 bool

			_gen_ctx := CacheMockMockLookupProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, key)

				_gen_res := _gen_cf(_gen_ctx, base.Lookup, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(any)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

				_gen_tmpr1, _gen_exist := _gen_res[1].(bool)
				if _gen_exist {
					_gen_r1 = _gen_tmpr1
				}

			} else {
				_gen_r0, _gen_r1 = base.Lookup(key)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0, _gen_r1
		},

		SetFunc: func(key string, value any) bool {
			_gen_begin := time.Now()

			var _gen_r0 bool

			_gen_ctx := CacheMockMockSetProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, key)

				_gen_params = append(_gen_params, value)

				_gen_res := _gen_cf(_gen_ctx, base.Set, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(bool)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.Set(key, value)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		SetNXFunc: func(key string, value any, expire time.Duration) bool {
			_gen_begin := time.Now()

			var _gen_r0 bool

			_gen_ctx := CacheMockMockSetNXProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, key)

				_gen_params = append(_gen_params, value)

				_gen_params = append(_gen_params, expire)

				_gen_res := _gen_cf(_gen_ctx, base.SetNX, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(bool)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.SetNX(key, value, expire)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},
	}
}

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

type memImplMock struct {
	GetFunc func(key string) any

	LookupFunc func(key string) (any, bool)

	SetFunc func(key string, value any) bool

	SetNXFunc func(key string, value any, expire time.Duration) bool
}

var (
	_ ImemImpl = &memImplMock{}

	memImplMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "github.com/donnol/tools/cache",
		InterfaceName: "ImemImpl",
	}
	memImplMockGetProxyContext = func() (pctx inject.ProxyContext) {
		pctx = memImplMockCommonProxyContext
		pctx.MethodName = "Get"
		return
	}()
	memImplMockLookupProxyContext = func() (pctx inject.ProxyContext) {
		pctx = memImplMockCommonProxyContext
		pctx.MethodName = "Lookup"
		return
	}()
	memImplMockSetProxyContext = func() (pctx inject.ProxyContext) {
		pctx = memImplMockCommonProxyContext
		pctx.MethodName = "Set"
		return
	}()
	memImplMockSetNXProxyContext = func() (pctx inject.ProxyContext) {
		pctx = memImplMockCommonProxyContext
		pctx.MethodName = "SetNX"
		return
	}()

	_ = getImemImplProxy
)

func getImemImplProxy(base ImemImpl) *memImplMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &memImplMock{
		GetFunc: func(key string) any {
			_gen_begin := time.Now()

			var _gen_r0 any

			_gen_ctx := memImplMockGetProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, key)

				_gen_res := _gen_cf(_gen_ctx, base.Get, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(any)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.Get(key)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		LookupFunc: func(key string) (any, bool) {
			_gen_begin := time.Now()

			var _gen_r0 any

			var _gen_r1 bool

			_gen_ctx := memImplMockLookupProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, key)

				_gen_res := _gen_cf(_gen_ctx, base.Lookup, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(any)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

				_gen_tmpr1, _gen_exist := _gen_res[1].(bool)
				if _gen_exist {
					_gen_r1 = _gen_tmpr1
				}

			} else {
				_gen_r0, _gen_r1 = base.Lookup(key)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0, _gen_r1
		},

		SetFunc: func(key string, value any) bool {
			_gen_begin := time.Now()

			var _gen_r0 bool

			_gen_ctx := memImplMockSetProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, key)

				_gen_params = append(_gen_params, value)

				_gen_res := _gen_cf(_gen_ctx, base.Set, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(bool)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.Set(key, value)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		SetNXFunc: func(key string, value any, expire time.Duration) bool {
			_gen_begin := time.Now()

			var _gen_r0 bool

			_gen_ctx := memImplMockSetNXProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, key)

				_gen_params = append(_gen_params, value)

				_gen_params = append(_gen_params, expire)

				_gen_res := _gen_cf(_gen_ctx, base.SetNX, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(bool)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.SetNX(key, value, expire)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},
	}
}

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

type memImplMockMock struct {
	GetFunc func(key string) any

	LookupFunc func(key string) (any, bool)

	SetFunc func(key string, value any) bool

	SetNXFunc func(key string, value any, expire time.Duration) bool
}

var (
	_ ImemImplMock = &memImplMockMock{}

	memImplMockMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "github.com/donnol/tools/cache",
		InterfaceName: "ImemImplMock",
	}
	memImplMockMockGetProxyContext = func() (pctx inject.ProxyContext) {
		pctx = memImplMockMockCommonProxyContext
		pctx.MethodName = "Get"
		return
	}()
	memImplMockMockLookupProxyContext = func() (pctx inject.ProxyContext) {
		pctx = memImplMockMockCommonProxyContext
		pctx.MethodName = "Lookup"
		return
	}()
	memImplMockMockSetProxyContext = func() (pctx inject.ProxyContext) {
		pctx = memImplMockMockCommonProxyContext
		pctx.MethodName = "Set"
		return
	}()
	memImplMockMockSetNXProxyContext = func() (pctx inject.ProxyContext) {
		pctx = memImplMockMockCommonProxyContext
		pctx.MethodName = "SetNX"
		return
	}()

	_ = getImemImplMockProxy
)

func getImemImplMockProxy(base ImemImplMock) *memImplMockMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &memImplMockMock{
		GetFunc: func(key string) any {
			_gen_begin := time.Now()

			var _gen_r0 any

			_gen_ctx := memImplMockMockGetProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, key)

				_gen_res := _gen_cf(_gen_ctx, base.Get, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(any)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.Get(key)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		LookupFunc: func(key string) (any, bool) {
			_gen_begin := time.Now()

			var _gen_r0 any

			var _gen_r1 bool

			_gen_ctx := memImplMockMockLookupProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, key)

				_gen_res := _gen_cf(_gen_ctx, base.Lookup, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(any)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

				_gen_tmpr1, _gen_exist := _gen_res[1].(bool)
				if _gen_exist {
					_gen_r1 = _gen_tmpr1
				}

			} else {
				_gen_r0, _gen_r1 = base.Lookup(key)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0, _gen_r1
		},

		SetFunc: func(key string, value any) bool {
			_gen_begin := time.Now()

			var _gen_r0 bool

			_gen_ctx := memImplMockMockSetProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, key)

				_gen_params = append(_gen_params, value)

				_gen_res := _gen_cf(_gen_ctx, base.Set, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(bool)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.Set(key, value)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		SetNXFunc: func(key string, value any, expire time.Duration) bool {
			_gen_begin := time.Now()

			var _gen_r0 bool

			_gen_ctx := memImplMockMockSetNXProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, key)

				_gen_params = append(_gen_params, value)

				_gen_params = append(_gen_params, expire)

				_gen_res := _gen_cf(_gen_ctx, base.SetNX, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(bool)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.SetNX(key, value, expire)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},
	}
}

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

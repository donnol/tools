package list

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

type StringListMock struct {
	FilterFunc func(s string) StringList
}

var (
	_ IStringList = &StringListMock{}

	stringListMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "github.com/donnol/tools/list",
		InterfaceName: "IStringList",
	}
	StringListMockFilterProxyContext = func() (pctx inject.ProxyContext) {
		pctx = stringListMockCommonProxyContext
		pctx.MethodName = "Filter"
		return
	}()

	_ = getIStringListProxy
)

func getIStringListProxy(base IStringList) *StringListMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &StringListMock{
		FilterFunc: func(s string) StringList {
			_gen_begin := time.Now()

			var _gen_r0 StringList

			_gen_ctx := StringListMockFilterProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, s)

				_gen_res := _gen_cf(_gen_ctx, base.Filter, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(StringList)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.Filter(s)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},
	}
}

func (mockRecv *StringListMock) Filter(s string) StringList {
	return mockRecv.FilterFunc(s)
}

type StringListMockMock struct {
	FilterFunc func(s string) StringList
}

var (
	_ IStringListMock = &StringListMockMock{}

	stringListMockMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "github.com/donnol/tools/list",
		InterfaceName: "IStringListMock",
	}
	StringListMockMockFilterProxyContext = func() (pctx inject.ProxyContext) {
		pctx = stringListMockMockCommonProxyContext
		pctx.MethodName = "Filter"
		return
	}()

	_ = getIStringListMockProxy
)

func getIStringListMockProxy(base IStringListMock) *StringListMockMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &StringListMockMock{
		FilterFunc: func(s string) StringList {
			_gen_begin := time.Now()

			var _gen_r0 StringList

			_gen_ctx := StringListMockMockFilterProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, s)

				_gen_res := _gen_cf(_gen_ctx, base.Filter, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(StringList)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.Filter(s)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},
	}
}

func (mockRecv *StringListMockMock) Filter(s string) StringList {
	return mockRecv.FilterFunc(s)
}

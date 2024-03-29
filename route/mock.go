package route

import (
	"context"
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

type CheckerMock struct {
	CheckFunc func(context.Context) error
}

var (
	_ Checker = &CheckerMock{}

	checkerMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "github.com/donnol/tools/route",
		InterfaceName: "Checker",
	}
	CheckerMockCheckProxyContext = func() (pctx inject.ProxyContext) {
		pctx = checkerMockCommonProxyContext
		pctx.MethodName = "Check"
		return
	}()

	_ = getCheckerProxy
)

func getCheckerProxy(base Checker) *CheckerMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &CheckerMock{
		CheckFunc: func(p0 context.Context) error {
			_gen_begin := time.Now()

			var _gen_r0 error

			_gen_ctx := CheckerMockCheckProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, p0)

				_gen_res := _gen_cf(_gen_ctx, base.Check, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(error)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.Check(p0)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},
	}
}

func (mockRecv *CheckerMock) Check(p0 context.Context) error {
	return mockRecv.CheckFunc(p0)
}

type FilterMock struct {
	FilterFunc func() any
}

var (
	_ Filter = &FilterMock{}

	filterMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "github.com/donnol/tools/route",
		InterfaceName: "Filter",
	}
	FilterMockFilterProxyContext = func() (pctx inject.ProxyContext) {
		pctx = filterMockCommonProxyContext
		pctx.MethodName = "Filter"
		return
	}()

	_ = getFilterProxy
)

func getFilterProxy(base Filter) *FilterMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &FilterMock{
		FilterFunc: func() any {
			_gen_begin := time.Now()

			var _gen_r0 any

			_gen_ctx := FilterMockFilterProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_res := _gen_cf(_gen_ctx, base.Filter, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(any)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.Filter()
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},
	}
}

func (mockRecv *FilterMock) Filter() any {
	return mockRecv.FilterFunc()
}

type NewerMock struct {
	NewFunc func() any
}

var (
	_ Newer = &NewerMock{}

	newerMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "github.com/donnol/tools/route",
		InterfaceName: "Newer",
	}
	NewerMockNewProxyContext = func() (pctx inject.ProxyContext) {
		pctx = newerMockCommonProxyContext
		pctx.MethodName = "New"
		return
	}()

	_ = getNewerProxy
)

func getNewerProxy(base Newer) *NewerMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &NewerMock{
		NewFunc: func() any {
			_gen_begin := time.Now()

			var _gen_r0 any

			_gen_ctx := NewerMockNewProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_res := _gen_cf(_gen_ctx, base.New, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(any)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.New()
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},
	}
}

func (mockRecv *NewerMock) New() any {
	return mockRecv.NewFunc()
}

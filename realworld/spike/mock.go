package spike

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

type SpikeMock struct {
	DetailFunc func() (Content, error)

	GrabFunc func(ctx context.Context) (uint, error)
}

var (
	_ Spike = &SpikeMock{}

	spikeMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "github.com/donnol/tools/realworld/spike",
		InterfaceName: "Spike",
	}
	SpikeMockDetailProxyContext = func() (pctx inject.ProxyContext) {
		pctx = spikeMockCommonProxyContext
		pctx.MethodName = "Detail"
		return
	}()
	SpikeMockGrabProxyContext = func() (pctx inject.ProxyContext) {
		pctx = spikeMockCommonProxyContext
		pctx.MethodName = "Grab"
		return
	}()

	_ = getSpikeProxy
)

func getSpikeProxy(base Spike) *SpikeMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &SpikeMock{
		DetailFunc: func() (Content, error) {
			_gen_begin := time.Now()

			var _gen_r0 Content

			var _gen_r1 error

			_gen_ctx := SpikeMockDetailProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_res := _gen_cf(_gen_ctx, base.Detail, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(Content)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

				_gen_tmpr1, _gen_exist := _gen_res[1].(error)
				if _gen_exist {
					_gen_r1 = _gen_tmpr1
				}

			} else {
				_gen_r0, _gen_r1 = base.Detail()
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0, _gen_r1
		},

		GrabFunc: func(ctx context.Context) (uint, error) {
			_gen_begin := time.Now()

			var _gen_r0 uint

			var _gen_r1 error

			_gen_ctx := SpikeMockGrabProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, ctx)

				_gen_res := _gen_cf(_gen_ctx, base.Grab, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(uint)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

				_gen_tmpr1, _gen_exist := _gen_res[1].(error)
				if _gen_exist {
					_gen_r1 = _gen_tmpr1
				}

			} else {
				_gen_r0, _gen_r1 = base.Grab(ctx)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0, _gen_r1
		},
	}
}

func (mockRecv *SpikeMock) Detail() (Content, error) {
	return mockRecv.DetailFunc()
}

func (mockRecv *SpikeMock) Grab(ctx context.Context) (uint, error) {
	return mockRecv.GrabFunc(ctx)
}

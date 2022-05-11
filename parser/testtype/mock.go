package testtype

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/donnol/tools/inject"
	"github.com/donnol/tools/parser/testtype/a"
	"github.com/donnol/tools/parser/testtype/big"
)

type Mock struct {
	HeadFunc func(ctx context.Context, p a.Pa) (big.Int, error)
}

var (
	_ I = &Mock{}

	mockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "github.com/donnol/tools/parser/testtype",
		InterfaceName: "I",
	}
	MockHeadProxyContext = func() (pctx inject.ProxyContext) {
		pctx = mockCommonProxyContext
		pctx.MethodName = "Head"
		return
	}()

	_                 = getIProxy
	_gen_customCtxMap = make(map[string]inject.CtxFunc)
)

func RegisterProxyMethod(pctx inject.ProxyContext, cf inject.CtxFunc) {
	_gen_customCtxMap[pctx.Uniq()] = cf
}

func getIProxy(base I) *Mock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &Mock{
		HeadFunc: func(ctx context.Context, p a.Pa) (big.Int, error) {
			_gen_begin := time.Now()

			var _gen_r0 big.Int

			var _gen_r1 error

			_gen_ctx := MockHeadProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, ctx)

				_gen_params = append(_gen_params, p)

				_gen_res := _gen_cf(_gen_ctx, base.Head, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(big.Int)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

				_gen_tmpr1, _gen_exist := _gen_res[1].(error)
				if _gen_exist {
					_gen_r1 = _gen_tmpr1
				}

			} else {
				_gen_r0, _gen_r1 = base.Head(ctx, p)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0, _gen_r1
		},
	}
}

func (mockRecv *Mock) Head(ctx context.Context, p a.Pa) (big.Int, error) {
	return mockRecv.HeadFunc(ctx, p)
}

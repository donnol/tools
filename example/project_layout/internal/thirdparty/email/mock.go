package email

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

type EmailMock struct {
	SendFunc func(ctx context.Context, msg Message) error
}

var (
	_ IEmail = &EmailMock{}

	emailMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "project_layout/internal/thirdparty/email",
		InterfaceName: "IEmail",
	}
	EmailMockSendProxyContext = func() (pctx inject.ProxyContext) {
		pctx = emailMockCommonProxyContext
		pctx.MethodName = "Send"
		return
	}()

	_ = getIEmailProxy
)

func getIEmailProxy(base IEmail) *EmailMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &EmailMock{
		SendFunc: func(ctx context.Context, msg Message) error {
			_gen_begin := time.Now()

			var _gen_r0 error

			_gen_ctx := EmailMockSendProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, ctx)

				_gen_params = append(_gen_params, msg)

				_gen_res := _gen_cf(_gen_ctx, base.Send, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(error)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.Send(ctx, msg)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},
	}
}

func (mockRecv *EmailMock) Send(ctx context.Context, msg Message) error {
	return mockRecv.SendFunc(ctx, msg)
}

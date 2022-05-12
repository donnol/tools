package service

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

type PingSrvMock struct {
	PingFunc func() string
}

var (
	_ PingSrv = &PingSrvMock{}

	pingSrvMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "project_layout/internal/service",
		InterfaceName: "PingSrv",
	}
	PingSrvMockPingProxyContext = func() (pctx inject.ProxyContext) {
		pctx = pingSrvMockCommonProxyContext
		pctx.MethodName = "Ping"
		return
	}()

	_ = getPingSrvProxy
)

func getPingSrvProxy(base PingSrv) *PingSrvMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &PingSrvMock{
		PingFunc: func() string {
			_gen_begin := time.Now()

			var _gen_r0 string

			_gen_ctx := PingSrvMockPingProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_res := _gen_cf(_gen_ctx, base.Ping, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(string)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.Ping()
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},
	}
}

func (mockRecv *PingSrvMock) Ping() string {
	return mockRecv.PingFunc()
}

type UserSrvMock struct {
	ModNameFunc func(id uint, name string) error
}

var (
	_ UserSrv = &UserSrvMock{}

	userSrvMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "project_layout/internal/service",
		InterfaceName: "UserSrv",
	}
	UserSrvMockModNameProxyContext = func() (pctx inject.ProxyContext) {
		pctx = userSrvMockCommonProxyContext
		pctx.MethodName = "ModName"
		return
	}()

	_ = getUserSrvProxy
)

func getUserSrvProxy(base UserSrv) *UserSrvMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &UserSrvMock{
		ModNameFunc: func(id uint, name string) error {
			_gen_begin := time.Now()

			var _gen_r0 error

			_gen_ctx := UserSrvMockModNameProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, id)

				_gen_params = append(_gen_params, name)

				_gen_res := _gen_cf(_gen_ctx, base.ModName, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(error)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.ModName(id, name)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},
	}
}

func (mockRecv *UserSrvMock) ModName(id uint, name string) error {
	return mockRecv.ModNameFunc(id, name)
}

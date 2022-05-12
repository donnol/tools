package store

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

type UserStoreMock struct {
	ModNameFunc func(id uint, name string) error
}

var (
	_ UserStore = &UserStoreMock{}

	userStoreMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "project_layout/internal/store",
		InterfaceName: "UserStore",
	}
	UserStoreMockModNameProxyContext = func() (pctx inject.ProxyContext) {
		pctx = userStoreMockCommonProxyContext
		pctx.MethodName = "ModName"
		return
	}()

	_ = getUserStoreProxy
)

func getUserStoreProxy(base UserStore) *UserStoreMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &UserStoreMock{
		ModNameFunc: func(id uint, name string) error {
			_gen_begin := time.Now()

			var _gen_r0 error

			_gen_ctx := UserStoreMockModNameProxyContext
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

func (mockRecv *UserStoreMock) ModName(id uint, name string) error {
	return mockRecv.ModNameFunc(id, name)
}

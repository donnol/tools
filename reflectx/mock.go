package reflectx

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

type StructMock struct {
	GetFieldsFunc func() []Field
}

var (
	_ IStruct = &StructMock{}

	structMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "github.com/donnol/tools/reflectx",
		InterfaceName: "IStruct",
	}
	StructMockGetFieldsProxyContext = func() (pctx inject.ProxyContext) {
		pctx = structMockCommonProxyContext
		pctx.MethodName = "GetFields"
		return
	}()

	_ = getIStructProxy
)

func getIStructProxy(base IStruct) *StructMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &StructMock{
		GetFieldsFunc: func() []Field {
			_gen_begin := time.Now()

			var _gen_r0 []Field

			_gen_ctx := StructMockGetFieldsProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_res := _gen_cf(_gen_ctx, base.GetFields, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].([]Field)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.GetFields()
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},
	}
}

func (mockRecv *StructMock) GetFields() []Field {
	return mockRecv.GetFieldsFunc()
}

type StructMockMock struct {
	GetFieldsFunc func() []Field
}

var (
	_ IStructMock = &StructMockMock{}

	structMockMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "github.com/donnol/tools/reflectx",
		InterfaceName: "IStructMock",
	}
	StructMockMockGetFieldsProxyContext = func() (pctx inject.ProxyContext) {
		pctx = structMockMockCommonProxyContext
		pctx.MethodName = "GetFields"
		return
	}()

	_ = getIStructMockProxy
)

func getIStructMockProxy(base IStructMock) *StructMockMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &StructMockMock{
		GetFieldsFunc: func() []Field {
			_gen_begin := time.Now()

			var _gen_r0 []Field

			_gen_ctx := StructMockMockGetFieldsProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_res := _gen_cf(_gen_ctx, base.GetFields, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].([]Field)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.GetFields()
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},
	}
}

func (mockRecv *StructMockMock) GetFields() []Field {
	return mockRecv.GetFieldsFunc()
}

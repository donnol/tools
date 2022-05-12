package dbdoc

import (
	"fmt"
	"io"
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

type TableMock struct {
	MakeGraphFunc func() *Table

	NewFunc func() *Table

	ResolveFunc func(v any) *Table

	SetCommentFunc func(comment string) *Table

	SetDescriptionFunc func(description string) *Table

	SetMapperFunc func(f Mapper) *Table

	SetTypeMapperFunc func(f Mapper) *Table

	WriteFunc func(w io.Writer) *Table
}

var (
	_ ITable = &TableMock{}

	tableMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "github.com/donnol/tools/dbdoc",
		InterfaceName: "ITable",
	}
	TableMockMakeGraphProxyContext = func() (pctx inject.ProxyContext) {
		pctx = tableMockCommonProxyContext
		pctx.MethodName = "MakeGraph"
		return
	}()
	TableMockNewProxyContext = func() (pctx inject.ProxyContext) {
		pctx = tableMockCommonProxyContext
		pctx.MethodName = "New"
		return
	}()
	TableMockResolveProxyContext = func() (pctx inject.ProxyContext) {
		pctx = tableMockCommonProxyContext
		pctx.MethodName = "Resolve"
		return
	}()
	TableMockSetCommentProxyContext = func() (pctx inject.ProxyContext) {
		pctx = tableMockCommonProxyContext
		pctx.MethodName = "SetComment"
		return
	}()
	TableMockSetDescriptionProxyContext = func() (pctx inject.ProxyContext) {
		pctx = tableMockCommonProxyContext
		pctx.MethodName = "SetDescription"
		return
	}()
	TableMockSetMapperProxyContext = func() (pctx inject.ProxyContext) {
		pctx = tableMockCommonProxyContext
		pctx.MethodName = "SetMapper"
		return
	}()
	TableMockSetTypeMapperProxyContext = func() (pctx inject.ProxyContext) {
		pctx = tableMockCommonProxyContext
		pctx.MethodName = "SetTypeMapper"
		return
	}()
	TableMockWriteProxyContext = func() (pctx inject.ProxyContext) {
		pctx = tableMockCommonProxyContext
		pctx.MethodName = "Write"
		return
	}()

	_ = getITableProxy
)

func getITableProxy(base ITable) *TableMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &TableMock{
		MakeGraphFunc: func() *Table {
			_gen_begin := time.Now()

			var _gen_r0 *Table

			_gen_ctx := TableMockMakeGraphProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_res := _gen_cf(_gen_ctx, base.MakeGraph, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(*Table)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.MakeGraph()
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		NewFunc: func() *Table {
			_gen_begin := time.Now()

			var _gen_r0 *Table

			_gen_ctx := TableMockNewProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_res := _gen_cf(_gen_ctx, base.New, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(*Table)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.New()
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		ResolveFunc: func(v any) *Table {
			_gen_begin := time.Now()

			var _gen_r0 *Table

			_gen_ctx := TableMockResolveProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, v)

				_gen_res := _gen_cf(_gen_ctx, base.Resolve, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(*Table)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.Resolve(v)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		SetCommentFunc: func(comment string) *Table {
			_gen_begin := time.Now()

			var _gen_r0 *Table

			_gen_ctx := TableMockSetCommentProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, comment)

				_gen_res := _gen_cf(_gen_ctx, base.SetComment, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(*Table)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.SetComment(comment)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		SetDescriptionFunc: func(description string) *Table {
			_gen_begin := time.Now()

			var _gen_r0 *Table

			_gen_ctx := TableMockSetDescriptionProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, description)

				_gen_res := _gen_cf(_gen_ctx, base.SetDescription, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(*Table)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.SetDescription(description)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		SetMapperFunc: func(f Mapper) *Table {
			_gen_begin := time.Now()

			var _gen_r0 *Table

			_gen_ctx := TableMockSetMapperProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, f)

				_gen_res := _gen_cf(_gen_ctx, base.SetMapper, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(*Table)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.SetMapper(f)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		SetTypeMapperFunc: func(f Mapper) *Table {
			_gen_begin := time.Now()

			var _gen_r0 *Table

			_gen_ctx := TableMockSetTypeMapperProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, f)

				_gen_res := _gen_cf(_gen_ctx, base.SetTypeMapper, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(*Table)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.SetTypeMapper(f)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		WriteFunc: func(w io.Writer) *Table {
			_gen_begin := time.Now()

			var _gen_r0 *Table

			_gen_ctx := TableMockWriteProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, w)

				_gen_res := _gen_cf(_gen_ctx, base.Write, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(*Table)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.Write(w)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},
	}
}

func (mockRecv *TableMock) MakeGraph() *Table {
	return mockRecv.MakeGraphFunc()
}

func (mockRecv *TableMock) New() *Table {
	return mockRecv.NewFunc()
}

func (mockRecv *TableMock) Resolve(v any) *Table {
	return mockRecv.ResolveFunc(v)
}

func (mockRecv *TableMock) SetComment(comment string) *Table {
	return mockRecv.SetCommentFunc(comment)
}

func (mockRecv *TableMock) SetDescription(description string) *Table {
	return mockRecv.SetDescriptionFunc(description)
}

func (mockRecv *TableMock) SetMapper(f Mapper) *Table {
	return mockRecv.SetMapperFunc(f)
}

func (mockRecv *TableMock) SetTypeMapper(f Mapper) *Table {
	return mockRecv.SetTypeMapperFunc(f)
}

func (mockRecv *TableMock) Write(w io.Writer) *Table {
	return mockRecv.WriteFunc(w)
}

package importpath

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

type ImportPathMock struct {
	GetByCurrentDirFunc func() (path string, err error)

	GetCurrentDirModFilePathFunc func() (modDir string, modPath string, err error)

	GetModFilePathFunc func(dir string) (modDir string, modPath string, err error)

	SplitImportPathWithTypeFunc func(importPathWithType string) (string, string)
}

var (
	_ IImportPath = &ImportPathMock{}

	importPathMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "github.com/donnol/tools/importpath",
		InterfaceName: "IImportPath",
	}
	ImportPathMockGetByCurrentDirProxyContext = func() (pctx inject.ProxyContext) {
		pctx = importPathMockCommonProxyContext
		pctx.MethodName = "GetByCurrentDir"
		return
	}()
	ImportPathMockGetCurrentDirModFilePathProxyContext = func() (pctx inject.ProxyContext) {
		pctx = importPathMockCommonProxyContext
		pctx.MethodName = "GetCurrentDirModFilePath"
		return
	}()
	ImportPathMockGetModFilePathProxyContext = func() (pctx inject.ProxyContext) {
		pctx = importPathMockCommonProxyContext
		pctx.MethodName = "GetModFilePath"
		return
	}()
	ImportPathMockSplitImportPathWithTypeProxyContext = func() (pctx inject.ProxyContext) {
		pctx = importPathMockCommonProxyContext
		pctx.MethodName = "SplitImportPathWithType"
		return
	}()

	_ = getIImportPathProxy
)

func getIImportPathProxy(base IImportPath) *ImportPathMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &ImportPathMock{
		GetByCurrentDirFunc: func() (path string, err error) {
			_gen_begin := time.Now()

			var _gen_r0 string

			var _gen_r1 error

			_gen_ctx := ImportPathMockGetByCurrentDirProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_res := _gen_cf(_gen_ctx, base.GetByCurrentDir, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(string)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

				_gen_tmpr1, _gen_exist := _gen_res[1].(error)
				if _gen_exist {
					_gen_r1 = _gen_tmpr1
				}

			} else {
				_gen_r0, _gen_r1 = base.GetByCurrentDir()
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0, _gen_r1
		},

		GetCurrentDirModFilePathFunc: func() (modDir string, modPath string, err error) {
			_gen_begin := time.Now()

			var _gen_r0 string

			var _gen_r1 string

			var _gen_r2 error

			_gen_ctx := ImportPathMockGetCurrentDirModFilePathProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_res := _gen_cf(_gen_ctx, base.GetCurrentDirModFilePath, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(string)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

				_gen_tmpr1, _gen_exist := _gen_res[1].(string)
				if _gen_exist {
					_gen_r1 = _gen_tmpr1
				}

				_gen_tmpr2, _gen_exist := _gen_res[2].(error)
				if _gen_exist {
					_gen_r2 = _gen_tmpr2
				}

			} else {
				_gen_r0, _gen_r1, _gen_r2 = base.GetCurrentDirModFilePath()
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0, _gen_r1, _gen_r2
		},

		GetModFilePathFunc: func(dir string) (modDir string, modPath string, err error) {
			_gen_begin := time.Now()

			var _gen_r0 string

			var _gen_r1 string

			var _gen_r2 error

			_gen_ctx := ImportPathMockGetModFilePathProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, dir)

				_gen_res := _gen_cf(_gen_ctx, base.GetModFilePath, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(string)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

				_gen_tmpr1, _gen_exist := _gen_res[1].(string)
				if _gen_exist {
					_gen_r1 = _gen_tmpr1
				}

				_gen_tmpr2, _gen_exist := _gen_res[2].(error)
				if _gen_exist {
					_gen_r2 = _gen_tmpr2
				}

			} else {
				_gen_r0, _gen_r1, _gen_r2 = base.GetModFilePath(dir)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0, _gen_r1, _gen_r2
		},

		SplitImportPathWithTypeFunc: func(importPathWithType string) (string, string) {
			_gen_begin := time.Now()

			var _gen_r0 string

			var _gen_r1 string

			_gen_ctx := ImportPathMockSplitImportPathWithTypeProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, importPathWithType)

				_gen_res := _gen_cf(_gen_ctx, base.SplitImportPathWithType, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(string)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

				_gen_tmpr1, _gen_exist := _gen_res[1].(string)
				if _gen_exist {
					_gen_r1 = _gen_tmpr1
				}

			} else {
				_gen_r0, _gen_r1 = base.SplitImportPathWithType(importPathWithType)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0, _gen_r1
		},
	}
}

func (mockRecv *ImportPathMock) GetByCurrentDir() (path string, err error) {
	return mockRecv.GetByCurrentDirFunc()
}

func (mockRecv *ImportPathMock) GetCurrentDirModFilePath() (modDir string, modPath string, err error) {
	return mockRecv.GetCurrentDirModFilePathFunc()
}

func (mockRecv *ImportPathMock) GetModFilePath(dir string) (modDir string, modPath string, err error) {
	return mockRecv.GetModFilePathFunc(dir)
}

func (mockRecv *ImportPathMock) SplitImportPathWithType(importPathWithType string) (string, string) {
	return mockRecv.SplitImportPathWithTypeFunc(importPathWithType)
}

type ImportPathMockMock struct {
	GetByCurrentDirFunc func() (path string, err error)

	GetCurrentDirModFilePathFunc func() (modDir string, modPath string, err error)

	GetModFilePathFunc func(dir string) (modDir string, modPath string, err error)
}

var (
	_ IImportPathMock = &ImportPathMockMock{}

	importPathMockMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "github.com/donnol/tools/importpath",
		InterfaceName: "IImportPathMock",
	}
	ImportPathMockMockGetByCurrentDirProxyContext = func() (pctx inject.ProxyContext) {
		pctx = importPathMockMockCommonProxyContext
		pctx.MethodName = "GetByCurrentDir"
		return
	}()
	ImportPathMockMockGetCurrentDirModFilePathProxyContext = func() (pctx inject.ProxyContext) {
		pctx = importPathMockMockCommonProxyContext
		pctx.MethodName = "GetCurrentDirModFilePath"
		return
	}()
	ImportPathMockMockGetModFilePathProxyContext = func() (pctx inject.ProxyContext) {
		pctx = importPathMockMockCommonProxyContext
		pctx.MethodName = "GetModFilePath"
		return
	}()

	_ = getIImportPathMockProxy
)

func getIImportPathMockProxy(base IImportPathMock) *ImportPathMockMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &ImportPathMockMock{
		GetByCurrentDirFunc: func() (path string, err error) {
			_gen_begin := time.Now()

			var _gen_r0 string

			var _gen_r1 error

			_gen_ctx := ImportPathMockMockGetByCurrentDirProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_res := _gen_cf(_gen_ctx, base.GetByCurrentDir, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(string)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

				_gen_tmpr1, _gen_exist := _gen_res[1].(error)
				if _gen_exist {
					_gen_r1 = _gen_tmpr1
				}

			} else {
				_gen_r0, _gen_r1 = base.GetByCurrentDir()
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0, _gen_r1
		},

		GetCurrentDirModFilePathFunc: func() (modDir string, modPath string, err error) {
			_gen_begin := time.Now()

			var _gen_r0 string

			var _gen_r1 string

			var _gen_r2 error

			_gen_ctx := ImportPathMockMockGetCurrentDirModFilePathProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_res := _gen_cf(_gen_ctx, base.GetCurrentDirModFilePath, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(string)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

				_gen_tmpr1, _gen_exist := _gen_res[1].(string)
				if _gen_exist {
					_gen_r1 = _gen_tmpr1
				}

				_gen_tmpr2, _gen_exist := _gen_res[2].(error)
				if _gen_exist {
					_gen_r2 = _gen_tmpr2
				}

			} else {
				_gen_r0, _gen_r1, _gen_r2 = base.GetCurrentDirModFilePath()
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0, _gen_r1, _gen_r2
		},

		GetModFilePathFunc: func(dir string) (modDir string, modPath string, err error) {
			_gen_begin := time.Now()

			var _gen_r0 string

			var _gen_r1 string

			var _gen_r2 error

			_gen_ctx := ImportPathMockMockGetModFilePathProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, dir)

				_gen_res := _gen_cf(_gen_ctx, base.GetModFilePath, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(string)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

				_gen_tmpr1, _gen_exist := _gen_res[1].(string)
				if _gen_exist {
					_gen_r1 = _gen_tmpr1
				}

				_gen_tmpr2, _gen_exist := _gen_res[2].(error)
				if _gen_exist {
					_gen_r2 = _gen_tmpr2
				}

			} else {
				_gen_r0, _gen_r1, _gen_r2 = base.GetModFilePath(dir)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0, _gen_r1, _gen_r2
		},
	}
}

func (mockRecv *ImportPathMockMock) GetByCurrentDir() (path string, err error) {
	return mockRecv.GetByCurrentDirFunc()
}

func (mockRecv *ImportPathMockMock) GetCurrentDirModFilePath() (modDir string, modPath string, err error) {
	return mockRecv.GetCurrentDirModFilePathFunc()
}

func (mockRecv *ImportPathMockMock) GetModFilePath(dir string) (modDir string, modPath string, err error) {
	return mockRecv.GetModFilePathFunc(dir)
}

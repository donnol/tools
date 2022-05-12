package apitest

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/donnol/tools/inject"
)

var (
	_gen_customCtxMap = make(map[string]inject.CtxFunc)
)

func RegisterProxyMethod(pctx inject.ProxyContext, cf inject.CtxFunc) {
	_gen_customCtxMap[pctx.Uniq()] = cf
}

type ATMock struct {
	DebugFunc func() *AT

	EqualFunc func(args ...any) *AT

	EqualCodeFunc func(wantCode int) *AT

	EqualThenFunc func(f func(*AT) error, args ...any) *AT

	ErrFunc func() error

	MonkeyRunFunc func() *AT

	NewFunc func() *AT

	PressureRunFunc func(n int, c int) *AT

	PressureRunBatchFunc func(param []PressureParam) *AT

	ResultFunc func(r any) *AT

	RunFunc func() *AT

	SetCookiesFunc func(cookies []*http.Cookie) *AT

	SetHeaderFunc func(header http.Header) *AT

	SetParamFunc func(param any) *AT

	SetPortFunc func(port string) *AT

	WriteFileFunc func(w io.Writer) *AT
}

var (
	_ IAT = &ATMock{}

	aTMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "github.com/donnol/tools/apitest",
		InterfaceName: "IAT",
	}
	ATMockDebugProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockCommonProxyContext
		pctx.MethodName = "Debug"
		return
	}()
	ATMockEqualProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockCommonProxyContext
		pctx.MethodName = "Equal"
		return
	}()
	ATMockEqualCodeProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockCommonProxyContext
		pctx.MethodName = "EqualCode"
		return
	}()
	ATMockEqualThenProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockCommonProxyContext
		pctx.MethodName = "EqualThen"
		return
	}()
	ATMockErrProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockCommonProxyContext
		pctx.MethodName = "Err"
		return
	}()
	ATMockMonkeyRunProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockCommonProxyContext
		pctx.MethodName = "MonkeyRun"
		return
	}()
	ATMockNewProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockCommonProxyContext
		pctx.MethodName = "New"
		return
	}()
	ATMockPressureRunProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockCommonProxyContext
		pctx.MethodName = "PressureRun"
		return
	}()
	ATMockPressureRunBatchProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockCommonProxyContext
		pctx.MethodName = "PressureRunBatch"
		return
	}()
	ATMockResultProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockCommonProxyContext
		pctx.MethodName = "Result"
		return
	}()
	ATMockRunProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockCommonProxyContext
		pctx.MethodName = "Run"
		return
	}()
	ATMockSetCookiesProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockCommonProxyContext
		pctx.MethodName = "SetCookies"
		return
	}()
	ATMockSetHeaderProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockCommonProxyContext
		pctx.MethodName = "SetHeader"
		return
	}()
	ATMockSetParamProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockCommonProxyContext
		pctx.MethodName = "SetParam"
		return
	}()
	ATMockSetPortProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockCommonProxyContext
		pctx.MethodName = "SetPort"
		return
	}()
	ATMockWriteFileProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockCommonProxyContext
		pctx.MethodName = "WriteFile"
		return
	}()

	_ = getIATProxy
)

func getIATProxy(base IAT) *ATMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &ATMock{
		DebugFunc: func() *AT {
			_gen_begin := time.Now()

			var _gen_r0 *AT

			_gen_ctx := ATMockDebugProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_res := _gen_cf(_gen_ctx, base.Debug, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(*AT)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.Debug()
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		EqualFunc: func(args ...any) *AT {
			_gen_begin := time.Now()

			var _gen_r0 *AT

			_gen_ctx := ATMockEqualProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, args...)

				_gen_res := _gen_cf(_gen_ctx, base.Equal, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(*AT)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.Equal(args...)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		EqualCodeFunc: func(wantCode int) *AT {
			_gen_begin := time.Now()

			var _gen_r0 *AT

			_gen_ctx := ATMockEqualCodeProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, wantCode)

				_gen_res := _gen_cf(_gen_ctx, base.EqualCode, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(*AT)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.EqualCode(wantCode)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		EqualThenFunc: func(f func(*AT) error, args ...any) *AT {
			_gen_begin := time.Now()

			var _gen_r0 *AT

			_gen_ctx := ATMockEqualThenProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, f)

				_gen_params = append(_gen_params, args...)

				_gen_res := _gen_cf(_gen_ctx, base.EqualThen, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(*AT)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.EqualThen(f, args...)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		ErrFunc: func() error {
			_gen_begin := time.Now()

			var _gen_r0 error

			_gen_ctx := ATMockErrProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_res := _gen_cf(_gen_ctx, base.Err, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(error)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.Err()
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		MonkeyRunFunc: func() *AT {
			_gen_begin := time.Now()

			var _gen_r0 *AT

			_gen_ctx := ATMockMonkeyRunProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_res := _gen_cf(_gen_ctx, base.MonkeyRun, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(*AT)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.MonkeyRun()
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		NewFunc: func() *AT {
			_gen_begin := time.Now()

			var _gen_r0 *AT

			_gen_ctx := ATMockNewProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_res := _gen_cf(_gen_ctx, base.New, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(*AT)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.New()
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		PressureRunFunc: func(n int, c int) *AT {
			_gen_begin := time.Now()

			var _gen_r0 *AT

			_gen_ctx := ATMockPressureRunProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, n)

				_gen_params = append(_gen_params, c)

				_gen_res := _gen_cf(_gen_ctx, base.PressureRun, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(*AT)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.PressureRun(n, c)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		PressureRunBatchFunc: func(param []PressureParam) *AT {
			_gen_begin := time.Now()

			var _gen_r0 *AT

			_gen_ctx := ATMockPressureRunBatchProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, param)

				_gen_res := _gen_cf(_gen_ctx, base.PressureRunBatch, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(*AT)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.PressureRunBatch(param)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		ResultFunc: func(r any) *AT {
			_gen_begin := time.Now()

			var _gen_r0 *AT

			_gen_ctx := ATMockResultProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, r)

				_gen_res := _gen_cf(_gen_ctx, base.Result, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(*AT)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.Result(r)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		RunFunc: func() *AT {
			_gen_begin := time.Now()

			var _gen_r0 *AT

			_gen_ctx := ATMockRunProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_res := _gen_cf(_gen_ctx, base.Run, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(*AT)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.Run()
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		SetCookiesFunc: func(cookies []*http.Cookie) *AT {
			_gen_begin := time.Now()

			var _gen_r0 *AT

			_gen_ctx := ATMockSetCookiesProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, cookies)

				_gen_res := _gen_cf(_gen_ctx, base.SetCookies, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(*AT)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.SetCookies(cookies)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		SetHeaderFunc: func(header http.Header) *AT {
			_gen_begin := time.Now()

			var _gen_r0 *AT

			_gen_ctx := ATMockSetHeaderProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, header)

				_gen_res := _gen_cf(_gen_ctx, base.SetHeader, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(*AT)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.SetHeader(header)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		SetParamFunc: func(param any) *AT {
			_gen_begin := time.Now()

			var _gen_r0 *AT

			_gen_ctx := ATMockSetParamProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, param)

				_gen_res := _gen_cf(_gen_ctx, base.SetParam, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(*AT)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.SetParam(param)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		SetPortFunc: func(port string) *AT {
			_gen_begin := time.Now()

			var _gen_r0 *AT

			_gen_ctx := ATMockSetPortProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, port)

				_gen_res := _gen_cf(_gen_ctx, base.SetPort, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(*AT)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.SetPort(port)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},

		WriteFileFunc: func(w io.Writer) *AT {
			_gen_begin := time.Now()

			var _gen_r0 *AT

			_gen_ctx := ATMockWriteFileProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, w)

				_gen_res := _gen_cf(_gen_ctx, base.WriteFile, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(*AT)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.WriteFile(w)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},
	}
}

func (mockRecv *ATMock) Debug() *AT {
	return mockRecv.DebugFunc()
}

func (mockRecv *ATMock) Equal(args ...any) *AT {
	return mockRecv.EqualFunc(args...)
}

func (mockRecv *ATMock) EqualCode(wantCode int) *AT {
	return mockRecv.EqualCodeFunc(wantCode)
}

func (mockRecv *ATMock) EqualThen(f func(*AT) error, args ...any) *AT {
	return mockRecv.EqualThenFunc(f, args...)
}

func (mockRecv *ATMock) Err() error {
	return mockRecv.ErrFunc()
}

func (mockRecv *ATMock) MonkeyRun() *AT {
	return mockRecv.MonkeyRunFunc()
}

func (mockRecv *ATMock) New() *AT {
	return mockRecv.NewFunc()
}

func (mockRecv *ATMock) PressureRun(n int, c int) *AT {
	return mockRecv.PressureRunFunc(n, c)
}

func (mockRecv *ATMock) PressureRunBatch(param []PressureParam) *AT {
	return mockRecv.PressureRunBatchFunc(param)
}

func (mockRecv *ATMock) Result(r any) *AT {
	return mockRecv.ResultFunc(r)
}

func (mockRecv *ATMock) Run() *AT {
	return mockRecv.RunFunc()
}

func (mockRecv *ATMock) SetCookies(cookies []*http.Cookie) *AT {
	return mockRecv.SetCookiesFunc(cookies)
}

func (mockRecv *ATMock) SetHeader(header http.Header) *AT {
	return mockRecv.SetHeaderFunc(header)
}

func (mockRecv *ATMock) SetParam(param any) *AT {
	return mockRecv.SetParamFunc(param)
}

func (mockRecv *ATMock) SetPort(port string) *AT {
	return mockRecv.SetPortFunc(port)
}

func (mockRecv *ATMock) WriteFile(w io.Writer) *AT {
	return mockRecv.WriteFileFunc(w)
}

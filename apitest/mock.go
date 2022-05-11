package apitest

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/donnol/tools/inject"
)

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

	_            = getIATProxy
	customCtxMap = make(map[string]inject.CtxFunc)
)

func getIATProxy(base IAT) *ATMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &ATMock{
		DebugFunc: func() *AT {
			begin := time.Now()

			var r0 *AT

			ctx := ATMockDebugProxyContext
			cf, ok := customCtxMap[ctx.Uniq()]
			if ok {
				params := []any{}

				res := cf(ctx, base.Debug, params)

				tmpr0, exist := res[0].(*AT)
				if exist {
					r0 = tmpr0
				}

			} else {
				r0 = base.Debug()
			}

			log.Printf("[ctx: %s]used time: %v\n", ctx.Uniq(), time.Since(begin))

			return r0
		},

		EqualFunc: func(args ...any) *AT {
			begin := time.Now()

			var r0 *AT

			ctx := ATMockEqualProxyContext
			cf, ok := customCtxMap[ctx.Uniq()]
			if ok {
				params := []any{}

				params = append(params, args...)

				res := cf(ctx, base.Equal, params)

				tmpr0, exist := res[0].(*AT)
				if exist {
					r0 = tmpr0
				}

			} else {
				r0 = base.Equal(args...)
			}

			log.Printf("[ctx: %s]used time: %v\n", ctx.Uniq(), time.Since(begin))

			return r0
		},

		EqualCodeFunc: func(wantCode int) *AT {
			begin := time.Now()

			var r0 *AT

			ctx := ATMockEqualCodeProxyContext
			cf, ok := customCtxMap[ctx.Uniq()]
			if ok {
				params := []any{}

				params = append(params, wantCode)

				res := cf(ctx, base.EqualCode, params)

				tmpr0, exist := res[0].(*AT)
				if exist {
					r0 = tmpr0
				}

			} else {
				r0 = base.EqualCode(wantCode)
			}

			log.Printf("[ctx: %s]used time: %v\n", ctx.Uniq(), time.Since(begin))

			return r0
		},

		EqualThenFunc: func(f func(*AT) error, args ...any) *AT {
			begin := time.Now()

			var r0 *AT

			ctx := ATMockEqualThenProxyContext
			cf, ok := customCtxMap[ctx.Uniq()]
			if ok {
				params := []any{}

				params = append(params, f)

				params = append(params, args...)

				res := cf(ctx, base.EqualThen, params)

				tmpr0, exist := res[0].(*AT)
				if exist {
					r0 = tmpr0
				}

			} else {
				r0 = base.EqualThen(f, args...)
			}

			log.Printf("[ctx: %s]used time: %v\n", ctx.Uniq(), time.Since(begin))

			return r0
		},

		ErrFunc: func() error {
			begin := time.Now()

			var r0 error

			ctx := ATMockErrProxyContext
			cf, ok := customCtxMap[ctx.Uniq()]
			if ok {
				params := []any{}

				res := cf(ctx, base.Err, params)

				tmpr0, exist := res[0].(error)
				if exist {
					r0 = tmpr0
				}

			} else {
				r0 = base.Err()
			}

			log.Printf("[ctx: %s]used time: %v\n", ctx.Uniq(), time.Since(begin))

			return r0
		},

		MonkeyRunFunc: func() *AT {
			begin := time.Now()

			var r0 *AT

			ctx := ATMockMonkeyRunProxyContext
			cf, ok := customCtxMap[ctx.Uniq()]
			if ok {
				params := []any{}

				res := cf(ctx, base.MonkeyRun, params)

				tmpr0, exist := res[0].(*AT)
				if exist {
					r0 = tmpr0
				}

			} else {
				r0 = base.MonkeyRun()
			}

			log.Printf("[ctx: %s]used time: %v\n", ctx.Uniq(), time.Since(begin))

			return r0
		},

		NewFunc: func() *AT {
			begin := time.Now()

			var r0 *AT

			ctx := ATMockNewProxyContext
			cf, ok := customCtxMap[ctx.Uniq()]
			if ok {
				params := []any{}

				res := cf(ctx, base.New, params)

				tmpr0, exist := res[0].(*AT)
				if exist {
					r0 = tmpr0
				}

			} else {
				r0 = base.New()
			}

			log.Printf("[ctx: %s]used time: %v\n", ctx.Uniq(), time.Since(begin))

			return r0
		},

		PressureRunFunc: func(n int, c int) *AT {
			begin := time.Now()

			var r0 *AT

			ctx := ATMockPressureRunProxyContext
			cf, ok := customCtxMap[ctx.Uniq()]
			if ok {
				params := []any{}

				params = append(params, n)

				params = append(params, c)

				res := cf(ctx, base.PressureRun, params)

				tmpr0, exist := res[0].(*AT)
				if exist {
					r0 = tmpr0
				}

			} else {
				r0 = base.PressureRun(n, c)
			}

			log.Printf("[ctx: %s]used time: %v\n", ctx.Uniq(), time.Since(begin))

			return r0
		},

		PressureRunBatchFunc: func(param []PressureParam) *AT {
			begin := time.Now()

			var r0 *AT

			ctx := ATMockPressureRunBatchProxyContext
			cf, ok := customCtxMap[ctx.Uniq()]
			if ok {
				params := []any{}

				params = append(params, param)

				res := cf(ctx, base.PressureRunBatch, params)

				tmpr0, exist := res[0].(*AT)
				if exist {
					r0 = tmpr0
				}

			} else {
				r0 = base.PressureRunBatch(param)
			}

			log.Printf("[ctx: %s]used time: %v\n", ctx.Uniq(), time.Since(begin))

			return r0
		},

		ResultFunc: func(r any) *AT {
			begin := time.Now()

			var r0 *AT

			ctx := ATMockResultProxyContext
			cf, ok := customCtxMap[ctx.Uniq()]
			if ok {
				params := []any{}

				params = append(params, r)

				res := cf(ctx, base.Result, params)

				tmpr0, exist := res[0].(*AT)
				if exist {
					r0 = tmpr0
				}

			} else {
				r0 = base.Result(r)
			}

			log.Printf("[ctx: %s]used time: %v\n", ctx.Uniq(), time.Since(begin))

			return r0
		},

		RunFunc: func() *AT {
			begin := time.Now()

			var r0 *AT

			ctx := ATMockRunProxyContext
			cf, ok := customCtxMap[ctx.Uniq()]
			if ok {
				params := []any{}

				res := cf(ctx, base.Run, params)

				tmpr0, exist := res[0].(*AT)
				if exist {
					r0 = tmpr0
				}

			} else {
				r0 = base.Run()
			}

			log.Printf("[ctx: %s]used time: %v\n", ctx.Uniq(), time.Since(begin))

			return r0
		},

		SetCookiesFunc: func(cookies []*http.Cookie) *AT {
			begin := time.Now()

			var r0 *AT

			ctx := ATMockSetCookiesProxyContext
			cf, ok := customCtxMap[ctx.Uniq()]
			if ok {
				params := []any{}

				params = append(params, cookies)

				res := cf(ctx, base.SetCookies, params)

				tmpr0, exist := res[0].(*AT)
				if exist {
					r0 = tmpr0
				}

			} else {
				r0 = base.SetCookies(cookies)
			}

			log.Printf("[ctx: %s]used time: %v\n", ctx.Uniq(), time.Since(begin))

			return r0
		},

		SetHeaderFunc: func(header http.Header) *AT {
			begin := time.Now()

			var r0 *AT

			ctx := ATMockSetHeaderProxyContext
			cf, ok := customCtxMap[ctx.Uniq()]
			if ok {
				params := []any{}

				params = append(params, header)

				res := cf(ctx, base.SetHeader, params)

				tmpr0, exist := res[0].(*AT)
				if exist {
					r0 = tmpr0
				}

			} else {
				r0 = base.SetHeader(header)
			}

			log.Printf("[ctx: %s]used time: %v\n", ctx.Uniq(), time.Since(begin))

			return r0
		},

		SetParamFunc: func(param any) *AT {
			begin := time.Now()

			var r0 *AT

			ctx := ATMockSetParamProxyContext
			cf, ok := customCtxMap[ctx.Uniq()]
			if ok {
				params := []any{}

				params = append(params, param)

				res := cf(ctx, base.SetParam, params)

				tmpr0, exist := res[0].(*AT)
				if exist {
					r0 = tmpr0
				}

			} else {
				r0 = base.SetParam(param)
			}

			log.Printf("[ctx: %s]used time: %v\n", ctx.Uniq(), time.Since(begin))

			return r0
		},

		SetPortFunc: func(port string) *AT {
			begin := time.Now()

			var r0 *AT

			ctx := ATMockSetPortProxyContext
			cf, ok := customCtxMap[ctx.Uniq()]
			if ok {
				params := []any{}

				params = append(params, port)

				res := cf(ctx, base.SetPort, params)

				tmpr0, exist := res[0].(*AT)
				if exist {
					r0 = tmpr0
				}

			} else {
				r0 = base.SetPort(port)
			}

			log.Printf("[ctx: %s]used time: %v\n", ctx.Uniq(), time.Since(begin))

			return r0
		},

		WriteFileFunc: func(w io.Writer) *AT {
			begin := time.Now()

			var r0 *AT

			ctx := ATMockWriteFileProxyContext
			cf, ok := customCtxMap[ctx.Uniq()]
			if ok {
				params := []any{}

				params = append(params, w)

				res := cf(ctx, base.WriteFile, params)

				tmpr0, exist := res[0].(*AT)
				if exist {
					r0 = tmpr0
				}

			} else {
				r0 = base.WriteFile(w)
			}

			log.Printf("[ctx: %s]used time: %v\n", ctx.Uniq(), time.Since(begin))

			return r0
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

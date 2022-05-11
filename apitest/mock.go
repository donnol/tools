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

	customCtxMap = make(map[string]inject.CtxFunc)
	_            = getIATProxy
)

func RegisterCustomProxyMethod(ctx inject.ProxyContext, f inject.CtxFunc) {
	customCtxMap[ctx.Uniq()] = f
}

func getIATProxy(base IAT) *ATMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &ATMock{
		DebugFunc: func() *AT {
			var r *AT
			begin := time.Now()

			cf, ok := customCtxMap[ATMockDebugProxyContext.Uniq()]
			if ok {
				// 自定义around，需要自己决定怎么调用，一般需要包含下面的默认调用
				res := cf(ATMockDebugProxyContext, base.Debug, nil)
				r = res[0].(*AT)
			} else {
				// 默认调用
				r = base.Debug()
			}

			log.Printf("[ctx: %s]used time: %v\n", ATMockDebugProxyContext.Uniq(), time.Since(begin))

			return r
		},
		EqualFunc: func(args ...any) *AT {
			var r *AT
			begin := time.Now()

			cf, ok := customCtxMap[ATMockEqualProxyContext.Uniq()]
			if ok {
				// 自定义around，需要自己决定怎么调用，一般需要包含下面的默认调用
				res := cf(ATMockEqualProxyContext, base.Equal, args)
				r = res[0].(*AT)
			} else {
				// 默认调用
				r = base.Equal(args...)
			}

			log.Printf("[ctx: %s]used time: %v\n", ATMockDebugProxyContext.Uniq(), time.Since(begin))

			return r
		},
		EqualCodeFunc: base.EqualCode,
		EqualThenFunc: func(fa func(*AT) error, args ...any) *AT {
			var r *AT
			begin := time.Now()

			cf, ok := customCtxMap[ATMockEqualThenProxyContext.Uniq()]
			if ok {
				// 自定义around，需要自己决定怎么调用，一般需要包含下面的默认调用
				allArg := []any{fa}
				allArg = append(allArg, args...)
				res := cf(ATMockEqualThenProxyContext, base.EqualThen, allArg)
				r = res[0].(*AT)
			} else {
				// 默认调用
				r = base.EqualThen(fa, args...)
			}

			log.Printf("[ctx: %s]used time: %v\n", ATMockDebugProxyContext.Uniq(), time.Since(begin))

			return r
		},
		ErrFunc:              base.Err,
		MonkeyRunFunc:        base.MonkeyRun,
		NewFunc:              base.New,
		PressureRunFunc:      base.PressureRun,
		PressureRunBatchFunc: base.PressureRunBatch,
		RunFunc:              base.Run,
		SetCookiesFunc:       base.SetCookies,
		SetHeaderFunc:        base.SetHeader,
		SetPortFunc:          base.SetPort,
		WriteFileFunc:        base.WriteFile,
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

type ATMockMock struct {
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
	_ IATMock = &ATMockMock{}

	aTMockMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "github.com/donnol/tools/apitest",
		InterfaceName: "IATMock",
	}
	ATMockMockDebugProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockMockCommonProxyContext
		pctx.MethodName = "Debug"
		return
	}()
	ATMockMockEqualProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockMockCommonProxyContext
		pctx.MethodName = "Equal"
		return
	}()
	ATMockMockEqualCodeProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockMockCommonProxyContext
		pctx.MethodName = "EqualCode"
		return
	}()
	ATMockMockEqualThenProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockMockCommonProxyContext
		pctx.MethodName = "EqualThen"
		return
	}()
	ATMockMockErrProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockMockCommonProxyContext
		pctx.MethodName = "Err"
		return
	}()
	ATMockMockMonkeyRunProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockMockCommonProxyContext
		pctx.MethodName = "MonkeyRun"
		return
	}()
	ATMockMockNewProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockMockCommonProxyContext
		pctx.MethodName = "New"
		return
	}()
	ATMockMockPressureRunProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockMockCommonProxyContext
		pctx.MethodName = "PressureRun"
		return
	}()
	ATMockMockPressureRunBatchProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockMockCommonProxyContext
		pctx.MethodName = "PressureRunBatch"
		return
	}()
	ATMockMockResultProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockMockCommonProxyContext
		pctx.MethodName = "Result"
		return
	}()
	ATMockMockRunProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockMockCommonProxyContext
		pctx.MethodName = "Run"
		return
	}()
	ATMockMockSetCookiesProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockMockCommonProxyContext
		pctx.MethodName = "SetCookies"
		return
	}()
	ATMockMockSetHeaderProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockMockCommonProxyContext
		pctx.MethodName = "SetHeader"
		return
	}()
	ATMockMockSetParamProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockMockCommonProxyContext
		pctx.MethodName = "SetParam"
		return
	}()
	ATMockMockSetPortProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockMockCommonProxyContext
		pctx.MethodName = "SetPort"
		return
	}()
	ATMockMockWriteFileProxyContext = func() (pctx inject.ProxyContext) {
		pctx = aTMockMockCommonProxyContext
		pctx.MethodName = "WriteFile"
		return
	}()
)

func (mockRecv *ATMockMock) Debug() *AT {
	return mockRecv.DebugFunc()
}

func (mockRecv *ATMockMock) Equal(args ...any) *AT {
	return mockRecv.EqualFunc(args...)
}

func (mockRecv *ATMockMock) EqualCode(wantCode int) *AT {
	return mockRecv.EqualCodeFunc(wantCode)
}

func (mockRecv *ATMockMock) EqualThen(f func(*AT) error, args ...any) *AT {
	return mockRecv.EqualThenFunc(f, args...)
}

func (mockRecv *ATMockMock) Err() error {
	return mockRecv.ErrFunc()
}

func (mockRecv *ATMockMock) MonkeyRun() *AT {
	return mockRecv.MonkeyRunFunc()
}

func (mockRecv *ATMockMock) New() *AT {
	return mockRecv.NewFunc()
}

func (mockRecv *ATMockMock) PressureRun(n int, c int) *AT {
	return mockRecv.PressureRunFunc(n, c)
}

func (mockRecv *ATMockMock) PressureRunBatch(param []PressureParam) *AT {
	return mockRecv.PressureRunBatchFunc(param)
}

func (mockRecv *ATMockMock) Result(r any) *AT {
	return mockRecv.ResultFunc(r)
}

func (mockRecv *ATMockMock) Run() *AT {
	return mockRecv.RunFunc()
}

func (mockRecv *ATMockMock) SetCookies(cookies []*http.Cookie) *AT {
	return mockRecv.SetCookiesFunc(cookies)
}

func (mockRecv *ATMockMock) SetHeader(header http.Header) *AT {
	return mockRecv.SetHeaderFunc(header)
}

func (mockRecv *ATMockMock) SetParam(param any) *AT {
	return mockRecv.SetParamFunc(param)
}

func (mockRecv *ATMockMock) SetPort(port string) *AT {
	return mockRecv.SetPortFunc(port)
}

func (mockRecv *ATMockMock) WriteFile(w io.Writer) *AT {
	return mockRecv.WriteFileFunc(w)
}

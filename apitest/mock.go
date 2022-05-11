package apitest

import (
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

func proxyHelper(ctx inject.ProxyContext, method any, args []any) (res []any) {
	// 执行前可以做很多东西

	begin := time.Now()

	switch ctx.Uniq() {
	case ATMockDebugProxyContext.Uniq():
		cf, ok := customCtxMap[ctx.Uniq()]
		if ok {
			// 自定义around，需要自己决定怎么调用，一般需要包含下面的默认调用
			res = cf(ctx, method, args)
		} else {
			// 默认调用
			f := method.(func() *AT)
			r1 := f()
			res = append(res, r1)
		}

	case ATMockEqualProxyContext.Uniq():
		cf, ok := customCtxMap[ctx.Uniq()]
		if ok {
			// 自定义around，需要自己决定怎么调用，一般需要包含下面的默认调用
			res = cf(ctx, method, args)
		} else {
			// 默认调用
			f := method.(func(args ...any) *AT)
			r1 := f(args...)
			res = append(res, r1)
		}

	case ATMockEqualCodeProxyContext.Uniq():
	case ATMockEqualThenProxyContext.Uniq():
		cf, ok := customCtxMap[ctx.Uniq()]
		if ok {
			// 自定义around，需要自己决定怎么调用，一般需要包含下面的默认调用
			res = cf(ctx, method, args)
		} else {
			// 默认调用
			f := method.(func(func(*AT) error, ...any) *AT)
			a1 := args[0].(func(*AT) error)
			r1 := f(a1, args[1:]...)
			res = append(res, r1)
		}

	case ATMockErrProxyContext.Uniq():
	case ATMockMonkeyRunProxyContext.Uniq():
	case ATMockNewProxyContext.Uniq():
	case ATMockPressureRunProxyContext.Uniq():
	case ATMockPressureRunBatchProxyContext.Uniq():
	case ATMockResultProxyContext.Uniq():
	case ATMockRunProxyContext.Uniq():
	case ATMockSetCookiesProxyContext.Uniq():
	case ATMockSetHeaderProxyContext.Uniq():
	case ATMockSetParamProxyContext.Uniq():
	case ATMockSetPortProxyContext.Uniq():
	case ATMockWriteFileProxyContext.Uniq():
	}

	// 执行后可以做很多东西

	used := time.Since(begin)
	log.Printf("[ctx: %s]used time: %v\n", ctx.Uniq(), used)

	return
}

func getIATProxy(base IAT) *ATMock {
	return &ATMock{
		DebugFunc: func() *AT {
			res := proxyHelper(ATMockDebugProxyContext, base.Debug, nil)
			return res[0].(*AT)
		},
		EqualFunc: func(args ...any) *AT {
			res := proxyHelper(ATMockEqualProxyContext, base.Equal, args)
			return res[0].(*AT)
		},
		EqualCodeFunc: base.EqualCode,
		EqualThenFunc: func(f func(*AT) error, args ...any) *AT {
			allArg := []any{f}
			allArg = append(allArg, args...)
			res := proxyHelper(ATMockEqualThenProxyContext, base.EqualThen, allArg)
			return res[0].(*AT)
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

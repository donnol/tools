package apitest

import (
	"io"
	"net/http"
)

type ATMock struct {
	DebugFunc func() *AT

	EqualFunc func(args ...interface{}) *AT

	EqualCodeFunc func(wantCode int) *AT

	EqualThenFunc func(f func(*AT) error, args ...interface{}) *AT

	ErrFunc func() error

	MonkeyRunFunc func() *AT

	NewFunc func() *AT

	PressureRunFunc func(n int, c int) *AT

	PressureRunBatchFunc func(param []PressureParam) *AT

	ResultFunc func(r interface{}) *AT

	RunFunc func() *AT

	SetCookiesFunc func(cookies []*http.Cookie) *AT

	SetHeaderFunc func(header http.Header) *AT

	SetParamFunc func(param interface{}) *AT

	SetPortFunc func(port string) *AT

	WriteFileFunc func(w io.Writer) *AT
}

var _ IAT = &ATMock{}

func (mockRecv *ATMock) Debug() *AT {
	return mockRecv.DebugFunc()
}

func (mockRecv *ATMock) Equal(args ...interface{}) *AT {
	return mockRecv.EqualFunc(args...)
}

func (mockRecv *ATMock) EqualCode(wantCode int) *AT {
	return mockRecv.EqualCodeFunc(wantCode)
}

func (mockRecv *ATMock) EqualThen(f func(*AT) error, args ...interface{}) *AT {
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

func (mockRecv *ATMock) Result(r interface{}) *AT {
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

func (mockRecv *ATMock) SetParam(param interface{}) *AT {
	return mockRecv.SetParamFunc(param)
}

func (mockRecv *ATMock) SetPort(port string) *AT {
	return mockRecv.SetPortFunc(port)
}

func (mockRecv *ATMock) WriteFile(w io.Writer) *AT {
	return mockRecv.WriteFileFunc(w)
}

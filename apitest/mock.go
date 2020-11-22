package apitest

import (
	"io"
	"net/http"
)

type IATMock struct {
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

var _ IAT = &IATMock{}

func (*IATMock) Debug() *AT {
	panic("Need to be implement!")
}

func (*IATMock) Equal(args ...interface{}) *AT {
	panic("Need to be implement!")
}

func (*IATMock) EqualCode(wantCode int) *AT {
	panic("Need to be implement!")
}

func (*IATMock) EqualThen(f func(*AT) error, args ...interface{}) *AT {
	panic("Need to be implement!")
}

func (*IATMock) Err() error {
	panic("Need to be implement!")
}

func (*IATMock) MonkeyRun() *AT {
	panic("Need to be implement!")
}

func (*IATMock) New() *AT {
	panic("Need to be implement!")
}

func (*IATMock) PressureRun(n int, c int) *AT {
	panic("Need to be implement!")
}

func (*IATMock) PressureRunBatch(param []PressureParam) *AT {
	panic("Need to be implement!")
}

func (*IATMock) Result(r interface{}) *AT {
	panic("Need to be implement!")
}

func (*IATMock) Run() *AT {
	panic("Need to be implement!")
}

func (*IATMock) SetCookies(cookies []*http.Cookie) *AT {
	panic("Need to be implement!")
}

func (*IATMock) SetHeader(header http.Header) *AT {
	panic("Need to be implement!")
}

func (*IATMock) SetParam(param interface{}) *AT {
	panic("Need to be implement!")
}

func (*IATMock) SetPort(port string) *AT {
	panic("Need to be implement!")
}

func (*IATMock) WriteFile(w io.Writer) *AT {
	panic("Need to be implement!")
}

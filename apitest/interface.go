package apitest

import (
	"io"
	"net/http"
)

type IAT interface {
	Debug() *AT
	Equal(args ...any) *AT
	EqualCode(wantCode int) *AT
	EqualThen(f func(*AT) error, args ...any) *AT
	Err() error
	MonkeyRun() *AT
	New() *AT
	PressureRun(n int, c int) *AT
	PressureRunBatch(param []PressureParam) *AT
	Result(r any) *AT
	Run() *AT
	SetCookies(cookies []*http.Cookie) *AT
	SetHeader(header http.Header) *AT
	SetParam(param any) *AT
	SetPort(port string) *AT
	WriteFile(w io.Writer) *AT
}

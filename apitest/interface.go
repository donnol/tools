package apitest

import (
	"io"
	"net/http"
)

type IAT interface {
	Debug() *AT
	Equal(args ...interface{}) *AT
	EqualCode(wantCode int) *AT
	EqualThen(f func(*AT) error, args ...interface{}) *AT
	Err() error
	MonkeyRun() *AT
	New() *AT
	PressureRun(n int, c int) *AT
	PressureRunBatch(param []PressureParam) *AT
	Result(r interface{}) *AT
	Run() *AT
	SetCookies(cookies []*http.Cookie) *AT
	SetHeader(header http.Header) *AT
	SetParam(param interface{}) *AT
	SetPort(port string) *AT
	WriteFile(w io.Writer) *AT
}

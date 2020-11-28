package inject

import "reflect"

type IocMock struct {
	InjectFunc func(v interface{}) (err error)

	RegisterProviderFunc func(v interface{}) (err error)
}

var _ IIoc = &IocMock{}

func (mockRecv *IocMock) Inject(v interface{}) (err error) {
	return mockRecv.InjectFunc(v)
}

func (mockRecv *IocMock) RegisterProvider(v interface{}) (err error) {
	return mockRecv.RegisterProviderFunc(v)
}

type ProxyContextMock struct {
	LogfFunc func(format string, args ...interface{})

	StringFunc func() string
}

var _ IProxyContext = &ProxyContextMock{}

func (mockRecv *ProxyContextMock) Logf(format string, args ...interface{}) {
	mockRecv.LogfFunc(format, args...)
}

func (mockRecv *ProxyContextMock) String() string {
	return mockRecv.StringFunc()
}

type AroundMock struct {
	AfterFunc func(pctx ProxyContext)

	BeforeFunc func(pctx ProxyContext)
}

var _ IAround = &AroundMock{}

func (mockRecv *AroundMock) After(pctx ProxyContext) {
	mockRecv.AfterFunc(pctx)
}

func (mockRecv *AroundMock) Before(pctx ProxyContext) {
	mockRecv.BeforeFunc(pctx)
}

type proxyImplMock struct {
	AddHookFunc func(hooks ...Hook)

	WrapFunc func(provider interface{}, mock interface{}, hooks ...Hook) interface{}
}

var _ IproxyImpl = &proxyImplMock{}

func (mockRecv *proxyImplMock) AddHook(hooks ...Hook) {
	mockRecv.AddHookFunc(hooks...)
}

func (mockRecv *proxyImplMock) Wrap(provider interface{}, mock interface{}, hooks ...Hook) interface{} {
	return mockRecv.WrapFunc(provider, mock, hooks...)
}

type ProxyMock struct {
	AddHookFunc func(...Hook)

	WrapFunc func(provider interface{}, mock interface{}, hooks ...Hook) interface{}
}

var _ Proxy = &ProxyMock{}

func (mockRecv *ProxyMock) AddHook(p0 ...Hook) {
	mockRecv.AddHookFunc(p0...)
}

func (mockRecv *ProxyMock) Wrap(provider interface{}, mock interface{}, hooks ...Hook) interface{} {
	return mockRecv.WrapFunc(provider, mock, hooks...)
}

type HookMock struct {
	AfterFunc func(ProxyContext)

	BeforeFunc func(ProxyContext)
}

var _ Hook = &HookMock{}

func (mockRecv *HookMock) After(p0 ProxyContext) {
	mockRecv.AfterFunc(p0)
}

func (mockRecv *HookMock) Before(p0 ProxyContext) {
	mockRecv.BeforeFunc(p0)
}

type CallerMock struct {
	CallFunc func(args []reflect.Value) []reflect.Value
}

var _ Caller = &CallerMock{}

func (mockRecv *CallerMock) Call(args []reflect.Value) []reflect.Value {
	return mockRecv.CallFunc(args)
}

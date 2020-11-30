package inject

import "reflect"

type ProxyContextMockMock struct {
	LogfFunc func(format string, args ...interface{})

	StringFunc func() string
}

var _ IProxyContextMock = &ProxyContextMockMock{}

func (mockRecv *ProxyContextMockMock) Logf(format string, args ...interface{}) {
	mockRecv.LogfFunc(format, args...)
}

func (mockRecv *ProxyContextMockMock) String() string {
	return mockRecv.StringFunc()
}

type ArounderMockMock struct {
	AroundFunc func(pctx ProxyContext, method reflect.Value, args []reflect.Value) []reflect.Value
}

var _ IArounderMock = &ArounderMockMock{}

func (mockRecv *ArounderMockMock) Around(pctx ProxyContext, method reflect.Value, args []reflect.Value) []reflect.Value {
	return mockRecv.AroundFunc(pctx, method, args)
}

type AroundMockMock struct {
	AfterFunc func(pctx ProxyContext)

	BeforeFunc func(pctx ProxyContext)
}

var _ IAroundMock = &AroundMockMock{}

func (mockRecv *AroundMockMock) After(pctx ProxyContext) {
	mockRecv.AfterFunc(pctx)
}

func (mockRecv *AroundMockMock) Before(pctx ProxyContext) {
	mockRecv.BeforeFunc(pctx)
}

type ProxyContextMock struct {
	LogShortfFunc func(format string, args ...interface{})

	LogfFunc func(format string, args ...interface{})

	StringFunc func() string
}

var _ IProxyContext = &ProxyContextMock{}

func (mockRecv *ProxyContextMock) LogShortf(format string, args ...interface{}) {
	mockRecv.LogShortfFunc(format, args...)
}

func (mockRecv *ProxyContextMock) Logf(format string, args ...interface{}) {
	mockRecv.LogfFunc(format, args...)
}

func (mockRecv *ProxyContextMock) String() string {
	return mockRecv.StringFunc()
}

type ProxyMockMock struct {
	AroundFunc func(provider interface{}, mock interface{}, arounder Arounder) interface{}
}

var _ IProxyMock = &ProxyMockMock{}

func (mockRecv *ProxyMockMock) Around(provider interface{}, mock interface{}, arounder Arounder) interface{} {
	return mockRecv.AroundFunc(provider, mock, arounder)
}

type proxyImplMock struct {
	AroundFunc func(provider interface{}, mock interface{}, arounder Arounder) interface{}
}

var _ IproxyImpl = &proxyImplMock{}

func (mockRecv *proxyImplMock) Around(provider interface{}, mock interface{}, arounder Arounder) interface{} {
	return mockRecv.AroundFunc(provider, mock, arounder)
}

type AroundFuncMock struct {
	AroundFunc func(pctx ProxyContext, method reflect.Value, args []reflect.Value) []reflect.Value
}

var _ IAroundFunc = &AroundFuncMock{}

func (mockRecv *AroundFuncMock) Around(pctx ProxyContext, method reflect.Value, args []reflect.Value) []reflect.Value {
	return mockRecv.AroundFunc(pctx, method, args)
}

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

type IocMockMock struct {
	InjectFunc func(v interface{}) (err error)

	RegisterProviderFunc func(v interface{}) (err error)
}

var _ IIocMock = &IocMockMock{}

func (mockRecv *IocMockMock) Inject(v interface{}) (err error) {
	return mockRecv.InjectFunc(v)
}

func (mockRecv *IocMockMock) RegisterProvider(v interface{}) (err error) {
	return mockRecv.RegisterProviderFunc(v)
}

type ProxyMock struct {
	AroundFunc func(provider interface{}, mock interface{}, arounder Arounder) interface{}
}

var _ Proxy = &ProxyMock{}

func (mockRecv *ProxyMock) Around(provider interface{}, mock interface{}, arounder Arounder) interface{} {
	return mockRecv.AroundFunc(provider, mock, arounder)
}

type ArounderMock struct {
	AroundFunc func(pctx ProxyContext, method reflect.Value, args []reflect.Value) []reflect.Value
}

var _ Arounder = &ArounderMock{}

func (mockRecv *ArounderMock) Around(pctx ProxyContext, method reflect.Value, args []reflect.Value) []reflect.Value {
	return mockRecv.AroundFunc(pctx, method, args)
}

package inject

import "reflect"

type IproxyImplMock struct {
	AddHookFunc func(hooks ...Hook)

	WrapFunc func(provider interface{}, mock interface{}, hooks ...Hook) interface{}
}

var _ IproxyImpl = &IproxyImplMock{}

func (*IproxyImplMock) AddHook(hooks ...Hook) {
	panic("Need to be implement!")
}

func (*IproxyImplMock) Wrap(provider interface{}, mock interface{}, hooks ...Hook) interface{} {
	panic("Need to be implement!")
}

type IIocMock struct {
	InjectFunc func(v interface{}) (err error)

	RegisterProviderFunc func(v interface{}) (err error)
}

var _ IIoc = &IIocMock{}

func (*IIocMock) Inject(v interface{}) (err error) {
	panic("Need to be implement!")
}

func (*IIocMock) RegisterProvider(v interface{}) (err error) {
	panic("Need to be implement!")
}

type ProxyMock struct {
	AddHookFunc func(...Hook)

	WrapFunc func(provider interface{}, mock interface{}, hooks ...Hook) interface{}
}

var _ Proxy = &ProxyMock{}

func (*ProxyMock) AddHook(...Hook) {
	panic("Need to be implement!")
}

func (*ProxyMock) Wrap(provider interface{}, mock interface{}, hooks ...Hook) interface{} {
	panic("Need to be implement!")
}

type HookMock struct {
	BeforeFunc func(ProxyContext)

	AfterFunc func(ProxyContext)
}

var _ Hook = &HookMock{}

func (*HookMock) Before(ProxyContext) {
	panic("Need to be implement!")
}

func (*HookMock) After(ProxyContext) {
	panic("Need to be implement!")
}

type CallerMock struct {
	CallFunc func(args []reflect.Value) []reflect.Value
}

var _ Caller = &CallerMock{}

func (*CallerMock) Call(args []reflect.Value) []reflect.Value {
	panic("Need to be implement!")
}

type IProxyContextMock struct {
	LogfFunc func(format string, args ...interface{})

	StringFunc func() string
}

var _ IProxyContext = &IProxyContextMock{}

func (*IProxyContextMock) Logf(format string, args ...interface{}) {
	panic("Need to be implement!")
}

func (*IProxyContextMock) String() string {
	panic("Need to be implement!")
}

type IAroundMock struct {
	AfterFunc func(pctx ProxyContext)

	BeforeFunc func(pctx ProxyContext)
}

var _ IAround = &IAroundMock{}

func (*IAroundMock) After(pctx ProxyContext) {
	panic("Need to be implement!")
}

func (*IAroundMock) Before(pctx ProxyContext) {
	panic("Need to be implement!")
}

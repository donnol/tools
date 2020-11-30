package inject

import "reflect"

type IIocMock interface {
	Inject(v interface{}) (err error)
	RegisterProvider(v interface{}) (err error)
}

type IProxyContext interface {
	LogShortf(format string, args ...interface{})
	Logf(format string, args ...interface{})
	String() string
}

type IProxyContextMock interface {
	Logf(format string, args ...interface{})
	String() string
}

type IArounderMock interface {
	Around(pctx ProxyContext, method reflect.Value, args []reflect.Value) []reflect.Value
}

type IAroundMock interface {
	After(pctx ProxyContext)
	Before(pctx ProxyContext)
}

type IProxyMock interface {
	Around(provider interface{}, mock interface{}, arounder Arounder) interface{}
}

type IAroundFunc interface {
	Around(pctx ProxyContext, method reflect.Value, args []reflect.Value) []reflect.Value
}

type IproxyImpl interface {
	Around(provider interface{}, mock interface{}, arounder Arounder) interface{}
}

type IIoc interface {
	Inject(v interface{}) (err error)
	RegisterProvider(v interface{}) (err error)
}

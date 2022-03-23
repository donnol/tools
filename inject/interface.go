package inject

import "reflect"

type IIocMock interface {
	Inject(v any) (err error)
	RegisterProvider(v any) (err error)
}

type IProxyContext interface {
	LogShortf(format string, args ...any)
	Logf(format string, args ...any)
	String() string
}

type IProxyContextMock interface {
	Logf(format string, args ...any)
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
	Around(provider any, mock any, arounder Arounder) any
}

type IAroundFunc interface {
	Around(pctx ProxyContext, method reflect.Value, args []reflect.Value) []reflect.Value
}

type IproxyImpl interface {
	Around(provider any, mock any, arounder Arounder) any
}

type IIoc interface {
	Inject(v any) (err error)
	RegisterProvider(v any) (err error)
}

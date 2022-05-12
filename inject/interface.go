package inject

import "reflect"

type IProxyContext interface {
	LogShortf(format string, args ...any)
	Logf(format string, args ...any)
	String() string
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

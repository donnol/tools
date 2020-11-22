package inject

type IProxyContext interface {
	Logf(format string, args ...interface{})
	String() string
}

type IAround interface {
	After(pctx ProxyContext)
	Before(pctx ProxyContext)
}

type IproxyImpl interface {
	AddHook(hooks ...Hook)
	Wrap(provider interface{}, mock interface{}, hooks ...Hook) interface{}
}

type IIoc interface {
	Inject(v interface{}) (err error)
	RegisterProvider(v interface{}) (err error)
}

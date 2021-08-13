package testtype

import (
	"context"

	"github.com/donnol/tools/inject"
	"github.com/donnol/tools/parser/testtype/a"
	"github.com/donnol/tools/parser/testtype/big"
)

type Mock struct {
	HeadFunc func(ctx context.Context, p a.Pa) (big.Int, error)
}

var (
	_ I = &Mock{}

	mockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "github.com/donnol/tools/parser/testtype",
		InterfaceName: "I",
	}
	MockHeadProxyContext = func() (pctx inject.ProxyContext) {
		pctx = mockCommonProxyContext
		pctx.MethodName = "Head"
		return
	}()
)

func (mockRecv *Mock) Head(ctx context.Context, p a.Pa) (big.Int, error) {
	return mockRecv.HeadFunc(ctx, p)
}

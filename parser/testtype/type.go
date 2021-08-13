package testtype

import (
	"context"

	"github.com/donnol/tools/parser/testtype/a"
	"github.com/donnol/tools/parser/testtype/b"
)

type I interface {
	Head(ctx context.Context, p a.Pa) (b.Pb, error)
}

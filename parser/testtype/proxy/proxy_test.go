package proxy

import (
	"testing"

	"github.com/donnol/tools/importpath"
	"github.com/donnol/tools/parser"
)

func TestInspectProxy(t *testing.T) {
	p := parser.New(parser.Option{
		ReplaceCallExpr: true,
	})
	ip := &importpath.ImportPath{}
	path, err := ip.GetByCurrentDir()
	if err != nil {
		t.Fatal(err)
	}
	pkgs, err := p.ParseByGoPackages(path)
	if err != nil {
		t.Fatal(err)
	}
	_ = pkgs
}

package parser

import (
	"testing"

	"github.com/donnol/tools/importpath"
)

func TestParseAST(t *testing.T) {
	p := New(Option{
		UseSourceImporter: true,
	})
	ip := &importpath.ImportPath{}
	path, err := ip.GetByCurrentDir()
	if err != nil {
		t.Fatal(err)
	}

	r, err := p.ParseAST(path)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	// For debug
	// if err := json.IndentToStdout(r); err != nil {
	// 	t.Fatal(err)
	// }

	for _, s := range r {
		id := s.MakeInterface()
		if id != "" {
			_ = id
			t.Logf("interface:\n%s\n", id)
		}
	}
}

func TestParseASTReplaceImportPath(t *testing.T) {
	p := New(Option{
		UseSourceImporter: true,
		ReplaceImportPath: true,
		FromPath:          "github.com/donnol/tools",
		ToPath:            "github.com/donnol/tools", // 替换为原来路径
		Output:            nil,                       // 写回原来文件
	})
	ip := &importpath.ImportPath{}
	path, err := ip.GetByCurrentDir()
	if err != nil {
		t.Fatal(err)
	}

	_, err = p.ParseAST(path)
	if err != nil {
		t.Fatalf("%+v", err)
	}
}

func TestParseByGoPackages(t *testing.T) {
	p := New(Option{})
	ip := &importpath.ImportPath{}
	path, err := ip.GetByCurrentDir()
	if err != nil {
		t.Fatal(err)
	}

	r, err := p.ParseByGoPackages(path)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	for _, single := range r.Pkgs {
		for _, stru := range single.Structs {
			if stru.Name == "Struct" {
				t.Logf("=== struct: %+v\n", stru)

				id := stru.MakeInterface()
				if id != "" {
					_ = id
					t.Logf("interface:\n%s\n", id)
				}
			}
		}
	}
}

// var _ IStruct2 = &Struct{}

// type IStruct2 interface {
// 	Demo(in types.Array) types.Basic
// 	MakeInterface() string
// 	PointerMethod(in types.Basic) types.Slice
// 	String(f Field, ip importpath.ImportPath)
// 	TypeAlias(p Field, ip importpath.ImportPath)
// }

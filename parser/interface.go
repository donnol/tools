package parser

import (
	"go/ast"
	"go/types"

	"github.com/donnol/tools/importpath"
)

type IParser interface {
	GetPkgInfo() PkgInfo
	ParseAST(importPath string) (structs []Struct, err error)
}

type IStruct interface {
	Demo(in types.Array) types.Basic
	MakeInterface() string
	String(f Field, ip importpath.ImportPath)
	TypeAlias(p Field, ip importpath.ImportPath)
}

type IPkgInfo interface {
	GetDir() string
	GetPkgName() string
}

type Ivisitor interface {
	Visit(node ast.Node) (w ast.Visitor)
}

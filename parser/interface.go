package parser

import (
	"go/ast"
	"go/types"

	"github.com/donnol/tools/importpath"
	"golang.org/x/tools/go/packages"
)

type IInterface interface{ MakeMock() string }

type IPackages interface {
	LookupPkg(name string) (Package, bool)
}

type IPackage interface {
	NewGoFileWithSuffix(suffix string) (file string)
	SaveInterface(file string) error
	SaveMock(file string) error
}

type IStruct interface {
	Demo(in types.Array) types.Basic
	MakeInterface() string
	PointerMethod(in types.Basic) types.Slice
	String(f Field, ip importpath.ImportPath)
	TypeAlias(p Field, ip importpath.ImportPath)
}

type IPkgInfo interface {
	GetDir() string
	GetPkgName() string
}

type IInspector interface {
	InspectFile(file *ast.File) (result FileResult)
	InspectPkg(pkg *packages.Package) Package
}

type IParser interface {
	GetPkgInfo() PkgInfo
	ParseAST(importPath string) (structs []Struct, err error)
	ParseByGoPackages(patterns ...string) (result Packages, err error)
}

type Ivisitor interface {
	Visit(node ast.Node) (w ast.Visitor)
}

package parser

import (
	"go/ast"
	"go/types"

	"github.com/donnol/tools/importpath"
	"golang.org/x/tools/go/packages"
)

type IPkgInfo interface {
	GetDir() string
	GetPkgName() string
}

type IInspector interface {
	InspectFile(file *ast.File) (result FileResult)
	InspectPkg(pkg *packages.Package) Package
}

type IInspectorMock interface {
	InspectFile(file *ast.File) (result FileResult)
	InspectPkg(pkg *packages.Package) Package
}

type IStructMock interface {
	Demo(in types.Array) types.Basic
	MakeInterface() string
	PointerMethod(in types.Basic) types.Slice
	String(f Field, ip importpath.ImportPath)
	TypeAlias(p Field, ip importpath.ImportPath)
}

type IParser interface {
	GetPkgInfo() PkgInfo
	ParseAST(importPath string) (structs []Struct, err error)
	ParseByGoPackages(patterns ...string) (result Packages, err error)
}

type IStruct interface {
	Demo(in types.Array) types.Basic
	MakeInterface() string
	PointerMethod(in types.Basic) types.Slice
	String(f Field, ip importpath.ImportPath)
	TypeAlias(p Field, ip importpath.ImportPath)
}

type IInterfaceMock interface{ MakeMock() string }

type IPackagesMock interface {
	LookupPkg(name string) (Package, bool)
}

type IInterface interface{ MakeMock() string }

type IPkgInfoMock interface {
	GetDir() string
	GetPkgName() string
}

type IPackageMock interface {
	NewGoFileWithSuffix(suffix string) (file string)
	SaveInterface(file string) error
	SaveMock(file string) error
}

type IParserMock interface {
	GetPkgInfo() PkgInfo
	ParseAST(importPath string) (structs []Struct, err error)
	ParseByGoPackages(patterns ...string) (result Packages, err error)
}

type IPackages interface {
	LookupPkg(name string) (Package, bool)
}

type Ivisitor interface {
	Visit(node ast.Node) (w ast.Visitor)
}

type IvisitorMock interface {
	Visit(node ast.Node) (w ast.Visitor)
}

type IPackage interface {
	NewGoFileWithSuffix(suffix string) (file string)
	SaveInterface(file string) error
	SaveMock(file string) error
}

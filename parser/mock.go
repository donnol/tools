package parser

import (
	"go/ast"
	"go/types"

	"github.com/donnol/tools/importpath"
	"golang.org/x/tools/go/packages"
)

type IInspectorMock struct {
	InspectFileFunc func(file *ast.File) (result FileResult)

	InspectPkgFunc func(pkg *packages.Package) Package
}

var _ IInspector = &IInspectorMock{}

func (*IInspectorMock) InspectFile(file *ast.File) (result FileResult) {
	panic("Need to be implement!")
}

func (*IInspectorMock) InspectPkg(pkg *packages.Package) Package {
	panic("Need to be implement!")
}

type IParserMock struct {
	GetPkgInfoFunc func() PkgInfo

	ParseASTFunc func(importPath string) (structs []Struct, err error)

	ParseByGoPackagesFunc func(patterns ...string) (result Packages, err error)
}

var _ IParser = &IParserMock{}

func (*IParserMock) GetPkgInfo() PkgInfo {
	panic("Need to be implement!")
}

func (*IParserMock) ParseAST(importPath string) (structs []Struct, err error) {
	panic("Need to be implement!")
}

func (*IParserMock) ParseByGoPackages(patterns ...string) (result Packages, err error) {
	panic("Need to be implement!")
}

type IvisitorMock struct {
	VisitFunc func(node ast.Node) (w ast.Visitor)
}

var _ Ivisitor = &IvisitorMock{}

func (*IvisitorMock) Visit(node ast.Node) (w ast.Visitor) {
	panic("Need to be implement!")
}

type IInterfaceMock struct {
	MakeMockFunc func() string
}

var _ IInterface = &IInterfaceMock{}

func (*IInterfaceMock) MakeMock() string {
	panic("Need to be implement!")
}

type IPackagesMock struct {
	LookupPkgFunc func(name string) (Package, bool)
}

var _ IPackages = &IPackagesMock{}

func (*IPackagesMock) LookupPkg(name string) (Package, bool) {
	panic("Need to be implement!")
}

type IPackageMock struct {
	NewGoFileWithSuffixFunc func(suffix string) (file string)

	SaveInterfaceFunc func(file string) error

	SaveMockFunc func(file string) error
}

var _ IPackage = &IPackageMock{}

func (*IPackageMock) NewGoFileWithSuffix(suffix string) (file string) {
	panic("Need to be implement!")
}

func (*IPackageMock) SaveInterface(file string) error {
	panic("Need to be implement!")
}

func (*IPackageMock) SaveMock(file string) error {
	panic("Need to be implement!")
}

type IStructMock struct {
	DemoFunc func(in types.Array) types.Basic

	MakeInterfaceFunc func() string

	PointerMethodFunc func(in types.Basic) types.Slice

	StringFunc func(f Field, ip importpath.ImportPath)

	TypeAliasFunc func(p Field, ip importpath.ImportPath)
}

var _ IStruct = &IStructMock{}

func (*IStructMock) Demo(in types.Array) types.Basic {
	panic("Need to be implement!")
}

func (*IStructMock) MakeInterface() string {
	panic("Need to be implement!")
}

func (*IStructMock) PointerMethod(in types.Basic) types.Slice {
	panic("Need to be implement!")
}

func (*IStructMock) String(f Field, ip importpath.ImportPath) {
	panic("Need to be implement!")
}

func (*IStructMock) TypeAlias(p Field, ip importpath.ImportPath) {
	panic("Need to be implement!")
}

type IPkgInfoMock struct {
	GetDirFunc func() string

	GetPkgNameFunc func() string
}

var _ IPkgInfo = &IPkgInfoMock{}

func (*IPkgInfoMock) GetDir() string {
	panic("Need to be implement!")
}

func (*IPkgInfoMock) GetPkgName() string {
	panic("Need to be implement!")
}

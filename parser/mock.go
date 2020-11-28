package parser

import (
	"go/ast"
	"go/types"

	"github.com/donnol/tools/importpath"
	"golang.org/x/tools/go/packages"
)

type PackageMock struct {
	NewGoFileWithSuffixFunc func(suffix string) (file string)

	SaveInterfaceFunc func(file string) error

	SaveMockFunc func(file string) error
}

var _ IPackage = &PackageMock{}

func (mockRecv *PackageMock) NewGoFileWithSuffix(suffix string) (file string) {
	return mockRecv.NewGoFileWithSuffixFunc(suffix)
}

func (mockRecv *PackageMock) SaveInterface(file string) error {
	return mockRecv.SaveInterfaceFunc(file)
}

func (mockRecv *PackageMock) SaveMock(file string) error {
	return mockRecv.SaveMockFunc(file)
}

type StructMock struct {
	DemoFunc func(in types.Array) types.Basic

	MakeInterfaceFunc func() string

	PointerMethodFunc func(in types.Basic) types.Slice

	StringFunc func(f Field, ip importpath.ImportPath)

	TypeAliasFunc func(p Field, ip importpath.ImportPath)
}

var _ IStruct = &StructMock{}

func (mockRecv *StructMock) Demo(in types.Array) types.Basic {
	return mockRecv.DemoFunc(in)
}

func (mockRecv *StructMock) MakeInterface() string {
	return mockRecv.MakeInterfaceFunc()
}

func (mockRecv *StructMock) PointerMethod(in types.Basic) types.Slice {
	return mockRecv.PointerMethodFunc(in)
}

func (mockRecv *StructMock) String(f Field, ip importpath.ImportPath) {
	mockRecv.StringFunc(f, ip)
}

func (mockRecv *StructMock) TypeAlias(p Field, ip importpath.ImportPath) {
	mockRecv.TypeAliasFunc(p, ip)
}

type PkgInfoMock struct {
	GetDirFunc func() string

	GetPkgNameFunc func() string
}

var _ IPkgInfo = &PkgInfoMock{}

func (mockRecv *PkgInfoMock) GetDir() string {
	return mockRecv.GetDirFunc()
}

func (mockRecv *PkgInfoMock) GetPkgName() string {
	return mockRecv.GetPkgNameFunc()
}

type InspectorMock struct {
	InspectFileFunc func(file *ast.File) (result FileResult)

	InspectPkgFunc func(pkg *packages.Package) Package
}

var _ IInspector = &InspectorMock{}

func (mockRecv *InspectorMock) InspectFile(file *ast.File) (result FileResult) {
	return mockRecv.InspectFileFunc(file)
}

func (mockRecv *InspectorMock) InspectPkg(pkg *packages.Package) Package {
	return mockRecv.InspectPkgFunc(pkg)
}

type ParserMock struct {
	GetPkgInfoFunc func() PkgInfo

	ParseASTFunc func(importPath string) (structs []Struct, err error)

	ParseByGoPackagesFunc func(patterns ...string) (result Packages, err error)
}

var _ IParser = &ParserMock{}

func (mockRecv *ParserMock) GetPkgInfo() PkgInfo {
	return mockRecv.GetPkgInfoFunc()
}

func (mockRecv *ParserMock) ParseAST(importPath string) (structs []Struct, err error) {
	return mockRecv.ParseASTFunc(importPath)
}

func (mockRecv *ParserMock) ParseByGoPackages(patterns ...string) (result Packages, err error) {
	return mockRecv.ParseByGoPackagesFunc(patterns...)
}

type visitorMock struct {
	VisitFunc func(node ast.Node) (w ast.Visitor)
}

var _ Ivisitor = &visitorMock{}

func (mockRecv *visitorMock) Visit(node ast.Node) (w ast.Visitor) {
	return mockRecv.VisitFunc(node)
}

type InterfaceMock struct {
	MakeMockFunc func() string
}

var _ IInterface = &InterfaceMock{}

func (mockRecv *InterfaceMock) MakeMock() string {
	return mockRecv.MakeMockFunc()
}

type PackagesMock struct {
	LookupPkgFunc func(name string) (Package, bool)
}

var _ IPackages = &PackagesMock{}

func (mockRecv *PackagesMock) LookupPkg(name string) (Package, bool) {
	return mockRecv.LookupPkgFunc(name)
}

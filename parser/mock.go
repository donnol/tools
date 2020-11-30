package parser

import (
	"go/ast"
	"go/types"

	"github.com/donnol/tools/importpath"
	"golang.org/x/tools/go/packages"
)

type PackagesMockMock struct {
	LookupPkgFunc func(name string) (Package, bool)
}

var _ IPackagesMock = &PackagesMockMock{}

func (mockRecv *PackagesMockMock) LookupPkg(name string) (Package, bool) {
	return mockRecv.LookupPkgFunc(name)
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

type ParserMockMock struct {
	GetPkgInfoFunc func() PkgInfo

	ParseASTFunc func(importPath string) (structs []Struct, err error)

	ParseByGoPackagesFunc func(patterns ...string) (result Packages, err error)
}

var _ IParserMock = &ParserMockMock{}

func (mockRecv *ParserMockMock) GetPkgInfo() PkgInfo {
	return mockRecv.GetPkgInfoFunc()
}

func (mockRecv *ParserMockMock) ParseAST(importPath string) (structs []Struct, err error) {
	return mockRecv.ParseASTFunc(importPath)
}

func (mockRecv *ParserMockMock) ParseByGoPackages(patterns ...string) (result Packages, err error) {
	return mockRecv.ParseByGoPackagesFunc(patterns...)
}

type visitorMock struct {
	VisitFunc func(node ast.Node) (w ast.Visitor)
}

var _ Ivisitor = &visitorMock{}

func (mockRecv *visitorMock) Visit(node ast.Node) (w ast.Visitor) {
	return mockRecv.VisitFunc(node)
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

type InspectorMockMock struct {
	InspectFileFunc func(file *ast.File) (result FileResult)

	InspectPkgFunc func(pkg *packages.Package) Package
}

var _ IInspectorMock = &InspectorMockMock{}

func (mockRecv *InspectorMockMock) InspectFile(file *ast.File) (result FileResult) {
	return mockRecv.InspectFileFunc(file)
}

func (mockRecv *InspectorMockMock) InspectPkg(pkg *packages.Package) Package {
	return mockRecv.InspectPkgFunc(pkg)
}

type InterfaceMockMock struct {
	MakeMockFunc func() string
}

var _ IInterfaceMock = &InterfaceMockMock{}

func (mockRecv *InterfaceMockMock) MakeMock() string {
	return mockRecv.MakeMockFunc()
}

type InterfaceMock struct {
	MakeMockFunc func() string
}

var _ IInterface = &InterfaceMock{}

func (mockRecv *InterfaceMock) MakeMock() string {
	return mockRecv.MakeMockFunc()
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

type StructMockMock struct {
	DemoFunc func(in types.Array) types.Basic

	MakeInterfaceFunc func() string

	PointerMethodFunc func(in types.Basic) types.Slice

	StringFunc func(f Field, ip importpath.ImportPath)

	TypeAliasFunc func(p Field, ip importpath.ImportPath)
}

var _ IStructMock = &StructMockMock{}

func (mockRecv *StructMockMock) Demo(in types.Array) types.Basic {
	return mockRecv.DemoFunc(in)
}

func (mockRecv *StructMockMock) MakeInterface() string {
	return mockRecv.MakeInterfaceFunc()
}

func (mockRecv *StructMockMock) PointerMethod(in types.Basic) types.Slice {
	return mockRecv.PointerMethodFunc(in)
}

func (mockRecv *StructMockMock) String(f Field, ip importpath.ImportPath) {
	mockRecv.StringFunc(f, ip)
}

func (mockRecv *StructMockMock) TypeAlias(p Field, ip importpath.ImportPath) {
	mockRecv.TypeAliasFunc(p, ip)
}

type visitorMockMock struct {
	VisitFunc func(node ast.Node) (w ast.Visitor)
}

var _ IvisitorMock = &visitorMockMock{}

func (mockRecv *visitorMockMock) Visit(node ast.Node) (w ast.Visitor) {
	return mockRecv.VisitFunc(node)
}

type PkgInfoMockMock struct {
	GetDirFunc func() string

	GetPkgNameFunc func() string
}

var _ IPkgInfoMock = &PkgInfoMockMock{}

func (mockRecv *PkgInfoMockMock) GetDir() string {
	return mockRecv.GetDirFunc()
}

func (mockRecv *PkgInfoMockMock) GetPkgName() string {
	return mockRecv.GetPkgNameFunc()
}

type PackagesMock struct {
	LookupPkgFunc func(name string) (Package, bool)
}

var _ IPackages = &PackagesMock{}

func (mockRecv *PackagesMock) LookupPkg(name string) (Package, bool) {
	return mockRecv.LookupPkgFunc(name)
}

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

type PackageMockMock struct {
	NewGoFileWithSuffixFunc func(suffix string) (file string)

	SaveInterfaceFunc func(file string) error

	SaveMockFunc func(file string) error
}

var _ IPackageMock = &PackageMockMock{}

func (mockRecv *PackageMockMock) NewGoFileWithSuffix(suffix string) (file string) {
	return mockRecv.NewGoFileWithSuffixFunc(suffix)
}

func (mockRecv *PackageMockMock) SaveInterface(file string) error {
	return mockRecv.SaveInterfaceFunc(file)
}

func (mockRecv *PackageMockMock) SaveMock(file string) error {
	return mockRecv.SaveMockFunc(file)
}

package parser

import (
	"fmt"
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/donnol/tools/format"
	"github.com/donnol/tools/importpath"
	"github.com/donnol/tools/internal/utils/debug"
	"golang.org/x/tools/go/packages"
)

type Packages struct {
	Patterns []string
	Pkgs     []Package
}

func (pkgs Packages) LookupPkg(name string) (Package, bool) {
	pkg := Package{}
	for _, single := range pkgs.Pkgs {
		if single.Name == name {
			return single, true
		}
	}
	return pkg, false
}

func (pkg Package) NewGoFileWithSuffix(suffix string) (file string) {
	part := strings.ReplaceAll(pkg.PkgPath, pkg.Module.Path, "")
	debug.Debug("pkg: %+v, module: %+v, %s\n", pkg.PkgPath, pkg.Module, part)

	dir := filepath.Join(pkg.Module.Dir, part)
	file = filepath.Join(dir, suffix+".go")

	return
}

func (pkg Package) SaveInterface(file string) error {
	var gocontent = "package " + pkg.Name + "\n"

	var content string
	for _, single := range pkg.Structs {
		content += single.MakeInterface() + "\n\n"
	}
	if content == "" {
		return nil
	}
	gocontent += content

	// TODO:检查是否重复

	if file == "" {
		file = pkg.NewGoFileWithSuffix("interface")
	}
	// 写入
	formatContent, err := format.Format(file, gocontent, false)
	if err != nil {
		return err
	}
	debug.Debug("content: %s, file: %s\n", formatContent, file)

	if err = ioutil.WriteFile(file, []byte(formatContent), os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (pkg Package) SaveMock(file string) error {
	var gocontent = "package " + pkg.Name + "\n"

	var content string
	for _, single := range pkg.Interfaces {
		content += single.MakeMock() + "\n\n"
	}
	if content == "" {
		return nil
	}
	gocontent += content

	// TODO:检查是否重复

	if file == "" {
		file = pkg.NewGoFileWithSuffix("mock")
	}
	// 写入
	formatContent, err := format.Format(file, gocontent, false)
	if err != nil {
		return err
	}
	debug.Debug("content: %s, file: %s\n", formatContent, file)

	if err = ioutil.WriteFile(file, []byte(formatContent), os.ModePerm); err != nil {
		return err
	}

	return nil
}

type Package struct {
	*packages.Package

	Structs    []Struct
	Interfaces []Interface
}

type Interface struct {
	*types.Interface

	PkgPath string
	PkgName string
	Name    string
	Methods []Method // 方法列表
}

func (s Interface) MakeMock() string {
	var is string
	var ms string
	for _, m := range s.Methods {
		is += fmt.Sprintf("\n%sFunc %s\n", m.Name, m.Signature)

		ms += fmt.Sprintf("\nfunc (*%s) %s%s{\n panic(\"Need to be implement!\") \n}\n", s.makeMockName(), m.Name, strings.TrimLeft(m.Signature, "func"))
	}

	is = mockPrefix(s.makeMockName(), is)

	is += `var _ ` + s.Name + ` = &` + s.makeMockName() + "{}\n"
	is += ms

	return is
}

func (s Interface) makeMockName() string {
	return s.Name + "Mock"
}

func mockPrefix(name, is string) string {
	return "type " + name + " struct{ " + is + "}\n"
}

// Struct 结构体
type Struct struct {
	// 如：github.com/pkg/errors
	PkgPath string // 包路径

	// 如: errors
	PkgName string // 包名

	Field

	Fields  []Field  // 字段列表
	Methods []Method // 方法列表
}

// --- 测试方法

// 让它传入本包里的另外一个结构体
// 传入本项目其它包的结构体
func (s Struct) String(f Field, ip importpath.ImportPath) {
	fmt.Printf("%s\n", s.PkgName)
}

func (s Struct) TypeAlias(p IIIIIIIInfo, ip ImportPathAlias) {

}

func (s Struct) Demo(in types.Array) types.Basic {
	return types.Basic{}
}

func (s *Struct) PointerMethod(in types.Basic) types.Slice {
	return types.Slice{}
}

// --- 测试方法

// MMakeInterface 根据结构体的方法生成相应接口
func (s Struct) MakeInterface() string {
	methods := make([]*types.Func, 0, len(s.Methods))
	for _, m := range s.Methods {
		if !m.Origin.Exported() {
			continue
		}
		methods = append(methods, m.Origin)
		// fmt.Printf("method: %+v, %s\n", m.Origin, m.Signature)
	}

	if len(methods) == 0 {
		return ""
	}

	i := types.NewInterfaceType(methods, nil)
	i = i.Complete()
	is := types.TypeString(i.Underlying(), pkgNameQualifier(qualifierParam{pkgPath: s.PkgPath}))

	is = interfacePrefix(s.makeInterfaceName(), is)

	return is
}

func (s Struct) makeInterfaceName() string {
	return "I" + s.Name
}

func interfacePrefix(name, is string) string {
	return "type " + name + " " + is
}

type Method struct {
	Origin    *types.Func
	Name      string
	Signature string
}

// Field 字段
type Field struct {
	Id   string // 唯一标志
	Name string // 名称

	TypesType types.Type // 原始类型
	Type      string     // 类型，包含包导入路径

	Doc     string // 文档
	Comment string // 注释
}

// IIIIIIIInfo 别名测试
type IIIIIIIInfo = Field // 别名测试注释

type ImportPathAlias = importpath.ImportPath

type PkgInfo struct {
	dir     string
	pkgName string
}

func (i PkgInfo) GetDir() string {
	return i.dir
}

func (i PkgInfo) GetPkgName() string {
	return i.pkgName
}

type FileResult struct {
	structMap    map[string]Struct    // 名称 -> 结构体
	methodMap    map[string][]Method  // 名称 -> 方法列表
	interfaceMap map[string]Interface // 名称 -> 接口
}

func MakeFileResult() FileResult {
	return FileResult{
		structMap:    make(map[string]Struct),
		methodMap:    make(map[string][]Method),
		interfaceMap: make(map[string]Interface),
	}
}

type DeclResult struct {
	structMap    map[string]Struct
	methodMap    map[string][]Method
	interfaceMap map[string]Interface // 名称 -> 接口
}

func MakeDeclResult() DeclResult {
	return DeclResult{
		structMap:    make(map[string]Struct),
		methodMap:    make(map[string][]Method),
		interfaceMap: make(map[string]Interface),
	}
}

type SpecResult struct {
	structMap    map[string]Struct    // 名称 -> 结构体
	interfaceMap map[string]Interface // 名称 -> 接口
}

func MakeSpecResult() SpecResult {
	return SpecResult{
		structMap:    make(map[string]Struct),
		interfaceMap: make(map[string]Interface),
	}
}

type ExprResult struct {
	Fields []Field
}

func MakeExprResult() ExprResult {
	return ExprResult{
		Fields: make([]Field, 0),
	}
}

type StmtResult struct {
}

func MakeStmtResult() StmtResult {
	return StmtResult{}
}

type FieldResult struct {
	RecvName string

	Fields []Field
}

func MakeFieldResult() FieldResult {
	return FieldResult{
		Fields: make([]Field, 0),
	}
}

type TokenResult struct {
}

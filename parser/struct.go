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
	mockType := s.makeMockName()
	mockRecv := s.makeMockRecv()

	cc := fmt.Sprintf(`%sCommonProxyContext`, LcFirst(mockType))

	var is string
	var pc string
	var ms string
	for _, m := range s.Methods {
		fieldName, fieldType, methodSig, returnStmt, call := s.processFunc(m)

		is += fmt.Sprintf("\n%s %s\n", fieldName, fieldType)

		pc += fmt.Sprintf(`%s%sProxyContext = func() (pctx inject.ProxyContext) { 
			pctx = %s
			pctx.MethodName = "%s"
			return
		} () 
		`, mockType, m.Name, cc, m.Name)

		ms += fmt.Sprintf("\nfunc (%s *%s) %s {\n %s %s.%s \n}\n", mockRecv, mockType, methodSig, returnStmt, mockRecv, call)
	}

	is = mockPrefix(mockType, is)

	is += `var (_ ` + s.Name + ` = &` + mockType + "{}\n\n"
	is += fmt.Sprintf(`%s = inject.ProxyContext {
		PkgPath: "%s",
		InterfaceName: "%s",
	}
	`, cc, s.PkgPath, s.Name)
	is += pc + `)`
	is += ms

	debug.Debug("is: %s\n", is)

	return is
}

const (
	sep         = ","
	leftParent  = "("
	rightParent = ")"
)

// func(ctx context.Context, m M) (err error) -> (ctx, m)
// func(context.Context,M) (error) -> (p0, p1)
func (s Interface) processFunc(m Method) (fieldName, fieldType, methodSig, returnStmt, call string) {

	fieldName = m.Name + "Func"
	fieldType = m.Signature

	sigType := m.Origin.Type().(*types.Signature)
	if sigType.Variadic() {
		//  在这里获取完整签名字符串时，还是正常的：func(interface{}, string, ...interface{}) error
		typStr := types.TypeString(sigType, pkgNameQualifier(qualifierParam{pkgPath: s.PkgPath}))
		debug.Debug("typ: %+v, str: %s\n", sigType, typStr)
	}
	params := sigType.Params()
	for i := 0; i < params.Len(); i++ {
		pvar := params.At(i)
		name := pvar.Name()

		// 参数名可能为空，需要置默认值
		if name == "" {
			name = fmt.Sprintf("p%d", i)
		}

		// 解析进来之后，不定参数类型变成了slice：[]interface{}
		typStr := types.TypeString(pvar.Type(), pkgNameQualifier(qualifierParam{pkgPath: s.PkgPath}))

		// 处理最后一个是不定参数的情况
		var paramTypePrefix string
		if sigType.Variadic() && i == params.Len()-1 {
			paramTypePrefix = "..."
			debug.Debug("typ: %+v, str: %s, params: %v\n", pvar.Type(), typStr, params.String())
		}

		// FIXME:感觉不太好，怎么办呢？
		// 当是不定参数，typStr会从...interface{}变为[]interface{}，因此，需要再将它重新变回来
		if paramTypePrefix != "" && strings.Index(typStr, "[]") == 0 {
			typStr = typStr[2:]
		}
		methodSig += name + " " + paramTypePrefix + typStr + sep

		call += name + paramTypePrefix + sep
	}
	methodSig = strings.TrimRight(methodSig, sep)
	methodSig = m.Name + leftParent + methodSig + rightParent

	res := sigType.Results()
	returnStmt = "return"
	if res.Len() == 0 {
		returnStmt = " "
	}
	var resString string
	for i := 0; i < res.Len(); i++ {
		rvar := res.At(i)
		name := rvar.Name()

		resString += name + " " + types.TypeString(rvar.Type(), pkgNameQualifier(qualifierParam{pkgPath: s.PkgPath})) + sep
	}
	resString = strings.TrimRight(resString, sep)
	resString = leftParent + resString + rightParent
	methodSig = methodSig + resString

	debug.Debug("methodSig: %v\n", methodSig)

	call = strings.TrimRight(call, sep)
	call = leftParent + call + rightParent
	call = fieldName + call

	return
}

func (s Interface) makeMockName() string {
	name := s.Name
	// 如果首个字符是I，则去掉
	index := strings.Index(name, "I")
	if index == 0 {
		name = name[1:]
	}
	return name + "Mock"
}

func (s Interface) makeMockRecv() string {
	return "mockRecv"
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
	// 考虑到结构体名称是非导出有后缀的，如：fileImpl
	// 1. 针对非导出，将首字母变大写
	// 2. 针对impl后缀，直接去掉
	name := s.Name
	name = strings.Title(name)
	index := strings.Index(name, "Impl")
	if index != -1 {
		name = name[:index]
	}
	return "I" + name
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

package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

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

func (pkg Package) NewGoFileWithSuffix(mode, dir, suffix string) (file string) {
	if pkg.Module == nil {
		fmt.Printf("pkg.Module is nil\n")
	}
	part := strings.ReplaceAll(pkg.PkgPath, pkg.Module.Path, "")
	debug.Printf("pkg: %+v, module: %+v, %s\n", pkg.PkgPath, pkg.Module, part)

	var mockdir string
	if mode == "offsite" {
		mockdir = dir
	} else {
		mockdir = filepath.Join(pkg.Module.Dir, part)
	}
	file = filepath.Join(mockdir, suffix+".go")

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
		file = pkg.NewGoFileWithSuffix("", "", "interface")
	}
	// 写入
	formatContent, err := format.Format(file, gocontent, false)
	if err != nil {
		return err
	}
	debug.Printf("content: %s, file: %s\n", formatContent, file)

	if err = ioutil.WriteFile(file, []byte(formatContent), os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (pkg Package) SaveMock(mode, dir, file string) error {
	pkgName := pkg.Name
	if mode == "offsite" {
		pkgName = filepath.Base(dir)
	}
	var gocontent = "package " + pkgName + "\n"

	// 找出所有外部包引用，生成import
	// 因为是生成mock结构体，所以有包引用的都是参数和返回值
	imports := make(map[string]struct{}, 4)

	debug.Printf("===test\n")
	var content string
	for _, single := range pkg.Interfaces {
		debug.Printf("have type set: %+v, embeds: %d\n", single.Interface, single.Interface.NumEmbeddeds())
		if single.Interface.NumEmbeddeds() != 0 {
			log.Printf("have type set: %+v\n", single.Interface)
			continue
		}
		mock, imps := single.MakeMock(mode)
		for imp := range imps {
			imports[imp] = struct{}{}
		}
		content += mock + "\n\n"
	}
	if content == "" {
		return nil
	}

	// 全局变量/函数
	globalContent := ``
	content = globalContent + content

	// 导入
	var impcontent string
	for imp := range imports {
		if imp == "" {
			continue
		}
		impcontent += `"` + imp + `"` + "\n"
	}
	if impcontent != "" {
		impcontent = "import (\n" + impcontent + ")\n"
		debug.Printf("import: %s\n", impcontent)
	}

	gocontent += impcontent

	// mock
	gocontent += content
	debug.Printf("gocontent: %s\n", gocontent)

	// TODO:检查是否重复

	if file == "" {
		file = pkg.NewGoFileWithSuffix(mode, dir, "mock")
	}
	// 写入
	formatContent, err := format.Format(file, gocontent, false)
	if err != nil {
		return fmt.Errorf("format failed: %w, content: \n%s", err, gocontent)
	}
	debug.Printf("content: %s, file: %s\n", formatContent, file)

	if err = ioutil.WriteFile(file, []byte(formatContent), os.ModePerm); err != nil {
		return err
	}

	return nil
}

type Package struct {
	*packages.Package

	Funcs      []Func
	Structs    []Struct
	Interfaces []Interface
}

type Interface struct {
	*types.Interface

	PkgPath    string
	PkgName    string
	Name       string
	TypeParams *ast.FieldList
	Methods    []Method // 方法列表
}

var (
	proxyMethodTmpl = `
	{{.methodName}}: {{.funcSignature}} {
		var _gen_ctx = {{.mockType}}{{.funcName}}ProxyContext

		_gen_stop := do.ProxyTraceBegin(_gen_ctx{{if .argNames}}, {{.argNames}} {{end}})
		defer func() {
			_gen_stop()
		}()

		{{.funcResult}}
		
		var _gen_actual_cf do.ProxyCtxFunc

		_gen_inner_cf, _gen_inner_ok := _gen_innerCtxMap.Lookup(_gen_ctx, typeParams...)
		_gen_cf, _gen_ok := do.GlobalProxyCtxMap().Lookup(_gen_ctx, typeParams...)
		if _gen_inner_ok {
			_gen_actual_cf = _gen_inner_cf
		} else if _gen_ok {
			_gen_actual_cf = _gen_cf
		}

		if _gen_actual_cf != nil {
			_gen_params := []any{}
			{{.params}}
			{{if .funcResultList}} _gen_res := {{else}}  {{end}} _gen_actual_cf(_gen_ctx, base.{{.funcName}}, _gen_params)
			{{.resultAssert}}
		} else {
			{{if .funcResultList}} {{.funcResultList}} = {{else}}  {{end}} base.{{.funcName}}({{.argNames}})
		}

		{{if .funcResultList}} return {{.funcResultList}} {{else}}  {{end}}
	},
	`

	proxyMethodParamsTmpl = `
	{{range $index, $ele := .args}}
	_gen_params = append(_gen_params, {{.Name}}){{end}}
	`
	proxyMethodResultTmpl = `
	 {{range $index, $ele := .reses}}
	 var _gen_r{{$index}} {{.Typ}}{{end}}
	`

	proxyMethodResultAssertTmpl = `
	{{range $index, $ele := .reses}}
			_gen_tmpr{{$index}}, _gen_exist := _gen_res[{{$index}}].({{.Typ}})
			if _gen_exist {
				_gen_r{{$index}} = _gen_tmpr{{$index}}
			}{{end}}
	`
)

type arg struct {
	Name     string
	Typ      string
	Variadic bool
}

func (s Interface) MakeMock(mode string) (string, map[string]struct{}) {
	if mode == "offsite" {
		// 异地模式，忽略非导出接口
		if !ast.IsExported(s.Name) {
			return "", nil
		}
	}
	fullTypeParam, partTypeParam := s.handleTypeParams()
	debug.Printf("interface : %#v, %#v, %v, %v\n", s.Interface, s.Name, fullTypeParam, partTypeParam)
	mockType := s.makeMockName()
	mockRecv := s.makeMockRecv()
	proxyFuncName := s.makeProxyFuncName()
	debug.Printf("proxyfuncname:%s\n", proxyFuncName)

	originTypeName := s.Name
	if mode == "offsite" {
		originTypeName = s.PkgName + "." + s.Name
	}
	proxyFunc := `
	// ` + proxyFuncName + ` 获取接口代理；若使用泛型，需传入typeParams，其值为类型参数的字符串字面量；若想进一步修改方法行为，可以使用RegisterProxyMethod函数注入自定义方法实现；如果想要为每个实例单独注入方法，则使用第二个返回值对象来设置
	func ` + proxyFuncName + fullTypeParam + "(base " + originTypeName + partTypeParam + ", typeParams ...string) (" + originTypeName + partTypeParam + ",  *do.ProxyCtxFuncStore) {" + `if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	_gen_innerCtxMap := do.NewProxyCtxMap()
	return &` + mockType + partTypeParam + `{`
	cc := fmt.Sprintf(`%sCommonProxyContext`, LcFirst(mockType))

	var is string
	var pc string
	var allpc = mockType + `ProxyContextAll = []do.ProxyContext{`
	var ms string
	var proxyMethod = new(bytes.Buffer)
	var imports = make(map[string]struct{}, 4)
	for _, m := range s.Methods {
		fieldName, fieldType, methodSig, returnStmt, call, args, reses, imps := s.processFunc(mode, m)

		for imp := range imps {
			imports[imp] = struct{}{}
		}

		is += fmt.Sprintf("\n%s %s\n", fieldName, fieldType)

		pcname := fmt.Sprintf("%s%sProxyContext", mockType, m.Name)
		allpc += `
		` + pcname + ","
		pc += fmt.Sprintf(`
		// represent %s.%s: %s
		%s = func() (pctx do.ProxyContext) { 
			pctx = %s
			pctx.MethodName = "%s"
			return
		} () 
		`, s.Name, m.Name, fieldType, pcname, cc, m.Name)

		ms += fmt.Sprintf("\nfunc (%s *%s) %s {\n %s %s.%s \n}\n", mockRecv, mockType+partTypeParam, methodSig, returnStmt, mockRecv, call)

		assertBuf := new(bytes.Buffer)
		assertTmpl, err := template.New("proxyMethodResultAssert").Parse(proxyMethodResultAssertTmpl)
		if err != nil {
			panic(err)
		}
		assertTmpl.Execute(assertBuf, map[string]interface{}{
			"reses": reses,
		})
		paramBuf := new(bytes.Buffer)
		paramTmpl, err := template.New("proxyMethodParam").Parse(proxyMethodParamsTmpl)
		if err != nil {
			panic(err)
		}
		methodArgs := make([]arg, 0, len(args))
		for _, arg := range args {
			// 这里把不定参数的...去掉，生成的代码会以切片的形式传递参数
			i := strings.Index(arg.Name, "...")
			if i != -1 {
				arg.Name = arg.Name[:i]
			}
			methodArgs = append(methodArgs, arg)
		}
		paramTmpl.Execute(paramBuf, map[string]interface{}{
			"args": methodArgs,
		})
		argNames := ""
		for i, arg := range args {
			argNames += arg.Name
			if i != len(args)-1 {
				argNames += ", "
			}
		}
		resBuf := new(bytes.Buffer)
		resTmpl, err := template.New("proxyMethodResult").Parse(proxyMethodResultTmpl)
		if err != nil {
			panic(err)
		}
		resTmpl.Execute(resBuf, map[string]interface{}{
			"reses": reses,
		})
		funcResultList := ""
		for i := range reses {
			funcResultList += "_gen_r" + strconv.Itoa(i)
			if i != len(reses)-1 {
				funcResultList += ", "
			}
		}
		tmpl, err := template.New("proxyMethod").Parse(proxyMethodTmpl)
		if err != nil {
			panic(err)
		}
		tmpl.Execute(proxyMethod, map[string]interface{}{
			"methodName":     fieldName,
			"funcSignature":  strings.Replace(methodSig, m.Name, "func", 1),
			"mockType":       mockType,
			"funcName":       m.Name,
			"funcResult":     resBuf.String(),
			"funcResultList": funcResultList,
			"argNames":       argNames,
			"params":         paramBuf.String(),
			"resultAssert":   assertBuf.String(),
		})
	}

	proxyFunc += proxyMethod.String() + "}, _gen_innerCtxMap}"
	is = mockStructPrefix(mockType+fullTypeParam, is)

	is += `var (`
	is += fmt.Sprintf(`%s = do.ProxyContext {
		PkgPath: "%s",
		InterfaceName: "%s",
	}
	`, cc, s.PkgPath, s.Name)
	is += pc + "\n" + allpc + "\n}\n" + `)`
	is += "\n" + "\n\n" + proxyFunc + "\n"
	is += ms

	debug.Printf("is: %s\n", is)

	return is, imports
}

const (
	sep         = ","
	leftParent  = "("
	rightParent = ")"
)

// func(ctx context.Context, m M) (err error) -> (ctx, m)
// func(context.Context,M) (error) -> (p0, p1)
func (s Interface) processFunc(mode string, m Method) (fieldName, fieldType, methodSig, returnStmt, call string, args []arg, reses []arg, imports map[string]struct{}) {

	imports = make(map[string]struct{}, 4) // 导入的包
	fieldName = m.Name + "Func"
	fieldType = m.Signature

	sigType := m.Origin.Type().(*types.Signature)
	if sigType.Variadic() {
		//  在这里获取完整签名字符串时，还是正常的：func(interface{}, string, ...interface{}) error
		typStr := types.TypeString(sigType, pkgNameQualifier(qualifierParam{pkgPath: s.PkgPath, keepPkgPathWhenIsSamePkg: mode == "offsite"}))
		debug.Printf("typ: %+v, str: %s\n", sigType, typStr)
	}
	params := sigType.Params()
	for i := 0; i < params.Len(); i++ {
		pvar := params.At(i)
		name := pvar.Name()

		// 参数名可能为空，需要置默认值
		if name == "" {
			name = fmt.Sprintf("p%d", i)
		}

		// 参数类型的包路径信息
		pkgPath := getTypesPkgPath(pvar.Type())
		imports[pkgPath] = struct{}{}

		// 解析进来之后，不定参数类型变成了slice：[]interface{}
		typStr := types.TypeString(pvar.Type(), pkgNameQualifier(qualifierParam{pkgPath: s.PkgPath, keepPkgPathWhenIsSamePkg: mode == "offsite"}))

		// 处理最后一个是不定参数的情况
		var paramTypePrefix string
		var variadic bool
		if sigType.Variadic() && i == params.Len()-1 {
			paramTypePrefix = "..."
			variadic = true
			debug.Printf("typ: %+v, str: %s, params: %v\n", pvar.Type(), typStr, params.String())
		}

		// FIXME:感觉不太好，怎么办呢？
		// 当是不定参数，typStr会从...interface{}变为[]interface{}，因此，需要再将它重新变回来
		if paramTypePrefix != "" && strings.Index(typStr, "[]") == 0 {
			typStr = typStr[2:]
		}
		methodSig += name + " " + paramTypePrefix + typStr
		if i != params.Len()-1 {
			methodSig += sep + " "
		}

		call += name + paramTypePrefix + sep

		args = append(args, arg{Name: name + paramTypePrefix, Typ: typStr, Variadic: variadic})
	}
	methodSig = strings.TrimRight(methodSig, sep)
	methodSig = m.Name + leftParent + methodSig + rightParent

	res := sigType.Results()
	returnStmt = "return"
	if res.Len() == 0 {
		returnStmt = " "
	}
	useNamedRet := false
	var resString string
	for i := 0; i < res.Len(); i++ {
		rvar := res.At(i)
		name := rvar.Name()
		if name != "" {
			useNamedRet = true
		}

		// 返回类型的包路径信息
		pkgPath := getTypesPkgPath(rvar.Type())
		imports[pkgPath] = struct{}{}

		typ := types.TypeString(rvar.Type(), pkgNameQualifier(qualifierParam{pkgPath: s.PkgPath, keepPkgPathWhenIsSamePkg: mode == "offsite"}))
		if name == "" {
			resString += typ
		} else {
			resString += name + " " + typ
		}
		if i != res.Len()-1 {
			resString += sep + " "
		}

		reses = append(reses, arg{
			Name: name,
			Typ:  typ,
		})
	}
	resString = strings.TrimRight(resString, sep)
	if res.Len() > 1 || useNamedRet {
		resString = leftParent + resString + rightParent
	}
	methodSig = methodSig + " " + resString

	debug.Printf("methodSig: %v\n", methodSig)
	if mode == "offsite" {
		li := strings.Index(methodSig, leftParent)
		fieldType = "func" + methodSig[li:]
	}

	call = strings.TrimRight(call, sep)
	call = leftParent + call + rightParent
	call = fieldName + call

	return
}

func (s Interface) makeProxyFuncName() string {
	return "Get" + s.Name + "Proxy"
}

func (s Interface) makeMockName() string {
	return s.Name + "Mock"
}

func (s Interface) handleTypeParams() (full string, part string) {
	if s.TypeParams == nil {
		return
	}
	for i, tp := range s.TypeParams.List {
		if len(tp.Names) == 0 {
			continue
		}
		for _, name := range tp.Names {
			full += name.Name
			part += name.Name
		}
		full += " " + types.ExprString(tp.Type)
		if i != len(s.TypeParams.List)-1 {
			full += ", "
			part += ", "
		}
	}
	full = "[" + full + "]"
	part = "[" + part + "]"
	return
}

func (s Interface) RemoveFirst(c string) string {
	name := s.Name
	// 如果首个字符是c，则去掉
	index := strings.Index(name, c)
	if index == 0 {
		name = name[1:]
	}
	return name
}

func (s Interface) makeMockRecv() string {
	return "mockRecv"
}

func mockStructPrefix(name, is string) string {
	return `
	// ===== ` + name + ` =====

	type ` + name + " struct{ " + is + "}\n"
}

// Struct 结构体
type Struct struct {
	// 如：github.com/pkg/errors
	PkgPath string `json:"pkgPath" toml:"pkg_path"` // 包路径

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

type Method = Func

type Func struct {
	Origin *types.Func

	PkgPath   string // 包路径
	Recv      string // 方法的receiver
	Name      string
	Signature string

	Calls []Func // 调用的函数/方法
}

func (f *Func) Set(fm map[string]Func, depth int) {
	l := 1
	setLowerCalls(f.Calls, fm, l, depth)
}

func setLowerCalls(calls []Func, fm map[string]Func, l, depth int) {
	if l > depth {
		return
	}
	for i, call := range calls {
		var key = call.Name
		if call.Recv != "" {
			key = call.Recv + "." + call.Name
		}
		if len(call.Calls) == 0 {
			calls[i].Calls = fm[key].Calls
			nl := l + 1
			setLowerCalls(calls[i].Calls, fm, nl, depth)
		}
	}
}

// PrintCallGraph 打印调用图，用ignore忽略包，用depth指定深度
func (f Func) PrintCallGraph(ignore []string, depth int) {
	ip := &importpath.ImportPath{}
	curPath, err := ip.GetByCurrentDir()
	if err != nil {
		panic(err) // 怎么知道这些内置函数是内置函数呢？
	}
	fmt.Printf("root module path: %s\n", curPath)

	fmt.Printf("root: %s(%s)\n", f.Name, f.PkgPath)
	l := 1

	printCallGraph(f.Calls, ignore, l, depth)
}

func printCallGraph(calls []Func, ignores []string, l, depth int) {
	for _, one := range calls {
		if l > depth {
			break
		}

		// 判断是否需要跳过
		pkgPath := one.PkgPath
		needIgnore := false
		for _, ignore := range ignores {
			if pkgPath != "" && ignore == pkgPath {
				needIgnore = true
				break
			}
		}
		if needIgnore {
			continue
		}

		fmt.Printf("%s -> %s(%s)\n", getIdent(l), one.Name, one.PkgPath)

		if len(one.Calls) > 0 {
			nl := l + 1
			printCallGraph(one.Calls, ignores, nl, depth)
		}
	}
}

const (
	ident = "	"
)

func getIdent(l int) string {
	s := ""
	for i := 0; i < l; i++ {
		if i == l-1 {
			s += "   " + strconv.Itoa(l)
		} else {
			s += ident
		}
	}
	return s
}

// Field 字段
type Field struct {
	Id        string // 唯一标志
	Name      string // 名称
	Anonymous bool   // 是否匿名

	TypesType types.Type // 原始类型
	Type      string     // 类型，包含包导入路径

	Tag         string        `json:"tag"` // 结构体字段的tag
	TagBasicLit *ast.BasicLit // ast的tag类型

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
	funcMap      map[string]Func      // 名称 -> 方法
}

func MakeFileResult() FileResult {
	return FileResult{
		structMap:    make(map[string]Struct),
		methodMap:    make(map[string][]Method),
		interfaceMap: make(map[string]Interface),
		funcMap:      make(map[string]Func),
	}
}

type DeclResult struct {
	structMap    map[string]Struct
	methodMap    map[string][]Method
	interfaceMap map[string]Interface // 名称 -> 接口
	funcMap      map[string]Func      // 名称 -> 方法
}

func MakeDeclResult() DeclResult {
	return DeclResult{
		structMap:    make(map[string]Struct),
		methodMap:    make(map[string][]Method),
		interfaceMap: make(map[string]Interface),
		funcMap:      make(map[string]Func),
	}
}

type SpecResult struct {
	structMap    map[string]Struct    // 名称 -> 结构体
	interfaceMap map[string]Interface // 名称 -> 接口
	funcMap      map[string]Func      // 名称 -> 方法
}

func MakeSpecResult() SpecResult {
	return SpecResult{
		structMap:    make(map[string]Struct),
		interfaceMap: make(map[string]Interface),
		funcMap:      make(map[string]Func),
	}
}

type ExprResult struct {
	Fields  []Field
	pkgPath string
	funcMap map[string]Func // 名称 -> 方法
}

func MakeExprResult() ExprResult {
	return ExprResult{
		Fields:  make([]Field, 0),
		funcMap: make(map[string]Func),
	}
}

func (er ExprResult) Merge(oer ExprResult) (ner ExprResult) {
	ner = er

	if ner.pkgPath == "" && oer.pkgPath != "" {
		ner.pkgPath = oer.pkgPath
	}
	ner.Fields = append(ner.Fields, oer.Fields...)
	for k, v := range oer.funcMap {
		ner.funcMap[k] = v
	}

	return
}

type StmtResult struct {
	pkgPath string
	funcMap map[string]Func // 名称 -> 方法
}

func MakeStmtResult() StmtResult {
	return StmtResult{
		funcMap: make(map[string]Func),
	}
}

func (er StmtResult) Merge(oer StmtResult) (ner StmtResult) {
	ner = er

	if ner.pkgPath == "" && oer.pkgPath != "" {
		ner.pkgPath = oer.pkgPath
	}
	for k, v := range oer.funcMap {
		ner.funcMap[k] = v
	}

	return
}

func (er StmtResult) MergeExprResult(oer ExprResult) (ner StmtResult) {
	ner = er

	if ner.pkgPath == "" && oer.pkgPath != "" {
		ner.pkgPath = oer.pkgPath
	}
	for k, v := range oer.funcMap {
		ner.funcMap[k] = v
	}

	return
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

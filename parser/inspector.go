package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/donnol/tools/internal/utils/debug"
	"golang.org/x/tools/go/packages"
)

type Inspector struct {
	parser *Parser

	pkg *packages.Package
}

type InspectOption struct {
	Parser *Parser
}

func NewInspector(opt InspectOption) *Inspector {
	return &Inspector{
		parser: opt.Parser,
	}
}

func (ins *Inspector) InspectPkg(pkg *packages.Package) Package {
	if pkg == nil {
		panic("input pkg is nil")
	}
	ins.pkg = pkg

	// 解析*ast.File信息
	structMap := make(map[string]Struct)
	methodsMap := make(map[string][]Method)
	interfaceMap := make(map[string]Interface)
	for i, astFile := range pkg.Syntax {
		// 替换import path
		if ins.parser.replaceImportPath {
			fileName := pkg.CompiledGoFiles[i]
			fmt.Printf("%v\n", pkg.CompiledGoFiles)
			if err := ins.parser.replaceFileImportPath(fileName, astFile); err != nil {
				panic(fmt.Errorf("replaceFileImportPath failed: %+v", err))
			}
			continue
		}

		fileResult := ins.InspectFile(astFile)

		for k, v := range fileResult.structMap {
			structMap[k] = v
		}
		for k, v := range fileResult.methodMap {
			methodsMap[k] = append(methodsMap[k], v...)
		}
		for k, v := range fileResult.interfaceMap {
			interfaceMap[k] = v
		}
	}

	structs := make([]Struct, 0, len(structMap))
	for structName, single := range structMap {
		methods := methodsMap[structName]
		single.Methods = methods
		structs = append(structs, single)
	}
	inters := make([]Interface, 0, len(interfaceMap))
	for _, single := range interfaceMap {
		inters = append(inters, single)
	}

	return Package{
		Package:    pkg,
		Structs:    structs,
		Interfaces: inters,
	}
}

func (ins *Inspector) InspectFile(file *ast.File) (result FileResult) {
	if file == nil {
		return
	}
	result = MakeFileResult()

	structMap := make(map[string]Struct)
	methodsMap := make(map[string][]Method)
	interfaceMap := make(map[string]Interface)
	for _, decl := range file.Decls {
		declResult := ins.inspectDecl(decl)
		for k, v := range declResult.structMap {
			structMap[k] = v
		}
		for k, v := range declResult.methodMap {
			methodsMap[k] = append(methodsMap[k], v...)
		}
		for k, v := range declResult.interfaceMap {
			interfaceMap[k] = v
		}
	}
	result.structMap = structMap
	result.methodMap = methodsMap
	result.interfaceMap = interfaceMap

	return
}

func (ins *Inspector) inspectDecl(decl ast.Decl) (result DeclResult) {
	if decl == nil {
		return
	}
	result = MakeDeclResult()

	switch declValue := decl.(type) {
	case *ast.BadDecl:
		panic(fmt.Errorf("BadDecl: %+v", declValue))

	case *ast.FuncDecl:
		debug.Debug("FundDecl name: %s, %s\n", declValue.Name, declValue.Doc.Text())

		funcType := &types.Func{}
		obj := ins.pkg.TypesInfo.Defs[declValue.Name]
		switch objTyp := obj.Type().(type) {
		case *types.Signature:
			debug.Debug("objTyp sig: %+v, %s\n", objTyp, toString(objTyp))
			funcType = types.NewFunc(declValue.Type.Func, ins.pkg.Types, obj.Name(), objTyp)
		}
		method := Method{
			Origin:    funcType,
			Name:      obj.Name(),
			Signature: toString(obj.Type()),
		}
		debug.Debug("method: %+v\n", method)

		ins.inspectExpr(declValue.Type) // 函数签名
		ins.inspectStmt(declValue.Body) // 函数体

		// method receiver: func (x *X) XXX()里的(x *X)部分
		recvName := ""
		if declValue.Recv != nil {
			debug.Debug("FundDecl recv: %v\n", declValue.Recv.List)

			fieldResult := ins.inspectFields(declValue.Recv)
			recvName = fieldResult.RecvName
		}
		result.methodMap[recvName] = append(result.methodMap[recvName], method)

	case *ast.GenDecl:
		switch declValue.Tok {
		case token.IMPORT:
		case token.CONST:
		case token.VAR:
		case token.TYPE:
			for _, spec := range declValue.Specs {
				specResult := ins.inspectSpec(spec)
				for k, v := range specResult.structMap {
					result.structMap[k] = v
				}
				for k, v := range specResult.interfaceMap {
					result.interfaceMap[k] = v
				}
			}
		}
	}

	return
}

func (ins *Inspector) inspectSpec(spec ast.Spec) (result SpecResult) {
	if spec == nil {
		return
	}
	result = MakeSpecResult()

	switch specValue := spec.(type) {
	case *ast.ImportSpec:
		debug.Debug("ImportSpec, name: %v, path: %v\n", specValue.Name, specValue.Path)

	case *ast.ValueSpec:
		debug.Debug("ValueSpec, name: %+v, type: %+v, value: %+v\n", specValue.Names, specValue.Type, specValue.Values)

	case *ast.TypeSpec:
		// 这里拿到类型信息: 名称，注释，文档
		debug.Debug("TypeSpec name: %s, type: %+v, comment: %s, doc: %s\n", specValue.Name, specValue.Type, specValue.Comment.Text(), specValue.Doc.Text())

		switch specValue.Type.(type) {
		case *ast.InterfaceType:
			exprResult := ins.inspectExpr(specValue.Type)
			debug.Debug("interface type name: %s, exprValue: %+v, type: %+v, result: %+v\n", specValue.Name, specValue, specValue.Type, exprResult)

			interType := ins.pkg.TypesInfo.TypeOf(specValue.Type)
			r := parseTypesType(interType, parseTypesTypeOption{pkgPath: ins.pkg.PkgPath})
			methods := r.methods

			inter := Interface{
				Interface: ins.pkg.TypesInfo.Types[specValue.Type].Type.(*types.Interface),
				Name:      specValue.Name.Name,
				PkgPath:   ins.pkg.PkgPath,
				PkgName:   ins.pkg.Name,
				Methods:   methods,
			}
			debug.Debug("mock: %s\n", inter.MakeMock())
			result.interfaceMap[specValue.Name.Name] = inter

		default:
			structOne := Struct{
				PkgPath: ins.pkg.PkgPath,
				PkgName: ins.pkg.Name,
				Field: Field{
					Id:        ins.pkg.TypesInfo.Types[specValue.Type].Type.String(),
					Name:      specValue.Name.Name,
					TypesType: ins.pkg.TypesInfo.Types[specValue.Type].Type,
					Type:      toString(specValue.Type),
					Doc:       specValue.Doc.Text(),
					Comment:   specValue.Comment.Text(),
				},
			}

			// 再拿field
			exprResult := ins.inspectExpr(specValue.Type)
			structOne.Fields = exprResult.Fields
			result.structMap[specValue.Name.Name] = structOne
		}
	}

	return
}

func (ins *Inspector) inspectExpr(expr ast.Expr) (result ExprResult) {
	if expr == nil {
		return
	}
	result = MakeExprResult()

	switch exprValue := expr.(type) {
	case *ast.StructType:
		fieldResult := ins.inspectFields(exprValue.Fields)
		result.Fields = fieldResult.Fields

	case *ast.StarExpr: // *T
		ins.inspectExpr(exprValue.X)
	case *ast.TypeAssertExpr: // X.(*T)
		ins.inspectExpr(exprValue.X)
		ins.inspectExpr(exprValue.Type)
	case *ast.ArrayType: // [L]T
		ins.inspectExpr(exprValue.Len)
		ins.inspectExpr(exprValue.Elt)
	case *ast.BadExpr:
		panic(fmt.Errorf("BadExpr: %+v", exprValue))
	case *ast.SelectorExpr: // X.M
		debug.Debug("SelectorExpr value: %v, typesString: %s\n", exprValue, toString(exprValue))
		ins.inspectExpr(exprValue.X)
	case *ast.SliceExpr: // []T, slice[1:3:5]
		ins.inspectExpr(exprValue.X)
		ins.inspectExpr(exprValue.Low)
		ins.inspectExpr(exprValue.High)
		ins.inspectExpr(exprValue.Max)
	case *ast.BasicLit: // 33 40.0 0x1f
	case *ast.BinaryExpr: // X+Y X-Y X*Y X/Y X%Y
		ins.inspectExpr(exprValue.X)
		ins.inspectExpr(exprValue.Y)
	case *ast.CallExpr: // M(1, 2)
		ins.inspectExpr(exprValue.Fun)
		for _, arg := range exprValue.Args {
			ins.inspectExpr(arg)
		}
	case *ast.ChanType: // chan T, <-chan T, chan<- T
		ins.inspectExpr(exprValue.Value)
	case *ast.CompositeLit: // T{Name: Value}
		ins.inspectExpr(exprValue.Type)
		for _, elt := range exprValue.Elts {
			ins.inspectExpr(elt)
		}
	case *ast.Ellipsis: // ...int, [...]Arr
		ins.inspectExpr(exprValue.Elt)
	case *ast.FuncLit:
		ins.inspectExpr(exprValue.Type)
		ins.inspectStmt(exprValue.Body)
	case *ast.FuncType:
		ins.inspectFields(exprValue.Params)
		ins.inspectFields(exprValue.Results)
	case *ast.Ident:
		if exprValue != nil {
			debug.Debug("Ident, name: %s, obj: %+v\n", exprValue.Name, exprValue.Obj)
		} else {
			debug.Debug("Ident is nil: %+v\n", expr)
		}
	case *ast.IndexExpr: // s[1], arr[1]
		ins.inspectExpr(exprValue.X)
		ins.inspectExpr(exprValue.Index)
	case *ast.InterfaceType: // interface { A(); B() }
		fieldResult := ins.inspectFields(exprValue.Methods)
		result.Fields = fieldResult.Fields
	case *ast.KeyValueExpr: // key:value
		ins.inspectExpr(exprValue.Key)
		ins.inspectExpr(exprValue.Value)
	case *ast.MapType: // map[string]T
		ins.inspectExpr(exprValue.Key)
		ins.inspectExpr(exprValue.Value)
	case *ast.ParenExpr: // (1==1)
		ins.inspectExpr(exprValue.X)
	case *ast.UnaryExpr: // *a
		ins.inspectExpr(exprValue.X)
	}

	return
}

func (ins *Inspector) inspectStmt(stmt ast.Stmt) (result StmtResult) {
	if stmt == nil {
		return
	}
	result = MakeStmtResult()

	switch stmtValue := stmt.(type) {
	case *ast.AssignStmt: // a, b := 1, 2
		for _, lhs := range stmtValue.Lhs {
			ins.inspectExpr(lhs)
		}
		for _, rhs := range stmtValue.Rhs {
			ins.inspectExpr(rhs)
		}
	case *ast.SelectStmt: // select { }
		ins.inspectStmt(stmtValue.Body)
	case *ast.SendStmt: // c <- 1
		ins.inspectExpr(stmtValue.Chan)
		ins.inspectExpr(stmtValue.Value)
	case *ast.SwitchStmt: // switch { }
		ins.inspectStmt(stmtValue.Init)
		ins.inspectExpr(stmtValue.Tag)
		ins.inspectStmt(stmtValue.Body)
	case *ast.BadStmt:
		panic(fmt.Errorf("BadStmt: %+v", stmtValue))
	case *ast.BlockStmt:
		for _, single := range stmtValue.List {
			ins.inspectStmt(single)
		}
	case *ast.BranchStmt:
		ins.inspectExpr(stmtValue.Label)
	case *ast.CaseClause:
		for _, one := range stmtValue.List {
			ins.inspectExpr(one)
		}
		for _, one := range stmtValue.Body {
			ins.inspectStmt(one)
		}
	case *ast.CommClause:
		ins.inspectStmt(stmtValue.Comm)
		for _, one := range stmtValue.Body {
			ins.inspectStmt(one)
		}
	case *ast.DeclStmt:
		ins.inspectDecl(stmtValue.Decl)
	case *ast.DeferStmt:
		ins.inspectExpr(stmtValue.Call)
	case *ast.EmptyStmt:
	case *ast.ExprStmt:
		ins.inspectExpr(stmtValue.X)
	case *ast.ForStmt: // for i:=0; i< l; i++ { }
		ins.inspectStmt(stmtValue.Init)
		ins.inspectExpr(stmtValue.Cond)
		ins.inspectStmt(stmtValue.Post)
		ins.inspectStmt(stmtValue.Body)
	case *ast.GoStmt:
		ins.inspectExpr(stmtValue.Call)
	case *ast.IfStmt:
		ins.inspectStmt(stmtValue.Init)
		ins.inspectExpr(stmtValue.Cond)
		ins.inspectStmt(stmtValue.Body)
		ins.inspectStmt(stmtValue.Else)
	case *ast.IncDecStmt:
		ins.inspectExpr(stmtValue.X)
	case *ast.LabeledStmt:
		ins.inspectExpr(stmtValue.Label)
		ins.inspectStmt(stmtValue.Stmt)
	case *ast.RangeStmt: // for key, value := range slice { }
		ins.inspectExpr(stmtValue.Key)
		ins.inspectExpr(stmtValue.Value)
		ins.inspectExpr(stmtValue.X)
		ins.inspectStmt(stmtValue.Body)
	case *ast.ReturnStmt:
		for _, one := range stmtValue.Results {
			ins.inspectExpr(one)
		}
	case *ast.TypeSwitchStmt: // switch x := m(); a := x.(type) { }
		ins.inspectStmt(stmtValue.Init)
		ins.inspectStmt(stmtValue.Assign)
		ins.inspectStmt(stmtValue.Body)
	}

	return
}

func (ins *Inspector) inspectFields(fields *ast.FieldList) (result FieldResult) {
	if fields == nil {
		return
	}
	result = MakeFieldResult()

	var _ *ast.Field // 是一个Node，但不是一个Expr

	for _, field := range fields.List {
		// 拿field的名称，类型，tag，注释，文档
		debug.Debug("StructType field name: %v, type: %+v, tag: %v, comment: %s, doc: %s\n", field.Names, field.Type, field.Tag, field.Comment.Text(), field.Doc.Text())

		// 获取receiver name
		fieldTyp := field.Type
		if singleTyp, ok := field.Type.(*ast.StarExpr); ok {
			fieldTyp = singleTyp.X
		}
		result.RecvName = toString(fieldTyp)

		ins.inspectExpr(field.Type)

		name := ""
		if len(field.Names) != 0 {
			for _, s := range field.Names {
				name += s.Name
			}
		} else {
			// 匿名结构体
			name = toString(field.Type)
		}

		result.Fields = append(result.Fields, Field{
			Id:        name,
			Name:      name,
			TypesType: ins.pkg.TypesInfo.TypeOf(field.Type),
			Type:      toString(field.Type),
			Doc:       field.Doc.Text(),
			Comment:   field.Comment.Text(),
		})
	}

	return
}

func toString(v interface{}) string {
	qualifier := pkgNameQualifier(qualifierParam{})

	switch vv := v.(type) {
	case ast.Expr:
		return types.ExprString(vv)
	case types.Type:
		return types.TypeString(vv, qualifier)
	case types.Object:
		return types.ObjectString(vv, qualifier)
	case *types.Selection:
		return types.SelectionString(vv, qualifier)
	default:
		return fmt.Sprintf("%v", v)
	}
}

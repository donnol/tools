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
	funcMap := make(map[string]Func)
	for i, astFile := range pkg.Syntax {
		// 替换import path
		if ins.parser.replaceImportPath {
			fileName := pkg.CompiledGoFiles[i]
			debug.Debug("%v\n", pkg.CompiledGoFiles)
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
		for k, v := range fileResult.funcMap {
			funcMap[k] = v
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
	funcs := make([]Func, 0, len(funcMap))
	for _, single := range funcMap {
		funcs = append(funcs, single)
	}

	return Package{
		Package:    pkg,
		Structs:    structs,
		Interfaces: inters,
		Funcs:      funcs,
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
	funcMap := make(map[string]Func)
	for _, decl := range file.Decls {
		declResult := ins.inspectDecl(decl, "")
		for k, v := range declResult.structMap {
			structMap[k] = v
		}
		for k, v := range declResult.methodMap {
			methodsMap[k] = append(methodsMap[k], v...)
		}
		for k, v := range declResult.interfaceMap {
			interfaceMap[k] = v
		}
		for k, v := range declResult.funcMap {
			funcMap[k] = v
		}
	}
	result.structMap = structMap
	result.methodMap = methodsMap
	result.interfaceMap = interfaceMap
	result.funcMap = funcMap

	return
}

func (ins *Inspector) inspectDecl(decl ast.Decl, from string) (result DeclResult) {
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
			PkgPath:   obj.Pkg().Path(),
			Name:      obj.Name(),
			Signature: toString(obj.Type()),
		}
		from = method.Name

		ins.inspectExpr(declValue.Type, from)               // 函数签名
		stmtResult := ins.inspectStmt(declValue.Body, from) // 函数体
		for _, oneFunc := range stmtResult.funcMap {
			method.Calls = append(method.Calls, oneFunc)
		}

		debug.Debug(from+"method: %+v\n", method)

		// method receiver: func (x *X) XXX()里的(x *X)部分
		var recvName string
		if declValue.Recv != nil { // 方法
			debug.Debug("FundDecl recv: %v\n", declValue.Recv.List)

			fieldResult := ins.inspectFields(declValue.Recv, from)
			recvName = fieldResult.RecvName
			method.Recv = recvName

			result.methodMap[recvName] = append(result.methodMap[recvName], method)
		}

		// 函数和方法
		result.funcMap[obj.Name()] = method

	case *ast.GenDecl:
		switch declValue.Tok {
		case token.IMPORT:
		case token.CONST:
		case token.VAR:
		case token.TYPE:
			for _, spec := range declValue.Specs {
				specResult := ins.inspectSpec(spec, from)
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

func (ins *Inspector) inspectSpec(spec ast.Spec, from string) (result SpecResult) {
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
			exprResult := ins.inspectExpr(specValue.Type, from)
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
			mock, imports := inter.MakeMock()
			debug.Debug("mock: %s, imports: %v\n", mock, imports)
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
			exprResult := ins.inspectExpr(specValue.Type, from)
			structOne.Fields = exprResult.Fields
			result.structMap[specValue.Name.Name] = structOne
		}
	}

	return
}

func (ins *Inspector) inspectExpr(expr ast.Expr, from string) (result ExprResult) {
	if expr == nil {
		return
	}
	result = MakeExprResult()

	switch exprValue := expr.(type) {
	case *ast.StructType:
		fieldResult := ins.inspectFields(exprValue.Fields, from)
		result.Fields = fieldResult.Fields

	case *ast.StarExpr: // *T
		exprResult := ins.inspectExpr(exprValue.X, from)
		result = result.Merge(exprResult)

	case *ast.TypeAssertExpr: // X.(*T)
		ins.inspectExpr(exprValue.X, from)
		ins.inspectExpr(exprValue.Type, from)

	case *ast.ArrayType: // [L]T
		ins.inspectExpr(exprValue.Len, from)
		ins.inspectExpr(exprValue.Elt, from)

	case *ast.BadExpr:
		panic(fmt.Errorf("BadExpr: %+v", exprValue))

	case *ast.SelectorExpr: // X.M
		debug.Debug("SelectorExpr value: %v, typesString: %s\n", exprValue, toString(exprValue))

		exprResult := ins.inspectExpr(exprValue.X, from)
		result = result.Merge(exprResult)

		pkgID, ok := exprValue.X.(*ast.Ident)
		if ok {
			if so, ok := ins.pkg.TypesInfo.Uses[pkgID].(*types.PkgName); ok {
				pkgPath := so.Imported().Path()
				debug.Debug(from+"SelectorExpr pkgPath: %#v\n", pkgPath)
				result.pkgPath = pkgPath
			}
		}

		debug.Debug(from+"SelectorExpr value: %#v, result: %#v\n", exprValue, result)

	case *ast.SliceExpr: // []T, slice[1:3:5]
		ins.inspectExpr(exprValue.X, from)
		ins.inspectExpr(exprValue.Low, from)
		ins.inspectExpr(exprValue.High, from)
		ins.inspectExpr(exprValue.Max, from)

	case *ast.BasicLit: // 33 40.0 0x1f

	case *ast.BinaryExpr: // X+Y X-Y X*Y X/Y X%Y
		exprResult := ins.inspectExpr(exprValue.X, from)
		result = result.Merge(exprResult)
		exprResult = ins.inspectExpr(exprValue.Y, from)
		result = result.Merge(exprResult)
		debug.Debug(from+"BinaryExpr: %+v\n", result)

	case *ast.CallExpr: // M(1, 2)
		debug.Debug(from+"funcMap 1: %#v, %+v\n", exprValue.Fun, result)
		exprResult := ins.inspectExpr(exprValue.Fun, from)
		debug.Debug(from+"funcMap mid: %#v, %+v\n", exprValue.Fun, exprResult)

		result.funcMap[toString(exprValue.Fun)] = Func{
			PkgPath: exprResult.pkgPath,
			Name:    toString(exprValue.Fun),
		}

		result = result.Merge(exprResult)
		debug.Debug(from+"funcMap 2: %#v, %+v\n", exprValue.Fun, result)

		for _, arg := range exprValue.Args {
			debug.Debug("CallExpr: %+v, %+v\n", exprValue.Fun, arg)
			exprResult := ins.inspectExpr(arg, from)
			result = result.Merge(exprResult)
		}
		debug.Debug(from+"funcMap: %+v\n", result)

	case *ast.ChanType: // chan T, <-chan T, chan<- T
		exprResult := ins.inspectExpr(exprValue.Value, from)
		result = result.Merge(exprResult)

	case *ast.CompositeLit: // T{Name: Value}
		ins.inspectExpr(exprValue.Type, from)
		for _, elt := range exprValue.Elts {
			exprResult := ins.inspectExpr(elt, from)
			result = result.Merge(exprResult)
		}

	case *ast.Ellipsis: // ...int, [...]Arr
		ins.inspectExpr(exprValue.Elt, from)

	case *ast.FuncLit:
		ins.inspectExpr(exprValue.Type, from)
		ins.inspectStmt(exprValue.Body, from)

	case *ast.FuncType:
		ins.inspectFields(exprValue.Params, from)
		ins.inspectFields(exprValue.Results, from)

	case *ast.Ident:

		if exprValue != nil {
			debug.Debug(from+"Ident, name: %s, obj: %+v\n", exprValue.Name, exprValue.Obj)
		} else {
			debug.Debug(from+"Ident is nil: %+v\n", expr)
		}

		obj, ok := ins.pkg.TypesInfo.Uses[exprValue]
		if ok {
			if obj.Pkg() != nil {
				_ = obj.Pkg().Path() // 变量的包路径

				// 变量类型的包路径
				var varTypePkgPath string
				if ptr, ok := obj.Type().(*types.Pointer); ok {
					// FIXME:改用parseTypesType统一处理types.Type信息
					switch ptrElem := ptr.Elem().(type) {
					case *types.Named:
						varTypePkgPath = ptrElem.Obj().Pkg().Path()
						debug.Debug(from+"Ident obj: %#v, ptr: %#v, pkgPath: %#v\n", obj.Type(), ptr, varTypePkgPath)
					}
				}
				result.pkgPath = varTypePkgPath
			}
		}

		debug.Debug(from+"Ident value: %#v, result: %#v\n", exprValue, result)

	case *ast.IndexExpr: // s[1], arr[1]
		exprResult := ins.inspectExpr(exprValue.X, from)
		result = result.Merge(exprResult)
		exprResult = ins.inspectExpr(exprValue.Index, from)
		result = result.Merge(exprResult)

	case *ast.InterfaceType: // interface { A(); B() }
		fieldResult := ins.inspectFields(exprValue.Methods, from)
		result.Fields = fieldResult.Fields

	case *ast.KeyValueExpr: // key:value
		ins.inspectExpr(exprValue.Key, from)
		exprResult := ins.inspectExpr(exprValue.Value, from)
		result = result.Merge(exprResult)

	case *ast.MapType: // map[string]T
		exprResult := ins.inspectExpr(exprValue.Key, from)
		result = result.Merge(exprResult)
		exprResult = ins.inspectExpr(exprValue.Value, from)
		result = result.Merge(exprResult)

	case *ast.ParenExpr: // (1==1)
		exprResult := ins.inspectExpr(exprValue.X, from)
		result = result.Merge(exprResult)

	case *ast.UnaryExpr: // *a
		exprResult := ins.inspectExpr(exprValue.X, from)
		result = result.Merge(exprResult)

	}

	return
}

func (ins *Inspector) inspectStmt(stmt ast.Stmt, from string) (result StmtResult) {
	if stmt == nil {
		return
	}
	result = MakeStmtResult()

	switch stmtValue := stmt.(type) {
	case *ast.AssignStmt: // a, b := 1, 2
		for _, lhs := range stmtValue.Lhs {
			ins.inspectExpr(lhs, from)
		}
		for _, rhs := range stmtValue.Rhs {
			exprResult := ins.inspectExpr(rhs, from)
			result = result.MergeExprResult(exprResult)
		}

	case *ast.SelectStmt: // select { }
		stmtResult := ins.inspectStmt(stmtValue.Body, from)
		result = result.Merge(stmtResult)

	case *ast.SendStmt: // c <- 1
		ins.inspectExpr(stmtValue.Chan, from)
		exprResult := ins.inspectExpr(stmtValue.Value, from)
		result = result.MergeExprResult(exprResult)

	case *ast.SwitchStmt: // switch { }
		stmtResult := ins.inspectStmt(stmtValue.Init, from)
		result = result.Merge(stmtResult)
		exprResult := ins.inspectExpr(stmtValue.Tag, from)
		result = result.MergeExprResult(exprResult)
		stmtResult = ins.inspectStmt(stmtValue.Body, from)
		result = result.Merge(stmtResult)

	case *ast.BadStmt:
		panic(fmt.Errorf("BadStmt: %+v", stmtValue))

	case *ast.BlockStmt:
		for _, single := range stmtValue.List {
			debug.Debug(from+"block stmt: %+v\n", single)
			res := ins.inspectStmt(single, from)
			result = result.Merge(res)
		}
		debug.Debug(from+"block funcMap: %+v\n", result.funcMap)

	case *ast.BranchStmt:
		exprResult := ins.inspectExpr(stmtValue.Label, from)
		result = result.MergeExprResult(exprResult)

	case *ast.CaseClause:
		for _, one := range stmtValue.List {
			exprResult := ins.inspectExpr(one, from)
			result = result.MergeExprResult(exprResult)
		}
		for _, one := range stmtValue.Body {
			stmtResult := ins.inspectStmt(one, from)
			result = result.Merge(stmtResult)
		}

	case *ast.CommClause:
		stmtResult := ins.inspectStmt(stmtValue.Comm, from)
		result = result.Merge(stmtResult)
		for _, one := range stmtValue.Body {
			stmtResult := ins.inspectStmt(one, from)
			result = result.Merge(stmtResult)
		}

	case *ast.DeclStmt:
		ins.inspectDecl(stmtValue.Decl, from)

	case *ast.DeferStmt:
		exprResult := ins.inspectExpr(stmtValue.Call, from)
		result = result.MergeExprResult(exprResult)

	case *ast.EmptyStmt:

	case *ast.ExprStmt:
		debug.Debug(from+"expr stmt: %+v\n", stmtValue.X)
		exprResult := ins.inspectExpr(stmtValue.X, from)
		result = result.MergeExprResult(exprResult)
		debug.Debug(from+"expr funcMap: %+v\n", result.funcMap)

	case *ast.ForStmt: // for i:=0; i< l; i++ { }
		ins.inspectStmt(stmtValue.Init, from)
		exprResult := ins.inspectExpr(stmtValue.Cond, from)
		result = result.MergeExprResult(exprResult)
		ins.inspectStmt(stmtValue.Post, from)
		stmtResult := ins.inspectStmt(stmtValue.Body, from)
		result = result.Merge(stmtResult)

	case *ast.GoStmt:
		exprResult := ins.inspectExpr(stmtValue.Call, from)
		result = result.MergeExprResult(exprResult)

	case *ast.IfStmt:
		stmtResult := ins.inspectStmt(stmtValue.Init, from)
		result = result.Merge(stmtResult)
		exprResult := ins.inspectExpr(stmtValue.Cond, from)
		result = result.MergeExprResult(exprResult)
		stmtResult = ins.inspectStmt(stmtValue.Body, from)
		result = result.Merge(stmtResult)
		stmtResult = ins.inspectStmt(stmtValue.Else, from)
		result = result.Merge(stmtResult)

	case *ast.IncDecStmt:
		exprResult := ins.inspectExpr(stmtValue.X, from)
		result = result.MergeExprResult(exprResult)

	case *ast.LabeledStmt:
		exprResult := ins.inspectExpr(stmtValue.Label, from)
		result = result.MergeExprResult(exprResult)
		ins.inspectStmt(stmtValue.Stmt, from)

	case *ast.RangeStmt: // for key, value := range slice { }
		ins.inspectExpr(stmtValue.Key, from)
		ins.inspectExpr(stmtValue.Value, from)
		exprResult := ins.inspectExpr(stmtValue.X, from)
		result = result.MergeExprResult(exprResult)
		stmtResult := ins.inspectStmt(stmtValue.Body, from)
		result = result.Merge(stmtResult)

	case *ast.ReturnStmt:
		for _, one := range stmtValue.Results {
			exprResult := ins.inspectExpr(one, from)
			result = result.MergeExprResult(exprResult)
			debug.Debug(from+"return stmt: %#v, %+v\n", one, result.funcMap)
		}

	case *ast.TypeSwitchStmt: // switch x := m(); a := x.(type) { }
		stmtResult := ins.inspectStmt(stmtValue.Init, from)
		result = result.Merge(stmtResult)
		stmtResult = ins.inspectStmt(stmtValue.Assign, from)
		result = result.Merge(stmtResult)
		stmtResult = ins.inspectStmt(stmtValue.Body, from)
		result = result.Merge(stmtResult)
	}

	return
}

func (ins *Inspector) inspectFields(fields *ast.FieldList, from string) (result FieldResult) {
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

		ins.inspectExpr(field.Type, from)

		name := ""
		if len(field.Names) != 0 {
			for _, s := range field.Names {
				name += s.Name
			}
		} else {
			// 匿名结构体
			name = toString(field.Type)
		}

		tag := ""
		if field.Tag != nil {
			tag = field.Tag.Value
		}
		result.Fields = append(result.Fields, Field{
			Id:          name,
			Name:        name,
			TypesType:   ins.pkg.TypesInfo.TypeOf(field.Type),
			Type:        toString(field.Type),
			Tag:         tag,
			TagBasicLit: field.Tag,
			Doc:         field.Doc.Text(),
			Comment:     field.Comment.Text(),
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

func getTypesPkgPath(t types.Type) string {
	debug.Debug("pvar type: %s\n", t)

	pkgPath := ""
	switch v := t.(type) {
	case *types.Named:
		if v.Obj().Pkg() != nil {
			pkgPath = v.Obj().Pkg().Path()
			debug.Debug("path: %s\n", pkgPath)
		}
	}

	return pkgPath
}

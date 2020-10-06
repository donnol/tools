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
	pkg *packages.Package
}

type InspectOption struct {
}

func NewInspector(opt InspectOption) *Inspector {
	return &Inspector{}
}

func (ins *Inspector) InspectPkg(pkg *packages.Package) Package {
	if pkg == nil {
		panic("input pkg is nil")
	}
	ins.pkg = pkg

	// 解析*ast.File信息
	structMap := make(map[string]Struct)
	methodsMap := make(map[string][]Method)
	for _, astFile := range pkg.Syntax {
		fileResult := ins.InspectFile(astFile)

		for k, v := range fileResult.structMap {
			structMap[k] = v
		}
		for k, v := range fileResult.methodMap {
			methodsMap[k] = append(methodsMap[k], v...)
		}
	}

	structs := make([]Struct, 0, len(structMap))
	for structName, single := range structMap {
		methods := methodsMap[structName]
		single.Methods = methods
		structs = append(structs, single)
	}

	return Package{
		Package: pkg,
		Structs: structs,
	}
}

func (ins *Inspector) InspectFile(file *ast.File) (result FileResult) {
	result = MakeFileResult()

	structMap := make(map[string]Struct)
	methodsMap := make(map[string][]Method)
	for _, decl := range file.Decls {
		declResult := ins.inspectDecl(decl)
		for k, v := range declResult.structMap {
			structMap[k] = v
		}
		for k, v := range declResult.methodMap {
			methodsMap[k] = append(methodsMap[k], v...)
		}
	}
	result.structMap = structMap
	result.methodMap = methodsMap

	return
}

func (ins *Inspector) inspectDecl(decl ast.Decl) (result DeclResult) {
	result = MakeDeclResult()

	switch declValue := decl.(type) {
	case *ast.BadDecl:
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
			}
		}
	}

	return
}

func (ins *Inspector) inspectSpec(spec ast.Spec) (result SpecResult) {
	result = MakeSpecResult()

	switch specValue := spec.(type) {
	case *ast.ImportSpec:
	case *ast.ValueSpec:
	case *ast.TypeSpec:
		// 这里拿到类型信息: 名称，注释，文档
		debug.Debug("TypeSpec name: %s, comment: %s, doc: %s\n", specValue.Name, specValue.Comment.Text(), specValue.Doc.Text())

		structOne := Struct{
			PkgPath: ins.pkg.PkgPath,
			PkgName: ins.pkg.Name,
			Field: Field{
				Id:      ins.pkg.TypesInfo.Types[specValue.Type].Type.String(),
				Name:    specValue.Name.Name,
				Type:    toString(specValue.Type),
				Doc:     specValue.Doc.Text(),
				Comment: specValue.Comment.Text(),
			},
		}

		// 再拿field
		exprResult := ins.inspectExpr(specValue.Type)
		structOne.Fields = exprResult.Fields
		result.structMap[specValue.Name.Name] = structOne
	}

	return
}

func (ins *Inspector) inspectExpr(expr ast.Expr) (result ExprResult) {
	result = MakeExprResult()

	switch exprValue := expr.(type) {
	case *ast.StructType:
		fieldResult := ins.inspectFields(exprValue.Fields)
		result.Fields = fieldResult.Fields

	case *ast.StarExpr:
		ins.inspectExpr(exprValue.X)
	case *ast.TypeAssertExpr:
	case *ast.ArrayType:
	case *ast.BadExpr:
	case *ast.SelectorExpr:
		debug.Debug("SelectorExpr value: %v, typesString: %s\n", exprValue, toString(exprValue))
	case *ast.SliceExpr:
	case *ast.BasicLit:
	case *ast.BinaryExpr:
	case *ast.CallExpr:
	case *ast.ChanType:
	case *ast.CompositeLit:
	case *ast.Ellipsis:
	case *ast.FuncLit:
	case *ast.FuncType:
		ins.inspectFields(exprValue.Params)
	case *ast.Ident:
	case *ast.IndexExpr:
	case *ast.InterfaceType:
	case *ast.KeyValueExpr:
	case *ast.MapType:
	case *ast.ParenExpr:
	case *ast.UnaryExpr:
	}

	return
}

func (ins *Inspector) inspectStmt(stmt ast.Stmt) (result StmtResult) {
	result = MakeStmtResult()

	switch stmtValue := stmt.(type) {
	case *ast.AssignStmt:
	case *ast.SelectStmt:
	case *ast.SendStmt:
	case *ast.SwitchStmt:
	case *ast.BadStmt:
	case *ast.BlockStmt:
		for _, single := range stmtValue.List {
			ins.inspectStmt(single)
		}
	case *ast.BranchStmt:
	case *ast.CaseClause:
	case *ast.CommClause:
	case *ast.DeclStmt:
		ins.inspectDecl(stmtValue.Decl)
	case *ast.DeferStmt:
	case *ast.EmptyStmt:
	case *ast.ExprStmt:
	case *ast.ForStmt:
	case *ast.GoStmt:
	case *ast.IfStmt:
	case *ast.IncDecStmt:
	case *ast.LabeledStmt:
	case *ast.RangeStmt:
	case *ast.ReturnStmt:
	case *ast.TypeSwitchStmt:
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
			Id:      name,
			Name:    name,
			Type:    toString(field.Type),
			Doc:     field.Doc.Text(),
			Comment: field.Comment.Text(),
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

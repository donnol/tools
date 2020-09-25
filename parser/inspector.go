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
}

type InspectOption struct {
}

func NewInspector(opt InspectOption) *Inspector {
	return &Inspector{}
}

func (ins *Inspector) InspectPkg(pkg *packages.Package) Package {
	// 用pkg.PkgPath和pkg.Module里的目录信息即可拿到导入路径对应的目录信息
	_ = pkg.PkgPath
	_ = pkg.Module

	// 类型信息
	_ = pkg.Types
	_ = pkg.TypesInfo

	// 解析*ast.File信息
	var structs []Struct
	for _, astFile := range pkg.Syntax {
		ins.InspectFile(astFile)
	}

	return Package{
		Package: pkg,
		Structs: structs,
	}
}

func (ins *Inspector) InspectFile(file *ast.File) {
	for _, decl := range file.Decls {
		ins.inspectDecl(decl)
	}
}

func (ins *Inspector) inspectDecl(decl ast.Decl) {
	switch declValue := decl.(type) {
	case *ast.BadDecl:
	case *ast.FuncDecl:
		debug.Debug("FundDecl name: %s, %s\n", declValue.Name, declValue.Doc.Text())

		ins.inspectExpr(declValue.Type) // 函数签名
		ins.inspectStmt(declValue.Body) // 函数体

		// method
		if declValue.Recv != nil {
			debug.Debug("FundDecl recv: %v\n", declValue.Recv.List)
			for _, single := range declValue.Recv.List {
				ins.inspectExpr(single.Type)
			}
		}

	case *ast.GenDecl:
		switch declValue.Tok {
		case token.IMPORT:
		case token.CONST:
		case token.VAR:
		case token.TYPE:
			for _, spec := range declValue.Specs {
				ins.inspectSpec(spec)
			}
		}
	}
}

func (ins *Inspector) inspectSpec(spec ast.Spec) {
	switch specValue := spec.(type) {
	case *ast.ImportSpec:
	case *ast.ValueSpec:
	case *ast.TypeSpec:
		// 这里拿到类型信息: 名称，注释，文档
		debug.Debug("TypeSpec name: %s, comment: %s, doc: %s\n", specValue.Name, specValue.Comment.Text(), specValue.Doc.Text())

		// 再拿field
		ins.inspectExpr(specValue.Type)
	}
}

func (ins *Inspector) inspectExpr(expr ast.Expr) {
	switch exprValue := expr.(type) {
	case *ast.StructType:
		var _ *ast.Field // 是一个Node，但不是一个Expr

		for _, field := range exprValue.Fields.List {
			// 拿field的名称，类型，tag，注释，文档
			debug.Debug("StructType field name: %v, type: %+v, tag: %v, comment: %s, doc: %s\n", field.Names, field.Type, field.Tag, field.Comment.Text(), field.Doc.Text())
			ins.inspectExpr(field.Type)
		}
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
		for _, field := range exprValue.Params.List {
			debug.Debug("FuncType params name: %v, type: %+v\n", field.Names, field.Type)
			ins.inspectExpr(field.Type)
		}
	case *ast.Ident:
	case *ast.IndexExpr:
	case *ast.InterfaceType:
	case *ast.KeyValueExpr:
	case *ast.MapType:
	case *ast.ParenExpr:
	case *ast.UnaryExpr:
	}
}

func (ins *Inspector) inspectStmt(stmt ast.Stmt) {
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

package parser

import (
	"go/ast"
	"go/token"

	"github.com/donnol/tools/internal/utils/debug"
)

type Inspector struct {
}

// TODO:
func inspectDecl(decl ast.Decl) {
	switch declValue := decl.(type) {
	case *ast.BadDecl:
	case *ast.FuncDecl:
		inspectStmt(declValue.Body)
	case *ast.GenDecl:
		if declValue.Tok == token.TYPE {
			for _, spec := range declValue.Specs {
				debug.Debug("spec of decl: %+v\n", spec)
				inspectSpec(spec)
			}
		}
	}
}

func inspectSpec(spec ast.Spec) {
	switch specValue := spec.(type) {
	case *ast.ImportSpec:
	case *ast.ValueSpec:
	case *ast.TypeSpec:
		inspectExpr(specValue.Type)
	}
}

func inspectExpr(expr ast.Expr) {
	switch exprValue := expr.(type) {
	case *ast.StructType:
		for _, field := range exprValue.Fields.List {
			switch field.Type.(type) {
			case *ast.SelectorExpr:
			case *ast.ArrayType:
			case *ast.StructType:
				// more...
			}
		}
	case *ast.StarExpr:
	case *ast.TypeAssertExpr:
	case *ast.ArrayType:
	case *ast.BadExpr:
	case *ast.SelectorExpr:
	case *ast.SliceExpr:
	case *ast.BasicLit:
	case *ast.BinaryExpr:
	case *ast.CallExpr:
	case *ast.ChanType:
	case *ast.CompositeLit:
	case *ast.Ellipsis:
	case *ast.FuncLit:
	case *ast.FuncType:
	case *ast.Ident:
	case *ast.IndexExpr:
	case *ast.InterfaceType:
	case *ast.KeyValueExpr:
	case *ast.MapType:
	case *ast.ParenExpr:
	case *ast.UnaryExpr:
	}
}

func inspectStmt(stmt ast.Stmt) {
	switch stmtValue := stmt.(type) {
	case *ast.AssignStmt:
	case *ast.SelectStmt:
	case *ast.SendStmt:
	case *ast.SwitchStmt:
	case *ast.BadStmt:
	case *ast.BlockStmt:
		for _, single := range stmtValue.List {
			inspectStmt(single)
		}
	case *ast.BranchStmt:
	case *ast.CaseClause:
	case *ast.CommClause:
	case *ast.DeclStmt:
		inspectDecl(stmtValue.Decl)
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

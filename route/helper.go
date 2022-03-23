package route

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"

	"github.com/pkg/errors"
)

// resolveCallExpr 解析函数调用字符串
func resolveCallExpr(funcCall string) (name string, v []any, t token.Token, err error) {

	expr, err := parser.ParseExpr(funcCall)
	if err != nil {
		err = errors.Errorf("parse vfunc failed, error: %+v\n", err)
		return
	}

	// 函数调用表达式
	callExpr, ok := expr.(*ast.CallExpr)
	if !ok {
		err = errors.Errorf("'%s' is not func call\n", funcCall)
		return
	}
	name = callExpr.Fun.(*ast.Ident).Name

	// 参数值和类型
	for _, arg := range callExpr.Args {
		switch arg := arg.(type) {
		case *ast.BasicLit:
			lit := arg
			t = lit.Kind
			var lv any
			switch t {
			case token.INT:
				lv, _ = strconv.Atoi(lit.Value)
			case token.FLOAT:
				lv, _ = strconv.ParseFloat(lit.Value, 64)
			}
			v = append(v, lv)
		case *ast.Ident:
			ident := arg
			v = append(v, ident.Name)
			t = token.STRING
		}
	}

	return
}

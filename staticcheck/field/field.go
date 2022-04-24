package field

import (
	"fmt"
	"go/ast"
	"go/types"
	"log"
	"strings"

	"golang.org/x/tools/go/packages"
)

// 针对github.com/fishedee/tools/query，对于Column(objs, "id")这种调用，检查obj结构体中是否存在id字段
type Checker interface {
	// 指定pkg，分析其中源码，找出含有指定包调用的地方，再检测参数里涉及的结构体和字段是否匹配
	CheckField(pkg string) error
	// 注册包，后续检查只针对这些包的调用
	RegisterPkg(pkg string)
}

func New() Checker {
	return &checkerImpl{
		pkgs: make([]string, 0, 8),
	}
}

type checkerImpl struct {
	pkgs []string
}

// ast遍历
// find CallExpr
// find param struct
// check if the field is exist in struct
func (impl *checkerImpl) CheckField(pkg string) error {
	// load pkg
	pkgInfo, err := packages.Load(&packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedExportsFile |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedTypesSizes |
			packages.NeedModule,
	}, pkg)
	if err != nil {
		return fmt.Errorf("package load failed: %w", err)
	}
	for _, pi := range pkgInfo {
		for _, file := range pi.Syntax {
			existImport := false
			for _, imp := range file.Imports {
				pathValue := removeDoubleQuote(imp.Path.Value)
				for _, regPkg := range impl.pkgs {
					if regPkg == pathValue {
						existImport = true
						break
					}
				}
			}
			if !existImport {
				continue
			}
			ast.Inspect(file, func(n ast.Node) bool {

				expr, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				//获取caller信息
				exprIdent, ok := expr.Fun.(*ast.Ident)
				if !ok {
					selectorExpr, ok := expr.Fun.(*ast.SelectorExpr)
					if !ok {
						indexListExpr, ok := expr.Fun.(*ast.IndexListExpr)
						if !ok {
							return true
						}
						exprIdent = indexListExpr.X.(*ast.SelectorExpr).Sel
					} else {
						exprIdent = selectorExpr.Sel
					}
				}

				typesInfo := pi.TypesInfo
				funcObj, ok := typesInfo.Uses[exprIdent].(*types.Func)
				if !ok {
					return true
				}

				// check field
				column := ""
				for i := len(expr.Args) - 1; i >= 0; i-- {
					arg := expr.Args[i]
					t1, isExist := typesInfo.Types[arg]
					if !isExist {
						panic(fmt.Errorf("unknown argument type:%v", expr.Args))
					}

					switch tt := t1.Type.(type) {
					case *types.Basic:
						column = removeDoubleQuote(t1.Value.String())
					case *types.Slice:
						elem, ok := tt.Elem().(*types.Named)
						if !ok {
							continue
						}

						struc, ok := elem.Underlying().(*types.Struct)
						if !ok {
							continue
						}

						existField := false
						for j := 0; j < struc.NumFields(); j++ {
							field := struc.Field(j)
							if field.Name() == column {
								existField = true
								break
							}
						}
						if !existField {
							log.Printf("[Param of %v]:field [%v] can not be found in struct [%v]\n", funcObj, column, elem)
							continue
						}
						log.Printf("[Param of %v]: field [%v] can be found in struct [%v]\n", funcObj, column, elem)
					}
				}

				return true
			})
		}
	}

	return nil
}

func removeDoubleQuote(in string) string {
	pathValue := in
	if strings.Index(pathValue, `"`) == 0 && strings.LastIndex(pathValue, `"`) == len(pathValue)-1 {
		pathValue = pathValue[1 : len(pathValue)-1]
	}
	return pathValue
}

func (impl *checkerImpl) RegisterPkg(pkg string) {
	impl.pkgs = append(impl.pkgs, pkg)
}

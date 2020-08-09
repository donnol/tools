package parser

import (
	"go/ast"
	"go/types"
	"reflect"
	"strings"

	"github.com/donnol/tools/internal/utils/debug"
)

type visitor struct {
	pkgPath string
	info    *types.Info
	structs []Struct

	methodMap map[string][]Method
	fieldMap  map[string]Field
}

// Visit 断言节点类型，然后从info里拿到types.Type信息，再使用parseTypesType方法获取具体类型信息
// 顺序是乱的
func (v *visitor) Visit(node ast.Node) (w ast.Visitor) {
	debug.Debug("=== node: %+v, %+v\n", node, reflect.TypeOf(node))

	switch n := node.(type) {
	case *ast.TypeSpec:
		// type and value
		tv := v.info.Types[n.Type]
		debug.Debug("=== type: %+v, %+v, %+v, %+v, %v, %v, %v\n", n, n.Type, tv, tv.Type, types.ExprString(n.Type), strings.TrimSpace(n.Doc.Text()), n.Comment.Text())

		// obj
		obj := v.info.Defs[n.Name]
		debug.Debug("=== obj: %+v\n", obj)

		doc := strings.TrimSpace(n.Doc.Text())
		v.parseTypesType(tv.Type, obj, doc, "")

	case *ast.Field:
		debug.Debug("=== field: %+v, %+v, %v, %v\n", n, n.Type, strings.TrimSpace(n.Doc.Text()), n.Comment.Text())
		tv := v.info.Types[n.Type]
		for _, name := range n.Names {
			obj, ok := v.info.Defs[name]
			if !ok {
				continue
			}
			v.fieldMap[obj.Id()] = Field{
				Id:      obj.Id(),
				Name:    obj.Name(),
				Type:    types.TypeString(obj.Type(), pkgNameQualifier(qualifierParam{pkgPath: v.pkgPath})),
				Doc:     n.Doc.Text(),
				Comment: n.Comment.Text(),
			}
		}
		v.parseTypesType(tv.Type, nil, "", "")

	case *ast.FieldList:
		debug.Debug("=== fieldList: %+v\n", n.List)

	case *ast.StructType:
		debug.Debug("=== struct: %+v\n", n)

	case *ast.DeclStmt:
		debug.Debug("=== decl: %+v, %+v\n", n, n.Decl)

	case *ast.Ident:
		debug.Debug("=== ident: %+v, %+v\n", n, n.Obj)
		if n.Obj != nil {
			debug.Debug("=== ident obj: %+v, %+v, %+v, %+v\n", n, n.Obj, n.Obj.Decl, reflect.TypeOf(n.Obj.Decl))
			switch nn := n.Obj.Decl.(type) {
			case *ast.TypeSpec:
				obj := v.info.ObjectOf(nn.Name)
				debug.Debug("=== ident type spec: %+v, %+v, %+v\n", nn, obj, reflect.TypeOf(obj))
				switch ot := obj.(type) {
				case *types.TypeName:
					debug.Debug("=== ident type spec type name: %+v, %+v\n", ot, ot.IsAlias())
				}
			}
		}

	case *ast.CommentGroup:
		debug.Debug("=== comment group: %+v, %v\n", n, n.Text())

	case *ast.ImportSpec:
		debug.Debug("=== import spec: %+v, %+v\n", n, n.Path)

	case *ast.BasicLit:
		debug.Debug("=== basicLit: %+v, %+v\n", n, n.Value)

	case *ast.AssignStmt:
		debug.Debug("=== assign stmt: %+v, %+v\n", n, n.Rhs)

	case *ast.GenDecl:
		debug.Debug("=== gen decl: %+v, %+v\n", n, n.Specs)
		for _, singleSpec := range n.Specs {
			switch singleSpec.(type) {
			case *ast.ImportSpec:
				debug.Debug("=== import spec: %+v, %+v\n", n, singleSpec)
			case *ast.TypeSpec:
				debug.Debug("=== type spec: %+v, %+v\n", n, singleSpec)
			case *ast.ValueSpec:
				debug.Debug("=== value spec: %+v, %+v\n", n, singleSpec)
			default:
				debug.Debug("=== gen decl: %+v, %+v\n", n, singleSpec)
			}
		}

	default:
		debug.Debug("=== default: %+v\n", n)
	}

	return v
}

type qualifierParam struct {
	pkgPath string
}

var (
	// 包名
	// 包名有可能不等于包路径的最后一部分的（最后一个'/'后面的部分）
	pkgNameQualifier = func(qp qualifierParam) types.Qualifier {
		return func(pkg *types.Package) string {
			name := pkg.Name()

			// 如果是同一个包内的，省略包名
			if pkg.Path() == qp.pkgPath {
				return ""
			}

			return name
		}
	}
)

// 解析类型提取所需信息
func (v *visitor) parseTypesType(t types.Type, obj types.Object, doc, comment string) {
	pkgSelType := types.TypeString(t, pkgNameQualifier(qualifierParam{pkgPath: v.pkgPath}))
	debug.Debug("pkgSelType: %+v\n", pkgSelType)

	switch tv := t.(type) {
	case *types.Signature:
		debug.Debug("=== signature: %+v, %+v, %+v\n", tv, tv.Params(), tv.Results())

	case *types.Pointer:
		debug.Debug("=== pointer: %+v, %+v\n", tv, tv.Elem())
		v.parseTypesType(tv.Elem(), obj, "", "")

	case *types.Named:
		methods := []Method{}
		for i := 0; i < tv.NumMethods(); i++ {
			met := tv.Method(i)
			methods = append(methods, Method{
				Origin:    met,
				Signature: types.TypeString(met.Type(), pkgNameQualifier(qualifierParam{pkgPath: v.pkgPath})),
			})
		}
		debug.Debug("=== named: %+v, is alias: %v, methods: %+v\n", tv, tv.Obj().IsAlias(), methods)
		if tv.Obj().IsAlias() {
			debug.Debug("===============================: %+v\n", tv)
		}

		// 怎么把methods关联到struct呢？
		v.methodMap[tv.Obj().Id()] = methods

	case *types.Struct:
		fields := []Field{}
		for i := 0; i < tv.NumFields(); i++ {
			field := tv.Field(i)

			if _, ok := v.fieldMap[field.Id()]; ok {
				continue
			}

			tmpField := Field{
				Id:   field.Id(),
				Name: field.Name(),
				Type: types.TypeString(field.Type(), pkgNameQualifier(qualifierParam{pkgPath: v.pkgPath})),
			}
			fields = append(fields, tmpField)
			v.fieldMap[field.Id()] = tmpField
		}
		debug.Debug("=== struct: %+v, fields: %+v\n", tv, fields)

		if obj == nil {
			return
		}
		v.structs = append(v.structs, Struct{
			PkgPath: obj.Pkg().Path(),
			PkgName: obj.Pkg().Name(),
			Field: Field{
				Id:      obj.Id(),
				Name:    obj.Name(),
				Type:    types.TypeString(obj.Type(), pkgNameQualifier(qualifierParam{pkgPath: v.pkgPath})),
				Comment: doc,
			},
			Fields: fields,
		})
	}
}

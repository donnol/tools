package parser

import (
	"go/ast"
	"go/types"
	"reflect"
	"strings"

	"github.com/donnol/tools/internal/utils/debug"
)

// visitor 将要弃用
type visitor struct {
	// 如: github.com/pkg/errors
	pkgPath string // 包路径

	// 如：errors.Frame
	info *types.Info // 类型信息

	// 如：[]Struct{...}
	structs []Struct // 结构体信息

	methodMap map[string][]Method // 方法集
	fieldMap  map[string]Field    // 字段集
}

func newVisitor() *visitor {
	return &visitor{
		structs:   make([]Struct, 0, 16),
		methodMap: make(map[string][]Method),
		fieldMap:  make(map[string]Field),
	}
}

// Visit 断言节点类型，然后从info里拿到types.Type信息，再使用parseTypesType方法获取具体类型信息
// 顺序是乱的
func (v *visitor) Visit(node ast.Node) (w ast.Visitor) {
	debug.Printf("=== node: %+v, %+v\n", node, reflect.TypeOf(node))

	switch n := node.(type) {
	case *ast.TypeSpec:
		// type and value
		tv := v.info.Types[n.Type]
		debug.Printf("=== type: %+v, %+v, %+v, %+v, %v, %v, %v\n", n, n.Type, tv, tv.Type, types.ExprString(n.Type), strings.TrimSpace(n.Doc.Text()), n.Comment.Text())

		// obj
		obj := v.info.Defs[n.Name]
		debug.Printf("=== obj: %+v\n", obj)

		doc := strings.TrimSpace(n.Doc.Text())
		comment := n.Comment.Text()
		v.parseTypesType(tv.Type, obj, doc, comment)

	case *ast.Field:
		debug.Printf("=== field: %+v, %+v, %v, %v\n", n, n.Type, strings.TrimSpace(n.Doc.Text()), n.Comment.Text())
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
		v.parseTypesType(tv.Type, nil, n.Doc.Text(), n.Comment.Text())

	case *ast.FieldList:
		debug.Printf("=== fieldList: %+v\n", n.List)

	case *ast.StructType:
		debug.Printf("=== struct: %+v\n", n)

	case *ast.DeclStmt:
		debug.Printf("=== decl: %+v, %+v\n", n, n.Decl)

	case *ast.Ident:
		debug.Printf("=== ident: %+v, %+v\n", n, n.Obj)
		if n.Obj != nil {
			debug.Printf("=== ident obj: %+v, %+v, %+v, %+v\n", n, n.Obj, n.Obj.Decl, reflect.TypeOf(n.Obj.Decl))
			switch nn := n.Obj.Decl.(type) {
			case *ast.TypeSpec:
				obj := v.info.ObjectOf(nn.Name)
				debug.Printf("=== ident type spec: %+v, %+v, %+v\n", nn, obj, reflect.TypeOf(obj))
				switch ot := obj.(type) {
				case *types.TypeName:
					debug.Printf("=== ident type spec type name: %+v, %+v\n", ot, ot.IsAlias())
				}
			}
		}

	case *ast.CommentGroup:
		debug.Printf("=== comment group: %+v, %v\n", n, n.Text())

	case *ast.ImportSpec:
		debug.Printf("=== import spec: %+v, %+v\n", n, n.Path)

	case *ast.BasicLit:
		debug.Printf("=== basicLit: %+v, %+v\n", n, n.Value)

	case *ast.AssignStmt:
		debug.Printf("=== assign stmt: %+v, %+v\n", n, n.Rhs)

	case *ast.GenDecl:
		debug.Printf("=== gen decl: %+v, %+v\n", n, n.Specs)
		for _, singleSpec := range n.Specs {
			switch singleSpec.(type) {
			case *ast.ImportSpec:
				debug.Printf("=== import spec: %+v, %+v\n", n, singleSpec)
			case *ast.TypeSpec:
				debug.Printf("=== type spec: %+v, %+v\n", n, singleSpec)
			case *ast.ValueSpec:
				debug.Printf("=== value spec: %+v, %+v\n", n, singleSpec)
			default:
				debug.Printf("=== gen decl: %+v, %+v\n", n, singleSpec)
			}
		}

	default:
		debug.Printf("=== default: %+v\n", n)
	}

	return v
}

// 解析类型提取所需信息
func (v *visitor) parseTypesType(t types.Type, obj types.Object, doc, comment string) {
	pkgSelType := types.TypeString(t, pkgNameQualifier(qualifierParam{pkgPath: v.pkgPath}))
	debug.Printf("pkgSelType: %+v\n", pkgSelType)

	switch tv := t.(type) {
	case *types.Signature:
		debug.Printf("=== signature: %+v, %+v, %+v\n", tv, tv.Params(), tv.Results())

	case *types.Pointer:
		debug.Printf("=== pointer: %+v, %+v\n", tv, tv.Elem())
		v.parseTypesType(tv.Elem(), obj, doc, comment)

	case *types.Named:
		methods := []Method{}
		for i := 0; i < tv.NumMethods(); i++ {
			met := tv.Method(i)
			methods = append(methods, Method{
				Origin:    met,
				Signature: types.TypeString(met.Type(), pkgNameQualifier(qualifierParam{pkgPath: v.pkgPath})),
			})
		}
		debug.Printf("=== named: %+v, is alias: %v, methods: %+v\n", tv, tv.Obj().IsAlias(), methods)
		if tv.Obj().IsAlias() {
			debug.Printf("===============================: %+v\n", tv)
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
		debug.Printf("=== struct: %+v, fields: %+v\n", tv, fields)

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
				Doc:     doc,
				Comment: comment,
			},
			Fields: fields,
		})
	}
}

package parser

import (
	"fmt"
	"go/types"
	"unicode"

	"github.com/donnol/tools/internal/utils/debug"
)

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

type parseTypesTypeOption struct {
	_       string
	pkgPath string
}

func parseTypesType(t types.Type, opt parseTypesTypeOption) (r struct {
	methods []Method
}) {
	switch tv := t.(type) {
	case *types.Interface:
		methods := make([]Method, 0, tv.NumMethods())
		for i := 0; i < tv.NumMethods(); i++ {
			met := tv.Method(i)
			methods = append(methods, Method{
				Origin:    met,
				Name:      met.Name(),
				Signature: types.TypeString(met.Type(), pkgNameQualifier(qualifierParam{pkgPath: opt.pkgPath})),
			})
		}
		debug.Printf("| parseTypesType | interface methods: %+v\n", methods)
		r.methods = methods

	case *types.Signature:
		debug.Printf("=== signature: %+v, %+v, %+v\n", tv, tv.Params(), tv.Results())

	case *types.Pointer:
		debug.Printf("=== pointer: %+v, %+v\n", tv, tv.Elem())
		parseTypesType(tv.Elem(), opt)

	case *types.Named:
		methods := []Method{}
		for i := 0; i < tv.NumMethods(); i++ {
			met := tv.Method(i)
			methods = append(methods, Method{
				Origin:    met,
				Signature: types.TypeString(met.Type(), pkgNameQualifier(qualifierParam{pkgPath: opt.pkgPath})),
			})
		}
		debug.Printf("=== named: %+v, is alias: %v, pkgPath: %v, methods: %+v\n", tv, tv.Obj().IsAlias(), tv.Obj().Pkg().Path(), methods)
		if tv.Obj().IsAlias() {
			debug.Printf("===============================: %+v\n", tv)
		}

	case *types.Struct:
		fields := []Field{}
		for i := 0; i < tv.NumFields(); i++ {
			field := tv.Field(i)

			tmpField := Field{
				Id:   field.Id(),
				Name: field.Name(),
				Type: types.TypeString(field.Type(), pkgNameQualifier(qualifierParam{pkgPath: opt.pkgPath})),
			}
			fields = append(fields, tmpField)
		}
		debug.Printf("=== struct: %+v, fields: %+v\n", tv, fields)

	case *types.Slice:
		debug.Printf("| parseTypesType | elem: %+v\n", tv.Elem())
		parseTypesType(tv.Elem(), opt)

	case *types.Array:
		debug.Printf("| parseTypesType | elem: %+v\n", tv.Elem())
		parseTypesType(tv.Elem(), opt)

	case *types.Basic:
		debug.Printf("| parseTypesType | elem: %+v\n", tv.Info())

	case *types.Chan:
		debug.Printf("| parseTypesType | elem: %+v\n", tv.Elem())
		parseTypesType(tv.Elem(), opt)

	case *types.Map:
		debug.Printf("| parseTypesType | key: %+v, value: %+v\n", tv.Key(), tv.Elem())
		parseTypesType(tv.Key(), opt)
		parseTypesType(tv.Elem(), opt)

	case *types.Tuple:
		debug.Printf("| parseTypesType | len: %+v\n", tv.Len())

	default:
		fmt.Printf("| parseTypesType | tv: %+v\n", tv)
	}

	return
}

func UcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func LcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

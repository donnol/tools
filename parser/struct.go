package parser

import (
	"fmt"
	"go/types"

	"github.com/donnol/tools/format"
	"github.com/donnol/tools/importpath"
)

// Struct 结构体
type Struct struct {
	// 如：github.com/pkg/errors
	PkgPath string // 包路径

	// 如: errors
	PkgName string // 包名

	Field

	Fields  []Field  // 字段列表
	Methods []Method // 方法列表

	isAlias bool
}

// --- 测试方法

// 让它传入本包里的另外一个结构体
// 传入本项目其它包的结构体
func (s Struct) String(f Field, ip importpath.ImportPath) {
	fmt.Printf("%s\n", s.PkgName)
}

func (s Struct) TypeAlias(p IIIIIIIInfo) {

}

func (s Struct) Demo(in types.Array) types.Basic {
	return types.Basic{}
}

// --- 测试方法

// MMakeInterface 根据结构体的方法生成相应接口
func (s Struct) MakeInterface() string {
	methods := make([]*types.Func, 0, len(s.Methods))
	for _, m := range s.Methods {
		if !m.Origin.Exported() {
			continue
		}
		methods = append(methods, m.Origin)
		fmt.Printf("method: %+v, %s\n", m.Origin, m.Signature)
	}

	if len(methods) == 0 {
		return ""
	}

	i := types.NewInterfaceType(methods, nil)
	i = i.Complete()
	// FIXME:这里拿到的isAlias不是结构体的，而是参数的，所以不准确，没法用
	is := types.TypeString(i.Underlying(), pkgNameQualifier(qualifierParam{pkgPath: s.PkgPath, isAlias: s.isAlias}))
	fmt.Printf("=== alias: %v, %s\n", s.isAlias, is)

	is = interfacePrefix(s.makeInterfaceName(), is)

	// 检查获得的接口定义是否规范
	formatContent, err := format.Format("", is, true)
	if err != nil {
		panic(err)
	}

	return string(formatContent)
}

func (s Struct) makeInterfaceName() string {
	return "I" + s.Name
}

func interfacePrefix(name, is string) string {
	return "type " + name + " " + is
}

type Method struct {
	Origin    *types.Func
	Signature string
}

// Field 字段
type Field struct {
	Id      string // 唯一标志
	Name    string // 名称
	Type    string // 类型，包含包导入路径
	Doc     string // 文档
	Comment string // 注释
}

// IIIIIIIInfo 别名测试
type IIIIIIIInfo = Field // 别名测试注释

type PkgInfo struct {
	dir     string
	pkgName string
}

func (i PkgInfo) GetDir() string {
	return i.dir
}

func (i PkgInfo) GetPkgName() string {
	return i.pkgName
}

package parser

import (
	"io"
	"os"
)

// Option 选项
type Option struct {
	Op Op // 操作，如生成接口，生成实现等

	Filter            func(os.FileInfo) bool // 过滤器
	UseSourceImporter bool                   // 使用源码importer

	ReplaceImportPath bool // 替换导入路径
	FromPath          string
	ToPath            string
	Output            io.Writer

	NeedCall bool // 需要记录调用了哪些函数/方法

	ReplaceCallExpr bool // 替换调用表达式
}

type Op string

const (
	OpReplace          Op = "replace"
	OpMock             Op = "mock"
	OpImpl             Op = "impl"
	OpInterface        Op = "interface"
	OpCallgraph        Op = "callgraph"
	OpGenProject       Op = "genproject"
	OpGenProxy         Op = "genproxy"
	OpFind             Op = "find"
	OpGenStructFromSQL Op = "sql2struct"
	OpGenDataForTable  Op = "gendata"
)

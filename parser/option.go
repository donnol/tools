package parser

import (
	"io"
	"os"
)

// Option 选项
type Option struct {
	Filter            func(os.FileInfo) bool // 过滤器
	UseSourceImporter bool                   // 使用源码importer

	ReplaceImportPath bool // 替换导入路径
	FromPath          string
	ToPath            string
	Output            io.Writer
}

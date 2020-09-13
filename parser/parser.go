package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "bytes"

	"github.com/donnol/tools/format"
	"github.com/donnol/tools/internal/utils/debug"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

// 为什么要这样写，因为：https://github.com/golang/go/issues/27477
type (
	// Parser 解析器
	// 解析指定的包导入路径，获取go源码信息
	Parser struct {
		filter func(os.FileInfo) bool // 过滤器

		fset *token.FileSet

		useSourceImporter bool // 使用源码importer

		replaceImportPath bool // 替换导入路径
		fromPath          string
		toPath            string
		output            io.Writer

		PkgInfo
	}
)

// New 新建
func New(opt Option) *Parser {
	return &Parser{
		filter:            opt.Filter,
		useSourceImporter: opt.UseSourceImporter,
		replaceImportPath: opt.ReplaceImportPath,
		fromPath:          opt.FromPath,
		toPath:            opt.ToPath,
		output:            opt.Output,
	}
}

// ParseAST 解析导入路径的ast，返回目录里的结构体信息
func (p *Parser) ParseAST(importPath string) (structs []Struct, err error) {

	fullDir, err := p.getFullDir(importPath)
	if err != nil {
		return
	}
	p.PkgInfo.dir = fullDir

	fset := token.NewFileSet()
	pkgs, err := p.parseDir(fset, fullDir)
	if err != nil {
		return
	}
	p.fset = fset

	// 提取文件信息
	for pkgName, pkg := range pkgs {
		p.PkgInfo.pkgName = pkgName
		fmt.Printf("pkgName: %s\n", pkgName)

		files := make([]*ast.File, 0, len(pkg.Files))
		for fileName, file := range pkg.Files {
			files = append(files, file)

			if p.replaceImportPath {
				// 替换import path
				if err = p.replaceFileImportPath(fileName, file); err != nil {
					return structs, err
				}
			}
		}
		// 如果执行了路径替换，马上退出
		if p.replaceImportPath {
			continue
		}

		// 获取类型信息
		info, err := p.typesCheck(pkgName, files)
		if err != nil {
			return structs, err
		}

		// 遍历ast
		vis := &visitor{
			pkgPath:   importPath,
			info:      info,
			structs:   make([]Struct, 0, 16),
			methodMap: make(map[string][]Method),
			fieldMap:  make(map[string]Field),
		}
		ast.Walk(vis, pkg)

		tmpStructs := make([]Struct, 0, len(vis.structs))
		for _, single := range vis.structs {
			single.Methods = vis.methodMap[single.Id]
			fields := make([]Field, 0, len(single.Fields))
			for _, field := range single.Fields {
				tmpField := vis.fieldMap[field.Id]
				field.Doc = tmpField.Doc
				field.Comment = tmpField.Comment
				fields = append(fields, field)
			}
			single.Fields = fields
			tmpStructs = append(tmpStructs, single)
		}
		structs = append(structs, tmpStructs...)
	}

	return
}

func (p *Parser) GetPkgInfo() PkgInfo {
	return p.PkgInfo
}

func (p *Parser) getFullDir(importPath string) (fullDir string, err error) {
	// 根据导入路径，获取完整路径等信息
	buildPkg, err := build.Import(importPath, "", build.ImportComment)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	fullDir = buildPkg.Dir

	return
}

// ParseByGoPackages 使用x/tools/go/packages解析指定导入路径
func (p *Parser) ParseByGoPackages(importPath string) (err error) {
	cfg := &packages.Config{
		Mode: packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedSyntax |
			packages.NeedModule |
			packages.NeedImports |
			packages.NeedTypes |
			packages.NeedTypesInfo |
			packages.NeedName,
	}
	pkgs, err := packages.Load(cfg, importPath)
	if err != nil {
		return
	}
	// TODO:
	for _, pkg := range pkgs {
		fmt.Printf("pkg: %#v, %+v\n", pkg, pkg.Module)
		// 用pkg.PkgPath和pkg.Module里的目录信息即可拿到导入路径对应的目录信息
		// 这里已经拿到*ast.File信息了
	}

	return
}

func (p *Parser) parseDir(fset *token.FileSet, fullDir string) (pkgs map[string]*ast.Package, err error) {
	const (
		testSuffix = "_test"
	)

	// 解析目录
	pkgs, err = parser.ParseDir(fset, fullDir, func(fi os.FileInfo) bool {
		li := strings.LastIndex(fi.Name(), filepath.Ext(fi.Name()))

		// 跳过test文件
		testi := strings.LastIndex(fi.Name(), testSuffix)
		if testi != -1 && li-testi == len(testSuffix) {
			return false
		}

		return true
	}, parser.ParseComments)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return
}

func (p *Parser) typesCheck(path string, files []*ast.File) (info *types.Info, err error) {
	imp := importer.Default()
	if p.useSourceImporter {
		fset := token.NewFileSet()
		imp = importer.ForCompiler(fset, "source", nil)
	}

	// 获取类型信息
	conf := types.Config{
		IgnoreFuncBodies: true,

		// 默认是用go install安装后生成的.a文件，可以选择使用source，但是会慢很多
		Importer: imp,

		Error: func(err error) {
			log.Printf("Check Failed: %+v\n", errors.WithStack(err))
		},
		DisableUnusedImportCheck: true,
	}
	info = &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}
	// conf.Check的path参数传入的是包名，而不是导入路径
	pkg, err := conf.Check(path, p.fset, files, info)
	if err != nil {
		err = errors.WithMessagef(err, "info: %+v\n", info)
		return
	}
	p.methodSet(pkg)

	return
}

// 根据types.Type找到method set，但是怎么将它转为interface呢？
func (p *Parser) methodSet(pkg *types.Package) {
	if pkg.Scope() == nil {
		return
	}
	for _, name := range pkg.Scope().Names() {
		obj := pkg.Scope().Lookup(name)
		if obj == nil {
			continue
		}
		typ := obj.Type()
		for _, t := range []types.Type{typ, types.NewPointer(typ)} {
			// fmt.Printf("Method set of %s:\n", t)
			mset := types.NewMethodSet(t)
			for i := 0; i < mset.Len(); i++ {
				sel := mset.At(i)
				_ = sel
				// fmt.Println("sel: ", sel, "type:", sel.Type(), reflect.TypeOf(sel.Type().Underlying()), "obj:", sel.Obj())
			}
			// fmt.Println()
		}
	}
}

func (p *Parser) replaceFileImportPath(fileName string, file *ast.File) error {
	// 替换import path
	for _, fi := range file.Imports {
		path := strings.Trim(fi.Path.Value, `"`)

		if strings.HasPrefix(path, p.fromPath) {
			topath := strings.Replace(path, p.fromPath, p.toPath, 1)

			rewrote := astutil.RewriteImport(p.fset, file, path, topath)

			debug.Debug("From %s to %s, rewrote: %v\n", p.fromPath, p.toPath, rewrote)

		}
	}

	// 获取file的ast内容并格式化
	buf := bytes.NewBuffer([]byte{})
	printer.Fprint(buf, p.fset, file)
	content, err := format.Format(fileName, buf.String(), true)
	if err != nil {
		return err
	}

	// 输出ast节点信息到文件
	var output = p.output
	if p.output == nil {
		// 将内容输出到原文件
		output, err = os.OpenFile(fileName, os.O_RDWR|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}
	}
	_, err = output.Write([]byte(content))
	if err != nil {
		return err
	}

	return nil
}

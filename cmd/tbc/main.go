package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/davecgh/go-spew/spew"
	"github.com/donnol/tools/format"
	"github.com/donnol/tools/importpath"
	"github.com/donnol/tools/internal/utils/debug"
	"github.com/donnol/tools/parser"
	"github.com/donnol/tools/sqlparser"
	"github.com/twpayne/go-jsonstruct"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "tbc",
		Short: "a tool named to be continued",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.UsageString())
		},
	}
)

func main() {
	log.SetFlags(log.Llongfile | log.LstdFlags)

	// 解析标签
	var rFlag bool
	rootCmd.PersistentFlags().BoolVarP(&rFlag, "recursive", "r", false, "recursively process dir from current")
	var path string
	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", "", "specify import path")
	var from string
	rootCmd.PersistentFlags().StringVarP(&from, "from", "", "", "specify from path with replace")
	var to string
	rootCmd.PersistentFlags().StringVarP(&to, "to", "", "", "specify to path with replace")
	var output string
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "specify output file")
	var interf string
	rootCmd.PersistentFlags().StringVarP(&interf, "interface", "", "", "specify interface")
	var funcName string
	rootCmd.PersistentFlags().StringVarP(&funcName, "func", "", "", "specify func or method name")
	var ignore string
	rootCmd.PersistentFlags().StringVarP(&ignore, "ignore", "", "", "specify ignore package or field")
	var depth int
	rootCmd.PersistentFlags().IntVarP(&depth, "depth", "", 0, "specify depth")
	var amount int64
	rootCmd.PersistentFlags().Int64VarP(&amount, "amount", "", 1, "specify amount")
	var file string
	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "specify input file")
	var pkgName string
	rootCmd.PersistentFlags().StringVar(&pkgName, "pkg", "", "specify package name")

	// 添加子命令
	addSubCommand(rootCmd)

	// 执行
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func addSubCommand(rootCmd *cobra.Command) {
	rootCmd.AddCommand(&cobra.Command{
		Use:   string(parser.OpGenStructFromJSON),
		Short: "gen struct from json",
		Long:  `tbc json2struct '{"name": "jack"}'`,
		Run: func(cmd *cobra.Command, args []string) {
			// 标志
			flags := cmd.Flags()
			file, _ := flags.GetString("file")
			output, _ := flags.GetString("output")

			jsonstr := ""
			if len(args) > 0 {
				jsonstr = args[0]
			} else if file != "" {
				data, err := os.ReadFile(file)
				if err != nil {
					fmt.Printf("read file failed: %v\n", err)
					os.Exit(1)
				}
				jsonstr = string(data)
			}

			if jsonstr == "" {
				fmt.Printf("please specify json like '{\"name\": \"jack\"}' or input file by --file=xxx.json\n")
				os.Exit(1)
			}

			obr, err := jsonstruct.ObserveJSON(bytes.NewReader([]byte(jsonstr)))
			if err != nil {
				fmt.Printf("parse json failed: %+v\n", err)
				os.Exit(1)
			}
			gen := jsonstruct.NewGenerator()
			data, err := gen.GoCode(obr)
			if err != nil {
				fmt.Printf("gen go code failed: %+v\n", err)
				os.Exit(1)
			}

			w := os.Stdout
			if output != "" {
				f, err := os.OpenFile(output, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
				if err != nil {
					fmt.Printf("open file %s failed: %v\n", output, err)
					os.Exit(1)
				}
				defer f.Close()

				w = f
			}

			content, err := format.Format(output, string(data), false)
			if err != nil {
				fmt.Printf("format failed: %v\n", err)
				os.Exit(1)
			}
			_, err = w.WriteString(content)
			if err != nil {
				fmt.Printf("write to w failed: %v\n", err)
				os.Exit(1)
			}
			fmt.Println()
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   string(parser.OpGenDataForTable),
		Short: "gen insert statement with data for table",
		Long:  `tbc gendata 'create table user(id int not null)'`,
		Run: func(cmd *cobra.Command, args []string) {
			// 标志
			flags := cmd.Flags()
			ignoreField, _ := flags.GetString("ignore")
			amount, _ := flags.GetInt64("amount")
			file, _ := flags.GetString("file")
			output, _ := flags.GetString("output")

			sql := ""
			if len(args) > 0 {
				sql = args[0]
			} else if file != "" {
				data, err := os.ReadFile(file)
				if err != nil {
					fmt.Printf("read file failed: %v\n", err)
					os.Exit(1)
				}
				sql = string(data)
			}

			if sql == "" {
				fmt.Printf("please specify sql like 'create table user(id int not null)' or input file by --file=xxx.sql\n")
				os.Exit(1)
			}

			opt := sqlparser.Option{}
			if ignoreField != "" {
				opt.IgnoreField = append(opt.IgnoreField, ignoreField)
			}

			ss := sqlparser.ParseCreateSQLBatch(sql)
			if ss == nil {
				fmt.Printf("parse sql failed\n")
				os.Exit(1)
			}
			w := os.Stdout
			if output != "" {
				f, err := os.OpenFile(output, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
				if err != nil {
					fmt.Printf("open file %s failed: %v\n", output, err)
					os.Exit(1)
				}
				defer f.Close()

				w = f
			}
			fmt.Printf("----- Begin generate %d -----\n", amount)
			for _, s := range ss {
				if err := s.GenData(w, amount, opt); err != nil {
					fmt.Printf("gen struct failed: %v\n", err)
					os.Exit(1)
				}
			}
			if output == "" {
				fmt.Println("----- Generate finish -----")
			} else {
				fmt.Printf("----- Generate finish: %s -----\n", output)
			}
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   string(parser.OpGenStructFromSQL),
		Short: "gen struct from sql",
		Long:  `tbc sql2struct 'create table user(id int not null)'`,
		Run: func(cmd *cobra.Command, args []string) {
			// 标志
			flags := cmd.Flags()
			ignoreField, _ := flags.GetString("ignore")
			file, _ := flags.GetString("file")
			output, _ := flags.GetString("output")
			pkg, _ := flags.GetString("pkg")

			sql := ""
			if len(args) > 0 {
				sql = args[0]
			} else if file != "" {
				data, err := os.ReadFile(file)
				if err != nil {
					fmt.Printf("read file failed: %v\n", err)
					os.Exit(1)
				}
				sql = string(data)
			}

			if sql == "" {
				fmt.Printf("please specify sql like 'create table user(id int not null)' or input file by --file=xxx.sql\n")
				os.Exit(1)
			}

			opt := sqlparser.Option{}
			if ignoreField != "" {
				opt.IgnoreField = append(opt.IgnoreField, ignoreField)
			}

			ss := sqlparser.ParseCreateSQLBatch(sql)
			if ss == nil {
				fmt.Printf("parse sql failed\n")
				os.Exit(1)
			}

			w := os.Stdout
			if output != "" {
				f, err := os.OpenFile(output, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
				if err != nil {
					fmt.Printf("open file %s failed: %v\n", output, err)
					os.Exit(1)
				}
				defer f.Close()

				w = f
			}
			buf := new(bytes.Buffer)
			if pkg != "" {
				_, err := buf.WriteString("package " + pkg)
				if err != nil {
					fmt.Printf("write package name failed: %v\n", err)
					os.Exit(1)
				}
			}
			for _, s := range ss {
				if err := s.Gen(buf, opt); err != nil {
					fmt.Printf("gen struct failed: %v\n", err)
					os.Exit(1)
				}
			}
			content, err := format.Format(output, buf.String(), false)
			if err != nil {
				fmt.Printf("format failed: %v\ncontent: %s\n", err, buf.String())
				os.Exit(1)
			}
			_, err = w.WriteString(content)
			if err != nil {
				fmt.Printf("write to w failed: %v\n", err)
				os.Exit(1)
			}
			fmt.Println()
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   string(parser.OpInterface),
		Short: "gen struct interface",
		Long: `gen struct interface, like: 
			type M struct {
				// ...
			}
			func (m *M) String() string {
				return "m.name"
			}
			got: 
			type IM interface {
				String() string
			}
		`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			// 标志
			flags := cmd.Flags()
			path, _ := flags.GetString("path")
			rec, _ := flags.GetBool("recursive")
			fmt.Printf("| interface | %+v, %+v\n", path, rec)

			ip := &importpath.ImportPath{}
			paths, err := getPaths(ip, path, rec)
			if err != nil {
				log.Fatal(err)
			}
			if len(paths) == 0 {
				log.Fatalf("找不到有效路径，请使用-p指定或设置-r！")
			}
			fmt.Printf("dirs: %+v\n", paths)

			// 解析
			p := parser.New(parser.Option{
				Op: parser.OpInterface,
			})
			pkgs, err := p.ParseByGoPackages(paths...)
			if err != nil {
				log.Fatal(err)
			}
			for _, pkg := range pkgs.Pkgs {
				if err = pkg.SaveInterface(""); err != nil {
					log.Fatal(err)
				}
			}
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   string(parser.OpReplace),
		Short: "replace import path",
		Long:  "replace import path, like: from 'import \"zzz.com/xxx/yyy\"' to 'import \"lll.com/mmm/nnn\"'",
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			// 获取from path和to path
			flags := cmd.Flags()
			from, _ := flags.GetString("from")
			to, _ := flags.GetString("to")
			path, _ := flags.GetString("path")

			ip := &importpath.ImportPath{}
			if path == "" {
				path, err = ip.GetByCurrentDir()
				if err != nil {
					log.Fatal(err)
				}
			}
			fmt.Printf("| replace | %+v, %+v, %v\n", from, to, path)

			// 初始化
			p := parser.New(parser.Option{
				Op:                parser.OpReplace,
				ReplaceImportPath: true,
				FromPath:          from,
				ToPath:            to,
			})

			// 执行
			_, err = p.ParseByGoPackages(path)
			if err != nil {
				log.Fatal(err)
			}
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   string(parser.OpMock),
		Short: "gen interface mock struct",
		Long: `gen interface mock struct, like: type I interface { String() string }, 
			gen mock: 
				type Mock struct { StringFunc func() string } 
				var _ I = &Mock{}
				func (mock *Mock) String() string {
					return mock.StringFunc()
				}
			after that, you can use like below:
				var mock = &Mock{
					// init the func like the normal field
					StringFunc: func() string {
						return "jd"								
					},	
				}
				fmt.Println(mock.String())
			`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			// 标志
			flags := cmd.Flags()
			path, _ := flags.GetString("path")
			rec, _ := flags.GetBool("recursive")
			fmt.Printf("| mock | %+v, %+v\n", path, rec)

			ip := &importpath.ImportPath{}
			paths, err := getPaths(ip, path, rec)
			if err != nil {
				log.Fatal(err)
			}
			if len(paths) == 0 {
				log.Fatalf("找不到有效路径，请使用-p指定或设置-r！")
			}
			fmt.Printf("dirs: %+v\n", paths)

			// 解析
			p := parser.New(parser.Option{
				Op:                parser.OpMock,
				UseSourceImporter: true,
			})
			pkgs, err := p.ParseByGoPackages(paths...)
			if err != nil {
				log.Fatal(err)
			}
			for _, pkg := range pkgs.Pkgs {
				if err = pkg.SaveMock(""); err != nil {
					log.Printf("gen mock failed, pkg: %+v, err: %+v\n", pkg, err)
				}
			}
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   string(parser.OpImpl),
		Short: "find implement by given interface in specify path",
		Long: `find implement by given interface in specify path, like: 
			'tbc impl --interface=io.Writer'
			will get some structs like
			type MyWriter struct {}
			func (w *MyWriter) Write(data []byte) (n int, err error)
		`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			// 标志
			flags := cmd.Flags()
			interf, _ := flags.GetString("interface")
			path, _ := flags.GetString("path")
			rec, _ := flags.GetBool("recursive")
			fmt.Printf("| impl | %+v, %+v, %+v\n", interf, path, rec)

			ip := &importpath.ImportPath{}
			paths, err := getPaths(ip, path, rec)
			if err != nil {
				log.Fatal(err)
			}
			if len(paths) == 0 {
				log.Fatalf("找不到有效路径，请使用-p指定或设置-r！")
			}
			fmt.Printf("dirs: %+v\n", paths)

			// 解析
			p := parser.New(parser.Option{
				Op: parser.OpImpl,
			})

			// 获取接口类型
			var interType *types.Interface
			importPath, typ := ip.SplitImportPathWithType(interf)
			inters, err := p.ParseByGoPackages(importPath)
			if err != nil {
				log.Fatal(err)
			}
			debug.Printf("inters: %+v\n", inters)
			for _, pkg := range inters.Pkgs {
				fmt.Printf("pkg: %+v\n", pkg)
				for _, one := range pkg.Interfaces {
					if one.Name != typ {
						continue
					}
					interType = one.Interface

					debug.Printf("interface: %+v\n", one.Interface)
				}
			}

			pkgs, err := p.ParseByGoPackages(paths...)
			if err != nil {
				log.Fatal(err)
			}
			for _, pkg := range pkgs.Pkgs {
				fmt.Printf("find %d structs in %s\n", len(pkg.Structs), pkg.PkgPath)
				for i, one := range pkg.Structs {
					if !types.Implements(one.TypesType, interType) {
						continue
					}
					fmt.Printf("	No.%d struct info: %+v\n", i, one.PkgPath+"."+one.Name)
				}
			}
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   string(parser.OpCallgraph),
		Short: "get the callgraph of a function or method",
		Long: `
	path: specify package path
	func：specify function or method, if method, use 'StructName.MethodName', like：A.GetByName
	ignore: ignore package, std means standart packages, others use themself's package path
	depth: call depth, use it if you want to skip deep call info

like:
	tbc callgraph --path=xxx.xxx.xxx/a/b --func=[main|normal_func|struct_method] --ignore=std;xxx.xxx.xxx/e/f --depth=2
		`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			// 标志
			flags := cmd.Flags()
			path, _ := flags.GetString("path")
			funcName, _ := flags.GetString("func")
			funcParts := strings.Split(funcName, ".")
			ignore, _ := flags.GetString("ignore")
			ignores := strings.Split(ignore, ";")
			depth, _ := flags.GetInt("depth")
			fmt.Printf("| callgraph | %+v, %+v, %+v, %v\n", path, funcParts, ignores, depth)

			// 从指定的函数/方法的定义处开始，途中忽略部分包调用，直到指定深度为止

			// 解析
			p := parser.New(parser.Option{
				Op:       parser.OpCallgraph,
				NeedCall: true,
			})

			// 获取接口类型
			pkgs, err := p.ParseByGoPackages(path)
			if err != nil {
				log.Fatal(err)
			}

			// 唯一键：struct name + method name
			// func name
			funcMap := make(map[string]parser.Func)
			for _, pkg := range pkgs.Pkgs {
				for _, oneFunc := range pkg.Funcs {
					var key = oneFunc.Name
					if oneFunc.Recv != "" {
						key = oneFunc.Recv + "." + oneFunc.Name
					}
					funcMap[key] = oneFunc
				}
			}

			var exist bool
			var targetFunc parser.Func
			for _, pkg := range pkgs.Pkgs {
				for _, oneFunc := range pkg.Funcs {
					debug.Printf("oneFunc: %+v\n", oneFunc)

					if len(funcParts) == 2 {
						if oneFunc.Recv == "" {
							continue
						}
						debug.Printf("%+v, %s, %s\n", funcParts, oneFunc.Recv, oneFunc.Name)
						if funcParts[0] != oneFunc.Recv ||
							funcParts[1] != oneFunc.Name {
							continue
						}
						exist = true
						targetFunc = oneFunc
						break
					} else {
						if oneFunc.Name != funcName {
							continue
						}
						exist = true
						targetFunc = oneFunc
						break
					}
				}
			}
			if !exist {
				log.Fatalf("can't find the func: %s\n", funcName)
			}

			ignoreStd := false
			newIgnores := make([]string, 0, len(ignores))
			for _, ignore := range ignores {
				if ignore == "std" {
					ignoreStd = true
					continue
				}
				newIgnores = append(newIgnores, ignore)
			}
			if ignoreStd {
				stdPkgs := p.GetStandardPackages()
				newIgnores = append(newIgnores, stdPkgs...)
			}

			oneFuncPtr := &targetFunc
			oneFuncPtr.Set(funcMap, depth)
			debug.Printf("=== got: %+v\n", targetFunc)
			targetFunc.PrintCallGraph(newIgnores, depth)
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   string(parser.OpGenProject),
		Short: "gen layout of a new project",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			// 指定一个目录，生成项目结构
			// 指定项目的go module name，没有则使用目录名

			_ = err
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   string(parser.OpGenProxy),
		Short: "gen proxy by @proxy directive from source code, notice: this will change the origin code",
		Long: `
		-p path maybe a directory, or an import path of a package, like: '~/a/b/c', 'bytes', 'github.com/donnol/tools'...
		--func specify func name which will be added '[func]Proxy'
		`,
		Run: func(cmd *cobra.Command, args []string) {
			// 标志
			flags := cmd.Flags()
			path, _ := flags.GetString("path")
			funcName, _ := flags.GetString("func")
			if path == "" {
				panic("path cannot be empty")
			}
			if funcName == "" {
				panic("func cannot be empty")
			}
			fmt.Printf("| %s | %+v, %s\n", parser.OpGenProxy, path, funcName)

			cfg := &packages.Config{
				Mode: packages.NeedName |
					packages.NeedFiles |
					packages.NeedCompiledGoFiles |
					packages.NeedImports |
					packages.NeedDeps |
					packages.NeedExportFile |
					packages.NeedTypes |
					packages.NeedSyntax |
					packages.NeedTypesInfo |
					packages.NeedTypesSizes |
					packages.NeedModule,
			}
			pkgs, err := packages.Load(cfg, path)
			if err != nil {
				return
			}

			for _, pkg := range pkgs {

				file := ""
				if len(pkg.Syntax) > 0 {
					file = pkg.CompiledGoFiles[0]
				}
				genFile := filepath.Join(filepath.Dir(file), "gen_proxy.go")
				f, err := os.OpenFile(genFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
				if err != nil {
					panic(err)
				}
				defer f.Close()

				for i, fileSyntax := range pkg.Syntax {
					file := pkg.CompiledGoFiles[i]
					if file == genFile {
						continue
					}
					// https://blog.microfast.ch/refactoring-go-code-using-ast-replacement-e3cbacd7a331
					// 通过ast替换来修改代码
					// golang.org/x/tools/go/ast/astutil包的astutil.Apply方法

					newNode := astutil.Apply(fileSyntax, func(cr *astutil.Cursor) bool {
						var args []ast.Expr
						ce, ok := cr.Node().(*ast.CallExpr)
						if !ok {
							return true
						}
						args = ce.Args
						ident, isIdent := ce.Fun.(*ast.Ident)
						if !isIdent {
							return true
						}
						// TODO: others: *ast.SelectorExpr, *ast.IndexListExpr

						if ident.Name != funcName {
							return true
						}

						debug.Printf("into replace call expr astutil apply callexpr, args: %+v, name: %s\n", args, ident.Name)
						debug.Printf("c.Index: %d\n", cr.Index())
						newFuncName := ident.Name + "Proxy"

						var argWithType, res, resDefine, resVars, argWithoutType string

						tav, ok := pkg.TypesInfo.Types[ce.Fun]
						if !ok {
							debug.Printf("cannot find types info of arg: %v\n", ce)
							return true
						}
						debug.Printf("find types info of arg: %v, %#v\n", ce, tav)

						sig, ok := tav.Type.(*types.Signature)
						if ok {
							debug.Printf("sig: %#v\n", sig)

							for i := 0; i < sig.Params().Len(); i++ {
								p := sig.Params().At(i)
								debug.Printf("sig, p: %#v\n", p)

								argWithoutType += p.Name()
								if i == sig.Params().Len()-1 && sig.Variadic() {
									slice, ok := p.Type().(*types.Slice)
									if ok {
										argWithType += p.Name() + " ..." + slice.Elem().String()
										argWithoutType += "..."
									}
								} else {
									argWithType += p.Name() + " " + p.Type().String()
								}
								if i != sig.Params().Len()-1 {
									argWithType += ", "
									argWithoutType += ", "
								}
							}
							for i := 0; i < sig.Results().Len(); i++ {
								r := sig.Results().At(i)
								debug.Printf("sig, r: %#v\n", r)

								res += r.Name() + " " + r.Type().String()
								resVarName := "r" + strconv.Itoa(i)
								resDefine += "var " + resVarName + " " + r.Type().String() + "\n"
								resVars += resVarName
								if i != sig.Results().Len()-1 {
									res += ", "
									resVars += ", "
								}
							}
						}

						tmplBuf := new(bytes.Buffer)
						err = tmpl.Execute(tmplBuf, map[string]interface{}{
							"args":           argWithType,
							"res":            res,
							"resDefine":      resDefine,
							"resVars":        resVars,
							"funcName":       ident.Name,
							"argWithoutType": argWithoutType,
							"newFuncName":    newFuncName,
						})
						if err != nil {
							panic(err)
						}
						debug.Printf("tmplBuf: %s\n", tmplBuf.String())

						formatContent, err := format.Format("gen_proxy.go", tmplBuf.String(), false)
						if err != nil {
							panic(err)
						}

						_, err = f.Write([]byte(formatContent))
						if err != nil {
							panic(err)
						}

						// Replace values
						cr.Replace(&ast.CallExpr{
							Fun:      ast.NewIdent(newFuncName),
							Lparen:   ce.Lparen,
							Args:     args,
							Ellipsis: ce.Ellipsis,
							Rparen:   ce.Rparen,
						})
						return false
					}, nil).(*ast.File)

					if debug.IsDebug() {
						spew.Dump(newNode)
					}

					f, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
					if err != nil {
						panic(err)
					}
					printer.Fprint(f, token.NewFileSet(), newNode)
				}
			}
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   string(parser.OpFind),
		Short: "find all module under directory",
		Long: `
		-p path specify a directory like: '~/a/b/c'
		`,
		Run: func(cmd *cobra.Command, args []string) {
			// 标志
			flags := cmd.Flags()
			path, _ := flags.GetString("path")

			fmt.Printf("| %s | %+v\n", parser.OpFind, path)

			p := &importpath.ImportPath{}
			mods, err := p.FindAllModule(path)
			if err != nil {
				panic(err)
			}

			for _, mod := range mods {
				fmt.Printf("module: %s, %s, %s\n", mod.RelDir, mod.Mod.Path, mod.Mod.Version)
			}
		},
	})
}

var (
	proxyTmpl = `
func {{.newFuncName}}({{.args}}) ({{.res}}) {
	begin := time.Now()

	{{.resDefine}}

	{{.resVars}} = {{.funcName}}({{.argWithoutType}})

	log.Printf("used time: %v\n", time.Since(begin))

	return {{.resVars}}
}
`
	tmpl = func() *template.Template {
		tmpl, err := template.New("proxyTmpl").Parse(proxyTmpl)
		if err != nil {
			panic(err)
		}
		return tmpl
	}()
)

func getPaths(ip *importpath.ImportPath, path string, rec bool) ([]string, error) {
	var err error

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var paths []string

	if path == "" {
		path, err = ip.GetByCurrentDir()
		if err != nil {
			return nil, err
		}

		haveGoFile, err := checkDirHaveGoFile(dir)
		if err != nil {
			return nil, err
		}
		if haveGoFile {
			paths = append(paths, path)
		}
	} else {
		// 手动指定的path，不校验是否存在go文件，由用户自己保证
		paths = append(paths, path)
	}

	modDir, modPath, err := ip.GetModFilePath(dir)
	if err != nil {
		return nil, err
	}
	fmt.Printf("dir: %s, modDir: %s, modPath: %s\n", dir, modDir, modPath)

	if rec {
		dirs, err := collectGoFileDir(dir)
		if err != nil {
			return nil, err
		}
		for _, d := range dirs {
			paths = append(paths, strings.Replace(d, dir, modDir, -1))
		}
	}
	return paths, nil
}

// collectGoFileDir 在指定目录下收集含有go文件的子目录
func collectGoFileDir(dir string) ([]string, error) {
	var dirs []string
	if err := filepath.Walk(dir, filepath.WalkFunc(func(childDir string, info os.FileInfo, ierr error) error {
		if ierr != nil {
			fmt.Printf("walk got err: %+v\n", ierr)
		}

		if childDir == dir {
			return nil
		}
		// 获取所需目录
		if !info.IsDir() {
			return nil
		}
		haveGoFile, err := checkDirHaveGoFile(childDir)
		if err != nil {
			return err
		}
		// 过滤没有go文件的
		if !haveGoFile {
			return nil
		}

		dirs = append(dirs, childDir)

		return nil
	})); err != nil {
		return nil, err
	}

	return dirs, nil
}

func checkDirHaveGoFile(dir string) (bool, error) {

	fileInfos, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	haveGoFile := false
	for _, fi := range fileInfos {
		ext := filepath.Ext(fi.Name())
		if ext == ".go" {
			haveGoFile = true
			break
		}
	}

	return haveGoFile, nil
}

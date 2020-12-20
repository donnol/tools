package main

import (
	"fmt"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/donnol/tools/importpath"
	"github.com/donnol/tools/internal/utils/debug"
	"github.com/donnol/tools/parser"

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
	rootCmd.PersistentFlags().StringVarP(&ignore, "ignore", "", "", "specify ignore package")
	var depth int
	rootCmd.PersistentFlags().IntVarP(&depth, "depth", "", 0, "specify depth")

	// 添加子命令
	addSubCommand(rootCmd)

	// 执行
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func addSubCommand(rootCmd *cobra.Command) {
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
			debug.Debug("inters: %+v\n", inters)
			for _, pkg := range inters.Pkgs {
				fmt.Printf("pkg: %+v\n", pkg)
				for _, one := range pkg.Interfaces {
					if one.Name != typ {
						continue
					}
					interType = one.Interface

					debug.Debug("interface: %+v\n", one.Interface)
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
					debug.Debug("oneFunc: %+v\n", oneFunc)

					if len(funcParts) == 2 {
						if oneFunc.Recv == "" {
							continue
						}
						debug.Debug("%+v, %s, %s\n", funcParts, oneFunc.Recv, oneFunc.Name)
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
			debug.Debug("=== got: %+v\n", targetFunc)
			targetFunc.PrintCallGraph(newIgnores, depth)
		},
	})
}

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

	fileInfos, err := ioutil.ReadDir(dir)
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

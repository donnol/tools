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

	// 添加子命令
	addSubCommand(rootCmd)

	// 执行
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func addSubCommand(rootCmd *cobra.Command) {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "interface",
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
			p := parser.New(parser.Option{})
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
		Use:   "replace",
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
		Use:   "mock",
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
		Use:   "impl",
		Short: "TODO: find implement by given interface in specify path",
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
			p := parser.New(parser.Option{})

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

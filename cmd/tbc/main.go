package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/donnol/tools/format"
	"github.com/donnol/tools/importpath"
	"github.com/donnol/tools/parser"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "tbc",
		Short: "a tool named to be continued",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func main() {
	// 解析标签
	var rFlag bool
	rootCmd.PersistentFlags().BoolVarP(&rFlag, "recursive", "r", false, "recursively process dir from current")
	var path string
	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", "", "specify import path")

	// 添加子命令
	addSsubCommand(rootCmd)

	// 执行
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func addSsubCommand(rootCmd *cobra.Command) {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "interface",
		Short: "gen struct interface",
		Long:  "gen struct interface, like: type IM interface {...}.",
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			// 标志
			flags := cmd.Flags()
			path, _ := flags.GetString("path")
			if path == "" {
				ip := &importpath.ImportPath{}
				path, err = ip.GetByCurrentDir()
				if err != nil {
					log.Fatal(err)
				}
			}
			rec, _ := flags.GetBool("recursive")
			fmt.Printf("| interface | %+v, %+v\n", path, rec)

			// 解析
			p := parser.New(parser.Option{
				UseSourceImporter: true,
			})
			structs, err := p.ParseAST(path)
			if err != nil {
				log.Fatal(err)
			}

			// 写入
			var is string = "package " + p.GetPkgName() + "\n"
			for _, single := range structs {
				is += single.MakeInterface() + "\n\n"
			}
			fileName := filepath.Join(p.GetDir(), "interface.go")
			formatContent, err := format.Format(fileName, is, false)
			if err != nil {
				log.Fatalf("err: %+v, content: %s\n", err, is)
			}
			if err = ioutil.WriteFile(fileName, []byte(formatContent), os.ModePerm); err != nil {
				log.Fatal(err)
			}
		},
	})
}

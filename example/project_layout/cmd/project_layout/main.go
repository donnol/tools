// 一个终端命令行应用，里面包括有web子命令、监控子命令等
package main

import (
	"fmt"
	"os"
	"project_layout/model/config"

	"github.com/urfave/cli/v2"
)

const (
	cmdName = "bcwallet"
)

func main() {
	app := cli.NewApp()
	app.Name = cmdName
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			Value:       "./data/conf/project_layout.toml",
			DefaultText: "./data/conf/project_layout.toml",
		},
		&cli.StringFlag{
			Name:  "secret",
			Usage: "specify secret key to encrypt or decrypt the password field of config, like: 1234567890abcdef",
		},
	}
	app.EnableBashCompletion = true

	// 配置子命令
	app.Commands = cmds
	app.Setup()

	env, ok := config.EnvFromOS()
	if ok {
		app.Metadata[config.EnvOSKey] = env
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ooh, something have went wrong!\nERROR: %+v\n\n", err) // nolint:errcheck

		os.Exit(1)
	}
}

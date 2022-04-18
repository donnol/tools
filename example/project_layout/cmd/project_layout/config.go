package main

import (
	"fmt"
	"time"

	"project_layout/model/config"

	"github.com/BurntSushi/toml"
	"github.com/donnol/tools/timex"
	"github.com/urfave/cli/v2"
)

func mustGetConf(cctx *cli.Context) *config.Config {
	conf := config.New(
		config.Option{
			Source: config.SourceFile,
			Setters: []config.Setter{
				// 完全覆盖，需要返回新的conf
				func(c *config.Config) *config.Config {
					file := cctx.String("config")

					conf := &config.Config{}
					_, err := toml.DecodeFile(file, conf)
					if err != nil {
						panic(err)
					}

					return conf
				},
			},
		},
		config.Option{
			Source: config.SourceEnv,
			Setters: []config.Setter{
				// 部分覆盖，需要对传入的c赋值，并返回nil
				func(c *config.Config) *config.Config {

					var exist bool
					env, exist := cctx.App.Metadata[config.EnvOSKey]
					if exist {
						var ok bool
						c.Env, ok = env.(config.Env)
						if !ok {
							panic("LC_WALLET_ENV环境变量设置有误")
						}
					}

					return nil
				},
			},
		},
		config.Option{
			Source: config.SourceFlag,
			Setters: []config.Setter{
				// 部分覆盖，需要对传入的c赋值，并返回nil
				func(conf *config.Config) *config.Config {
					port := cctx.Int("port")
					if port != 0 {
						conf.Server.Port = port
					}

					return nil
				},
			},
		},
	)

	fmt.Printf("%+v\n", conf)

	return conf
}

func printNowTime() {
	now := time.Now()
	fmt.Printf("\nts: %d, time: %s, loc: %v\n", now.Unix(), now.Format(timex.DateTimeFormat), time.Local)
}

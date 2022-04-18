package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/urfave/cli/v2"
)

var (
	cmds = []*cli.Command{
		{
			Name: "server",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name: "port",
				},
			},
			Action: func(ctx *cli.Context) error {
				conf := mustGetConf(ctx)

				engine := gin.Default()
				registerAPI(engine)

				addr := fmt.Sprintf(":%d", conf.Server.Port)
				if err := engine.Run(addr); err != nil {
					return err
				}

				return nil
			},
		},
		{
			Name: "timer",
			Action: func(ctx *cli.Context) error {
				timer := cron.New(
					cron.WithLocation(time.Local),
					cron.WithSeconds(),
				)

				registerTimer(timer)

				timer.Run()

				return nil
			},
		},
	}
)

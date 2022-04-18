package main

import (
	"fmt"
	"log"
	"project_layout/internal/api"
	"project_layout/internal/service"
	"project_layout/internal/store"

	"github.com/donnol/tools/inject"
	"github.com/gin-gonic/gin"
	"github.com/gocraft/dbr/v2"
	"github.com/robfig/cron/v3"

	_ "github.com/go-sql-driver/mysql"
)

var (
	ioc = inject.NewIoc(true)

	apiObj = func() *api.API {
		obj := &api.API{}

		// 注册provider
		if err := registerProvider(ioc); err != nil {
			panic(err)
		}

		// 注入
		if err := ioc.Inject(obj); err != nil {
			panic(err)
		}

		return obj
	}()
)

func registerProvider(ioc *inject.Ioc) error {
	for _, provider := range []interface{}{
		func() *dbr.Connection {
			// TODO: use config
			conn, err := dbr.Open("mysql", "project_layout_man:project_layout_power@tcp(wondko-dms-pet-production.mysql.rds.aliyuncs.com:3306)/project_layout?charset=utf8&parseTime=true&loc=Local", nil)
			if err != nil {
				panic(err)
			}
			return conn
		},
		store.NewUserStore,

		service.NewPingSrv,
		service.NewUserSrv,

		// TODO: add more provider...
	} {
		if err := ioc.RegisterProvider(provider); err != nil {
			return fmt.Errorf("Register provider failed: %w", err)
		}
	}

	return nil
}

func registerAPI(engine *gin.Engine) {
	// api
	api := engine.Group("/api")
	{
		api.GET("/ping", apiObj.Ping())

		// TODO: add more route...
	}
}

func registerTimer(timer *cron.Cron) {
	for _, entry := range []struct {
		spec string
		job  cron.FuncJob
	}{
		{
			spec: "* * * * * *",
			job: cron.FuncJob(
				func() {
					fmt.Printf("test: %s\n", apiObj.PingSrv.Ping())
				},
			),
		},

		// TODO: add more cron entry...
	} {
		id, err := timer.AddFunc(entry.spec, entry.job)
		if err != nil {
			panic(err)
		}
		log.Printf("timer entry id: %v\n", id)
	}
}

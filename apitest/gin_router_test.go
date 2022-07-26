package apitest

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestGinHandlerAPIDoc(t *testing.T) {
	const (
		routePrefix = "/apidoc"
	)

	// 生成两个接口文档
	t.Run("gen doc user", func(t *testing.T) {
		f, err := OpenFile("doc/user.md", "用户模块接口文档")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		if err := NewAT(routePrefix+"/user", http.MethodGet, "获取用户信息", nil, nil).SetParam(&struct {
			Id uint `json:"id"`
		}{
			Id: 1,
		}).FakeRun().
			Result(&struct {
				Code int    `json:"code"`
				Msg  string `json:"msg"`
				User struct {
					Id   uint   `json:"id"`
					Name string `json:"string"`
				}
			}{
				User: struct {
					Id   uint   `json:"id"`
					Name string `json:"string"`
				}{
					Id:   1,
					Name: "jd",
				},
			}).
			Errors(&struct {
				Code int    `json:"code"`
				Msg  string `json:"msg"`
			}{
				Code: 1,
				Msg:  "认证失败",
			}, &struct {
				Code int    `json:"code"`
				Msg  string `json:"msg"`
			}{
				Code: 2,
				Msg:  "校验失败",
			}).
			WriteFile(f).
			Err(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("gen doc book", func(t *testing.T) {
		f, err := OpenFile("doc/book.md", "图书模块接口文档")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		if err := NewAT(routePrefix+"/book", http.MethodGet, "获取图书信息", nil, nil).SetParam(&struct {
			Id uint `json:"id"`
		}{
			Id: 1,
		}).FakeRun().
			Result(&struct {
				Code int    `json:"code"`
				Msg  string `json:"msg"`
				Book struct {
					Id   uint   `json:"id"`
					Name string `json:"string"`
					Page int    `json:"page"`
				}
			}{
				Book: struct {
					Id   uint   `json:"id"`
					Name string `json:"string"`
					Page int    `json:"page"`
				}{
					Id:   1,
					Name: "jd",
					Page: 100,
				},
			}).
			Errors(&struct {
				Code int    `json:"code"`
				Msg  string `json:"msg"`
			}{
				Code: 1,
				Msg:  "认证失败",
			}, &struct {
				Code int    `json:"code"`
				Msg  string `json:"msg"`
			}{
				Code: 2,
				Msg:  "校验失败",
			}).
			WriteFile(f).
			Err(); err != nil {
			t.Fatal(err)
		}
	})

	// 启动服务，注册路由
	buf := new(bytes.Buffer)
	gin.DefaultWriter = io.MultiWriter(buf, os.Stdout)
	engine := gin.Default()
	doc := engine.Group(routePrefix)
	{
		GinHandlerAPIDoc(doc, "doc", "jdlau")
	}
	go func() {
		if err := engine.Run(":8888"); err != nil {
			panic(err)
		}
	}()
	time.Sleep(1 * time.Second)

	// 验证路由是否已注册
	var haveUser, haveBook bool
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, "GET    "+routePrefix+"/book") {
			haveBook = true
		}

		if strings.Contains(text, "GET    "+routePrefix+"/user") {
			haveUser = true
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if !haveBook || !haveUser {
		t.Fatalf("bad result, don't have book or user route")
	}
}

package apitest

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/donnol/tools/log"
	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

const (
	prefixTmpl = `<!DOCTYPE html>
	<html lang="en">
	
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<!-- <meta name="viewport" content="width=device-width, initial-scale=1.0"> -->
		<meta name="viewport" content="width=device-width,minimum-scale=1.0,maximum-scale=1.0,user-scalable=no">
		<title>Document</title>
		<style>
			body,
			html {
				margin: 0;
				padding: 0;
				width: 100%;
				height: 100%;
			}
			.content{
				/* display: flex;
				flex-flow: column;
				justify-content: center;
				align-items: center; */
				width: 50%;
				margin: 0 auto;
				position: relative;
				min-height: calc(100% - 80px);
			}
			#heading{
				text-align: center;
				font-size: 40px;
				font-weight: bold;
				padding: 40px 0;
			}
			.api{
				color: #409EFF;
				font-weight: bolder;
			}
			.foot{
				border-top: 1px;
				height: 80px;
				line-height: 80px;
				text-align: center;
				background: #eee;
				color: #333;
				font-size: 16px;
				width: 100%;
			}
			.codebg{
				background: #2e2e1f;
				color: #fff;
				border-radius: 10px;
				padding: 10px;
			}
			@media screen and (max-width: 500px) {
				.content{
					width: 100%;
					padding: 0 20px;
				}
			}
		</style>
	</head>
	
	<body>
		<div class="content">
`

	suffixTmpl = `</div>
    <!-- raw HTML omitted -->
    <footer class="foot">
		Copyright © 2022, %s. All rights reserved. 
    </footer>
</body>
<script text="javascript">
    function directPath(path) {
        window.location.href = window.location.pathname.replace('/apidoc/index', path)
    }
</script>
</html>
	`

	indexTmpl = `<div class="content">
    <p id="heading">接口文档</p>
    <ul>
    {{range $k, $v := .list}}
        <li><p><a href="javascript:;" onclick="directPath('{{$v.Url}}')">{{$v.Title}}</a></p></li>
    {{end}}
    </ul>
</div>
`
)

type Link struct {
	Url   string
	Title string
}

// GinHandlerAPIDoc 针对指定目录下的md接口文档，生成对应的html文件，并注册到gin路由上
func GinHandlerAPIDoc(doc *gin.RouterGroup, dir string, brand string) {
	log.Infof("[apidoc] apidoc dir: %s", dir)

	// 遍历目录
	var links []Link
	if err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 找出md文件
		fi, err := d.Info()
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		if filepath.Ext(fi.Name()) != ".md" {
			return nil
		}
		fileName := strings.TrimSuffix(filepath.Base(fi.Name()), filepath.Ext(fi.Name()))
		mddata, err := os.ReadFile(filepath.Join(dir, fi.Name()))
		if err != nil {
			return err
		}

		// 获取标题
		var title string
		mddataBuf := bytes.NewBuffer(mddata)
		scanner := bufio.NewScanner(mddataBuf)
		if scanner.Scan() {
			title = scanner.Text()
		}
		title = strings.TrimLeft(title, "# ")

		// 转为html
		md := goldmark.New(
			goldmark.WithExtensions(extension.GFM),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
				parser.WithBlockParsers(),
			),
			goldmark.WithRendererOptions(
				html.WithHardWraps(),
				html.WithUnsafe(),
				html.WithXHTML(),
			),
		)
		var buf bytes.Buffer
		if err := md.Convert(mddata, &buf); err != nil {
			log.Errorf("[apidoc] convert md file to html failed: %v, file: %s", err, fi.Name())
			return err
		}
		htmlFileName := fileName + ".html"
		content := fillContent(buf.Bytes(), brand)
		os.WriteFile(filepath.Join(dir, htmlFileName), content, os.ModePerm)

		// 注册路由
		route := "/" + fileName
		doc.StaticFile(route, filepath.Join(dir, htmlFileName))

		links = append(links, Link{
			Url:   doc.BasePath() + route,
			Title: title,
		})
		return nil
	}); err != nil {
		log.Errorf("[apidoc] walk doc dir %s failed: %v", dir, err)
		return
	}

	// 索引页面
	indexRoute := "index"
	indexFileName := indexRoute + ".html"
	indexFilePath := filepath.Join(dir, indexFileName)
	temp, err := template.New(indexRoute).Parse(indexTmpl)
	if err != nil {
		log.Errorf("[apidoc] parse index template failed: %v", err)
		return
	}
	indexBuf := new(bytes.Buffer)
	if err := temp.ExecuteTemplate(indexBuf, indexRoute, map[string]interface{}{
		"list": links,
	}); err != nil {
		log.Errorf("[apidoc] exec index template failed: %v", err)
		return
	}
	content := fillContent(indexBuf.Bytes(), brand)
	if err := os.WriteFile(indexFilePath, content, os.ModePerm); err != nil {
		log.Errorf("[apidoc] write file failed: %v", err)
		return
	}
	doc.StaticFile("/"+indexRoute, indexFilePath)
}

func fillContent(in []byte, brand string) []byte {
	suffixContent := fmt.Sprintf(suffixTmpl, brand)

	content := []byte(prefixTmpl)
	content = append(content, in...)
	content = append(content, []byte(suffixContent)...)

	return content
}

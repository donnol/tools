package route

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/donnol/tools/log"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
)

// 参数相关
var (
	decoder = schema.NewDecoder()
)

// Param 参数
type Param struct {
	// 方法
	method string

	// 参数
	body   []byte
	values url.Values

	// 文件
	multipartReader *multipart.Reader
}

// Parse 解析
func (p *Param) Parse(ctx context.Context, v interface{}) error {
	var err error

	// 解析
	switch p.method {
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		err = json.Unmarshal(p.body, v)
	case http.MethodGet:
		fallthrough
	case http.MethodDelete:
		err = decoder.Decode(v, p.values)
	}
	if err != nil {
		return errors.WithStack(err)
	}

	// 检查参数
	if vv, ok := v.(Checker); ok {
		if err := vv.Check(ctx); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

// ParseMultipartForm 解析内容
func (p *Param) ParseMultipartForm(maxFileSize int64, v interface{}) (map[string][]byte, error) {
	if p.multipartReader == nil {
		return nil, fmt.Errorf("Bad multipart reader")
	}

	// 使用ReadForm
	form, err := p.multipartReader.ReadForm(maxFileSize)
	if err != nil {
		return nil, err
	}

	// 获取参数
	if err := decoder.Decode(v, form.Value); err != nil {
		return nil, err
	}

	// 获取内容
	body := make(map[string][]byte)
	buf := new(bytes.Buffer)
	for name, single := range form.File {
		for i, one := range single {
			file, err := one.Open()
			if err != nil {
				return nil, err
			}
			defer file.Close()

			_, err = buf.ReadFrom(file)
			if err != nil {
				return nil, err
			}
			log.Default().Debugf("No.%d, name: %s, content: %s\n", i, name, buf.Bytes())
		}
		body[name] = buf.Bytes()

		buf = new(bytes.Buffer) // 不能用 buf.Reset()，因为在下次写入的数据长度比本次的长度小时，会复用之前的内存，导致内容错误
	}

	return body, nil
}

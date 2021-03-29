package route

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"

	utilerrors "github.com/donnol/tools/errors"
)

// Result 结果
type Result struct {
	utilerrors.Error

	Data interface{} `json:"data"` // 正常返回时的数据

	// 给登陆接口使用
	CookieAfterLogin int `json:"-"` // 登陆时需要设置登陆态的用户信息

	// 时间戳
	Timestamp int64 `json:"timestamp"`

	// 请求ID，在请求到来时生成，处理过程传递，返回时一并返回
	RequestID string `json:"requestID"` // uuid

	// 下载内容时使用
	Content
}

// Content 内容
type Content struct {
	ContentLength int64             `json:"-"`
	ContentType   string            `json:"-"`
	ContentReader io.Reader         `json:"-"`
	ExtraHeaders  map[string]string `json:"-"`
}

// MakeContentFromBuffer 新建内容
func MakeContentFromBuffer(filename string, buf *bytes.Buffer) Content {
	var r Content

	writer := multipart.NewWriter(buf)
	r.ContentLength = int64(buf.Len())
	r.ContentType = writer.FormDataContentType()
	r.ContentReader = buf
	r.ExtraHeaders = map[string]string{
		ContentDispositionHeaderKey: fmt.Sprintf(
			ContentDispositionHeaderValueFormat,
			filename,
		),
	}

	return r
}

func MakeContentFromBytes(filename string, content []byte) (Content, error) {
	var r Content

	buf := new(bytes.Buffer)
	_, err := buf.Write(content)
	if err != nil {
		return r, err
	}
	writer := multipart.NewWriter(buf)
	r.ContentLength = int64(buf.Len())
	r.ContentType = writer.FormDataContentType()
	r.ContentReader = buf
	r.ExtraHeaders = map[string]string{
		ContentDispositionHeaderKey: fmt.Sprintf(
			ContentDispositionHeaderValueFormat,
			filename,
		),
	}

	return r, nil
}

// PresentData 用具体结构体展现数据
func (r *Result) PresentData(v interface{}) error {
	b, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, v); err != nil {
		return err
	}

	return nil
}

// AddResult 添加记录后的结果
type AddResult struct {
	ID int `json:"id"` // 新纪录ID
}

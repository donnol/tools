package route

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"reflect"

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

var (
	version = 0
)

func SetResultVersion(ver int) {
	version = ver
}

type Result1 struct {
	Header ResultCode  `json:"header"`
	Data   interface{} `json:"data"` // 可以是对象或数组
}

type ResultCode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Count   int    `json:"count"`
}

// From https://stackoverflow.com/questions/18287242/locking-an-object-during-json-marshal-in-go
// 通过定义一个以Result为基础的新类型，然后在MarshalJSON方法里做一个类型转换，从而避免调json.Marshal方法时无限循环
type resultHelper Result

func (r Result) MarshalJSON() ([]byte, error) {
	var rd interface{}
	switch version {
	case 1:
		rd = r.ToResult1()
	default:
		rd = resultHelper(r)
	}
	return json.Marshal(rd)
}

func (r *Result) ToResult1() Result1 {
	var r1 Result1
	r1.Header.Code = r.Code
	r1.Header.Message = "OK"
	if r.Msg != "" {
		r1.Header.Message = r.Msg
	}
	r1.Data = r.Data

	r1.Header.Count = 1
	rv := reflect.ValueOf(r.Data)
	rvt := rv.Type()
	if rvt.Kind() != reflect.Slice && rvt.Kind() != reflect.Struct {
		panic(fmt.Errorf("not support result data type: %v", rvt))
	}
	switch rvt.Kind() {
	case reflect.Slice:
		r1.Header.Count = rv.Len()
	case reflect.Struct:
		if rv.IsZero() {
			r1.Header.Count = 0
		}
	}

	return r1
}

// AddResult 添加记录后的结果
type AddResult struct {
	ID int `json:"id"` // 新纪录ID
}

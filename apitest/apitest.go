// Package apitest Usage:
// 	NewAT(xxx).
// 		SetParam(xxx).
// 		Debug().
// 		Run().
// 		EqualCode(xxx).
// 		Result(xxx).
// 		Equal(...).
// 		WriteFile(xxx).
// 		Err()
package apitest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/donnol/tools/worker"
)

// Predefined error
var (
	// ErrNilParam 参数为nil
	ErrNilParam = errors.New("Please input param, param is nil now")
)

// AT api test
type AT struct {
	// 服务器配置
	scheme string
	host   string
	port   string
	url    url.URL

	// 请求相关
	authHeaderKey   string
	authHeaderValue string
	path            string
	method          string
	comment         string
	header          http.Header
	cookies         []*http.Cookie
	param           interface{}
	result          interface{}
	ates            []ate

	// 请求和响应
	req     *http.Request
	reqBody []byte
	resp    *http.Response

	// 文档
	doc string

	// 调试
	debug bool

	// 慢请求数量
	slowNum int

	// 是否批量压力测试中
	isPressureBatch bool

	err error
}

// ate api test错误
type ate struct {
	Code int    `json:"Code"` // 错误码
	Msg  string `json:"msg"`  // 错误信息
}

// NewAT 新建
func NewAT(
	path,
	method,
	comment string,
	h http.Header,
	cookies []*http.Cookie,
) *AT {
	return &AT{
		path:    path,
		method:  method,
		comment: comment,
		header:  h,
		cookies: cookies,
	}
}

// New 克隆一个新的AT
func (at *AT) New() *AT {
	return at.clone()
}

// SetPort 设置端口，如":8080"
func (at *AT) SetPort(port string) *AT {
	at.port = port
	return at
}

// SetHeader 设置header
func (at *AT) SetHeader(header http.Header) *AT {
	at.header = header
	return at
}

func (at *AT) MarkAuthHeader(authHeaderKey, authHeaderValue string) *AT {
	at.authHeaderKey = authHeaderKey
	at.authHeaderValue = authHeaderValue
	return at
}

// SetCookies 设置cookies
func (at *AT) SetCookies(cookies []*http.Cookie) *AT {
	at.cookies = cookies
	return at
}

// SetParam 设置参数
func (at *AT) SetParam(param interface{}) *AT {
	if param == nil {
		at.setErr(fmt.Errorf("nil param"))
		return at
	}

	at.param = param
	return at
}

// Run 运行
func (at *AT) Run() *AT {
	return at.run(true)
}

// Run 运行
func (at *AT) FakeRun() *AT {
	return at.run(false)
}

// MonkeyRun 猴子运行
func (at *AT) MonkeyRun() *AT {
	if at.param == nil {
		at.setErr(ErrNilParam)
		return at
	}

	// 根据参数结构体随机生成测试值
	param, err := structRandomValue(at.param)
	if err != nil {
		at.setErr(err)
		return at
	}
	if at.debug { // 打印随机值
		JSONIndent(os.Stdout, param)
	}
	at.param = param

	return at.run(true)
}

// PressureRun 压力运行，n: 运行次数，c: 并发数
func (at *AT) PressureRun(n, c int) *AT {
	w := worker.New(c)
	w.Start()

	// 记录开始时间
	before := time.Now()

	var total int64
	for i := 0; i < n; i++ {
		if err := w.Push(worker.MakeJob(func() error {
			// 运行
			at.run(true)

			// 统计数量
			atomic.AddInt64(&total, 1)

			return nil
		}, 0, nil)); err != nil {
			at.setErr(err)
		}
	}

	w.Stop()

	// 记录结束时间，并计算耗时
	after := time.Now()
	used := after.Unix() - before.Unix()
	var avg int64
	if used != 0 {
		avg = total / used
	}
	fmt.Printf("\n=== Pressure Report ===\nNumber: %d\nConcurrency: %d\nCompleted: %d\nUsed time: %ds\nRPS: %v\n=== END ===\n\n", n, c, total, used, avg)

	return at
}

// PressureParam 压力测试参数
type PressureParam struct {
	N int // 运行次数
	C int // 并发数
}

// PressureRunBatch 批量压力运行
func (at *AT) PressureRunBatch(param []PressureParam) *AT {
	at.isPressureBatch = true

	for _, single := range param {
		at = at.PressureRun(single.N, single.C)
	}

	fmt.Printf("slowNum is %d\n", at.slowNum)
	at.isPressureBatch = false

	return at
}

// Debug 开启调试模式
func (at *AT) Debug() *AT {
	at.debug = true
	return at
}

// EqualCode 比较响应码
func (at *AT) EqualCode(wantCode int) *AT {
	// 复制resp.Body数据
	data, _, err := copyResponseBody(at.resp)
	if err != nil {
		at.setErr(err)
		return at
	}

	// 校验响应码
	if at.resp.StatusCode == wantCode {
		return at
	}

	at.setErr(fmt.Errorf("bad status code, got %+v\ndata is %s", at.resp, data))
	return at
}

// Result 获取结果
func (at *AT) Result(r interface{}) *AT {
	if r == nil {
		at.setErr(fmt.Errorf("nil r"))
		return at
	}

	// 复制resp.Body
	if at.resp != nil {
		data, _, err := copyResponseBody(at.resp)
		if err != nil {
			at.setErr(err)
			return at
		}

		// 解析data到r
		if err := json.Unmarshal(data, r); err != nil {
			at.setErr(err)
			return at
		}
	}
	at.result = r

	at.jsonIndent(os.Stdout, r)

	return at
}

// Equal 校验
func (at *AT) Equal(args ...interface{}) *AT {
	l := len(args)
	d := l % 2
	if d != 0 {
		at.setErr(fmt.Errorf("Please Input Double Args: %v", args))
		return at
	}
	for i := 0; i < l; i += 2 {
		if !reflect.DeepEqual(args[i], args[i+1]) {
			at.setErr(fmt.Errorf("No.%d Not Equal, Have %v, Want %v", i/2+1, args[i], args[i+1]))
			return at
		}
	}

	return at
}

// EqualThen 相等之后
func (at *AT) EqualThen(f func(*AT) error, args ...interface{}) *AT {
	// 先比较args
	at = at.Equal(args...)
	if at.err != nil {
		return at
	}

	// 成功之后才继续运行f
	if err := f(at); err != nil {
		at.setErr(err)
		return at
	}

	return at
}

// WriteFile 写入文件
func (at *AT) WriteFile(w io.Writer) *AT {
	if w == nil {
		at.setErr(fmt.Errorf("nil writer"))
		return at
	}

	if at.doc == "" {
		at.makeDoc() // 尝试一次生成文档
	}

	if at.doc == "" {
		at.setErr(fmt.Errorf("Empty doc"))
		return at
	}

	if _, err := w.Write([]byte(at.doc)); err != nil {
		at.setErr(err)
		return at
	}
	return at
}

// Err 获取error
func (at *AT) Err() error {
	return at.err
}

// === Private method ===

func (at *AT) makeURL() *AT {
	// 默认值
	scheme := "http"
	host := "localhost"
	port := ":80"

	if at.scheme != "" {
		scheme = at.scheme
	}
	if at.host != "" {
		host = at.host
	}
	if at.port != "" {
		port = at.port
	}

	at.url = url.URL{
		Scheme: scheme,
		Host:   host + port,
		Path:   at.path,
	}

	return at
}

func (at *AT) run(realDo bool) *AT {
	// 请求链接
	at = at.makeURL()
	u := at.url

	// 参数处理
	var body = new(bytes.Buffer)
	switch at.method {
	case http.MethodGet, http.MethodDelete:
		q := u.Query()
		params, err := structToMap(at.param)
		if err != nil {
			at.setErr(err)
			return at
		}
		var valueStr string
		for key, value := range params {
			switch v := value.(type) { // 类型断言，既不能用逗号分隔，也不可用fallthrough
			case []int: // 整型数组
				for _, s := range v {
					valueStr = fmt.Sprintf("%v", s)
					q.Add(key, valueStr)
				}
			case []string: // 字符串数组
				for _, s := range v {
					valueStr = fmt.Sprintf("%v", s)
					q.Add(key, valueStr)
				}
			default:
				valueStr = fmt.Sprintf("%v", value)
				q.Add(key, valueStr)
			}
		}
		u.RawQuery = q.Encode()
	case http.MethodPost, http.MethodPut:
		paramBytes, err := json.Marshal(at.param)
		if err != nil {
			at.setErr(err)
			return at
		}
		_, err = body.Write(paramBytes)
		if err != nil {
			at.setErr(err)
			return at
		}
	default:
		at.setErr(fmt.Errorf("not support method %s", at.method))
		return at
	}

	// 复制一份请求body
	reqBody := make([]byte, body.Len())
	copy(reqBody, body.Bytes())
	at.reqBody = reqBody

	// 新建请求
	req, err := http.NewRequest(at.method, u.String(), body)
	if err != nil {
		at.setErr(err)
		return at
	}

	// 设置header
	for headerKey, headerValue := range map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	} {
		req.Header.Set(headerKey, headerValue)
	}
	for k, v := range at.header {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}

	// 添加cookie, 支持设置多个
	for _, c := range at.cookies {
		req.AddCookie(c)
	}
	at.req = req

	// 发起请求
	// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
		MaxIdleConns:        100, // 最大空闲连接数
		MaxIdleConnsPerHost: 100, // 每个域名最大空闲连接数
	}
	client := &http.Client{
		Timeout:   time.Second * 10, // 超时
		Transport: transport,
	}
	if realDo {
		beforeDo := time.Now()
		resp, err := client.Do(req)
		if err != nil {
			at.setErr(err)
			return at
		}
		afterDo := time.Now()
		used := afterDo.UnixNano() - beforeDo.UnixNano()
		if used >= 1000000000 { // 不小于1s
			if at.isPressureBatch { // 统计数量
				at.slowNum++
			} else {
				fmt.Printf("WARNING: '%s' is slow, used %d ms\n", u.String(), used/1000000)
			}
		}

		// https://stackoverflow.com/questions/17948827/reusing-http-connections-in-golang
		// 只要不关闭response，client就不会重用连接，而是新建连接
		at.resp = resp
	}

	// 收集错误码
	at = at.collectATE()

	return at
}

// 收集错误码
func (at *AT) collectATE() *AT {
	// 复制resp.Body
	if at.resp != nil {
		data, _, err := copyResponseBody(at.resp)
		if err != nil {
			at.setErr(err)
			return at
		}
		tmpATE := ate{}
		if err := json.Unmarshal(data, &tmpATE); err != nil {
			at.setErr(err)
			return at
		}
		if tmpATE.Code != 0 {
			var exist bool
			for _, e := range at.ates {
				if tmpATE.Code == e.Code {
					exist = true
					break
				}
			}
			if !exist {
				at.ates = append(at.ates, tmpATE)
			}
		}
	}

	return at
}

// 生成文档
func (at *AT) makeDoc() *AT {
	const paramName = "Param"
	const returnName = "Return"
	const errorName = "Error"
	var doc string

	// 保存请求和响应
	key := apiKey(at.path, at.method)

	// 标题
	doc += "## " + at.comment + "\n\n"

	// 方法
	doc += "`" + key + "`\n\n"

	// req header
	h := "Request header:\n"
	for k, v := range at.req.Header {
		if k != "Content-Type" && k != at.authHeaderKey {
			continue
		}
		v1 := ""
		if len(v) > 0 {
			v1 = v[0]
		}
		if k == at.authHeaderKey && at.authHeaderValue != "" {
			v1 = at.authHeaderValue
		}
		h += fmt.Sprintf("- '%s': %s\n", k, v1)
	}
	doc += h + "\n"

	// resp header
	resph := "Response header:\n"
	if at.resp != nil {
		for k, v := range at.resp.Header {
			if k != "Content-Type" {
				continue
			}
			v1 := ""
			if len(v) > 0 {
				v1 = v[0]
			}
			resph += fmt.Sprintf("- '%s': %s\n", k, v1)
		}
		resph += "\n"
	} else {
		resph += "- Content-Type: application/json; charset=utf-8\n\n"
	}
	doc += resph

	// 参数
	block, err := structToBlock(paramName, at.param)
	if err != nil {
		at.setErr(err)
		return at
	}
	doc += block

	// 返回
	block, err = structToBlock(returnName, at.result)
	if err != nil {
		at.setErr(err)
		return at
	}
	doc += block

	// 错误码
	if len(at.ates) > 0 {
		block, err = structToList(errorName, at.ates...)
		if err != nil {
			at.setErr(err)
			return at
		}
		doc += block
	}

	// 参数和返回示例
	switch at.method {
	case http.MethodGet, http.MethodDelete:
		doc += dataToSummary(paramName, []byte(at.req.URL.RawQuery), false)
	case http.MethodPost, http.MethodPut:
		doc += dataToSummary(paramName, at.reqBody, true)
	}

	// 复制resp.Body
	var data []byte
	if at.resp != nil {
		data, _, err = copyResponseBody(at.resp)
		if err != nil {
			at.setErr(err)
			return at
		}
	} else {
		data, err = json.Marshal(at.result)
		if err != nil {
			at.setErr(err)
			return at
		}
	}
	doc += dataToSummary(returnName, data, true)

	at.doc = doc

	return at
}

func (at *AT) setErr(err error) *AT {
	if at.err == nil {
		at.err = err
	}
	return at
}

func (at *AT) jsonIndent(w io.Writer, r interface{}) *AT {
	if at.debug {
		JSONIndent(w, r)
	}
	return at
}

func (at *AT) clone() *AT {
	return NewAT(at.path, at.method, at.comment, at.header, at.cookies)
}

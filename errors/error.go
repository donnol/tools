package errors

import "fmt"

// 错误级别
const (
	LevelNormal = 1
	LevelFatal  = 2
)

// 错误码
const (
	ErrorCodeRouter = 10001 // 路由问题
	ErrorCodeAuth   = 10010 // 认证问题
)

// Error 错误
type Error struct {
	Code int    `json:"code"` // 请求返回码，一般0表示正常，非0表示异常
	Msg  string `json:"msg"`  // 信息，一般是出错时的描述信息

	level int `json:"-"` // 级别
}

func newError(code int, msg string, level int) error {
	return Error{
		Code:  code,
		Msg:   msg,
		level: level,
	}
}

// NewNormal 新建普通错误
func NewNormal(code int, msg string) error {
	return newError(code, msg, LevelNormal)
}

// NewFatal 新建严重错误
func NewFatal(code int, msg string) error {
	return newError(code, msg, LevelFatal)
}

// Error 实现error接口
func (e Error) Error() string {
	return fmt.Sprintf("[%s] Code: %d, Msg: %s", e.nameByLevel(), e.Code, e.Msg)
}

// IsNormal 是否普通错误
func (e Error) IsNormal() bool {
	return e.level == LevelNormal
}

// IsFatal 是否严重错误
func (e Error) IsFatal() bool {
	return e.level == LevelFatal
}

func (e Error) nameByLevel() string {
	switch e.level {
	case LevelNormal:
		return "Normal"
	case LevelFatal:
		return "Fatal"
	}
	return ""
}

var _ error = Error{}

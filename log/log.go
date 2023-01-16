package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
)

const (
	calldepth = 3
	skip      = 5
)

// Notifier 通知接口
type Notifier interface {
	Levels() []Level
	Notify(msg string)
}

// Logger 日志接口
type Logger interface {
	SetNotify(notify Notifier)
	Fatalf(format string, v ...any) // Mark the level fatal, but not panic or exit
	Errorf(format string, v ...any)
	Warnf(format string, v ...any)
	Infof(format string, v ...any)
	Debugf(format string, v ...any)
	Tracef(format string, v ...any)
}

// logger 日志
type logger struct {
	*log.Logger

	notify Notifier
}

var defaultLogger = New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

// Default 默认
func Default() Logger {
	return defaultLogger
}

// New 新建
func New(out io.Writer, prefix string, flag int) Logger {
	return &logger{
		Logger: log.New(out, prefix+" ", flag),
	}
}

// SetNotify 设置通知
func (l *logger) SetNotify(notify Notifier) {
	l.notify = notify
}

// Fatalf 致命
func (l *logger) Fatalf(format string, v ...any) {
	level := FatalLevel
	l.printf(level, format, v...)
}

// Errorf 错误
func (l *logger) Errorf(format string, v ...any) {
	level := ErrorLevel
	l.printf(level, format, v...)
}

// Warnf 警告
func (l *logger) Warnf(format string, v ...any) {
	level := WarnLevel
	l.printf(level, format, v...)
}

// Infof 信息
func (l *logger) Infof(format string, v ...any) {
	level := InfoLevel
	l.printf(level, format, v...)
}

// Debugf 调试
func (l *logger) Debugf(format string, v ...any) {
	level := DebugLevel
	l.printf(level, format, v...)
}

// Tracef 追踪
func (l *logger) Tracef(format string, v ...any) {
	level := TraceLevel
	l.printf(level, format, v...)
}

func (l *logger) printf(level Level, format string, v ...any) {
	format = getFormat(level, format)
	msg := fmt.Sprintf(format, v...)

	if err := l.Logger.Output(calldepth, msg); err != nil {
		panic(err)
	}

	// 发送通知
	l.notice(level, msg)
}

// 发送通知
func (l *logger) notice(level Level, msg string) {
	if l.notify == nil {
		return
	}

	levels := l.notify.Levels()
	if !InLevel(levels, level) {
		return
	}

	stack := collectStack()
	l.notify.Notify(msg + stack)
}

// Fatalf 致命
func Fatalf(format string, v ...any) {
	level := FatalLevel
	printf(level, format, v...)
}

// Errorf 错误
func Errorf(format string, v ...any) {
	level := ErrorLevel
	printf(level, format, v...)
}

// Warnf 警告
func Warnf(format string, v ...any) {
	level := WarnLevel
	printf(level, format, v...)
}

// Infof 信息
func Infof(format string, v ...any) {
	level := InfoLevel
	printf(level, format, v...)
}

// Debugf 调试
func Debugf(format string, v ...any) {
	level := DebugLevel
	printf(level, format, v...)
}

// Tracef 追踪
func Tracef(format string, v ...any) {
	level := TraceLevel
	printf(level, format, v...)
}

func getFormat(level Level, format string) string {
	return fmt.Sprintf("[%s] %s", level, format)
}

func printf(level Level, format string, v ...any) {
	format = getFormat(level, format)
	msg := fmt.Sprintf(format, v...)

	if err := log.Output(calldepth, msg); err != nil {
		panic(err)
	}
}

func collectStack() string {
	var stack string

	// 收集栈信息
	var pcs = make([]uintptr, 10)
	n := runtime.Callers(skip, pcs)
	if n == 0 {
		return stack
	}
	pcs = pcs[:n]
	frames := runtime.CallersFrames(pcs)
	for {
		next, more := frames.Next()
		if !more {
			break
		}
		stack += fmt.Sprintf("%v:%v\n%v\n", next.File, next.Line, next.Function)
	}

	return fmt.Sprintf("\n%s\n", stack)
}

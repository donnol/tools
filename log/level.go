package log

// Level 日志等级
type Level string

// 日志级别--https://stackoverflow.com/questions/2031163/when-to-use-the-different-log-levels
const (
	FatalLevel Level = "FATAL"
	ErrorLevel Level = "ERROR"
	WarnLevel  Level = "WARN"
	InfoLevel  Level = "INFO"
	DebugLevel Level = "DEBUG"
	TraceLevel Level = "TRACE"
)

var (
	allLevel = []Level{
		FatalLevel,
		ErrorLevel,
		WarnLevel,
		InfoLevel,
		DebugLevel,
		TraceLevel,
	}
)

func (l Level) All() []Level {
	return allLevel
}

// InLevel 是否在指定Level里
func InLevel(levels []Level, level Level) bool {
	for _, l := range levels {
		if l == level {
			return true
		}
	}

	return false
}

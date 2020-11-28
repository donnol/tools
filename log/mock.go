package log

type NotifierMock struct {
	LevelsFunc func() []Level

	NotifyFunc func(msg string)
}

var _ Notifier = &NotifierMock{}

func (mockRecv *NotifierMock) Levels() []Level {
	return mockRecv.LevelsFunc()
}

func (mockRecv *NotifierMock) Notify(msg string) {
	mockRecv.NotifyFunc(msg)
}

type LoggerMock struct {
	DebugfFunc func(format string, v ...interface{})

	ErrorfFunc func(format string, v ...interface{})

	FatalfFunc func(format string, v ...interface{})

	InfofFunc func(format string, v ...interface{})

	SetNotifyFunc func(notify Notifier)

	TracefFunc func(format string, v ...interface{})

	WarnfFunc func(format string, v ...interface{})
}

var _ Logger = &LoggerMock{}

func (mockRecv *LoggerMock) Debugf(format string, v ...interface{}) {
	mockRecv.DebugfFunc(format, v...)
}

func (mockRecv *LoggerMock) Errorf(format string, v ...interface{}) {
	mockRecv.ErrorfFunc(format, v...)
}

func (mockRecv *LoggerMock) Fatalf(format string, v ...interface{}) {
	mockRecv.FatalfFunc(format, v...)
}

func (mockRecv *LoggerMock) Infof(format string, v ...interface{}) {
	mockRecv.InfofFunc(format, v...)
}

func (mockRecv *LoggerMock) SetNotify(notify Notifier) {
	mockRecv.SetNotifyFunc(notify)
}

func (mockRecv *LoggerMock) Tracef(format string, v ...interface{}) {
	mockRecv.TracefFunc(format, v...)
}

func (mockRecv *LoggerMock) Warnf(format string, v ...interface{}) {
	mockRecv.WarnfFunc(format, v...)
}

type LevelMock struct {
	AllFunc func() []Level
}

var _ ILevel = &LevelMock{}

func (mockRecv *LevelMock) All() []Level {
	return mockRecv.AllFunc()
}

type loggerMock struct {
	DebugfFunc func(format string, v ...interface{})

	ErrorfFunc func(format string, v ...interface{})

	FatalfFunc func(format string, v ...interface{})

	InfofFunc func(format string, v ...interface{})

	SetNotifyFunc func(notify Notifier)

	TracefFunc func(format string, v ...interface{})

	WarnfFunc func(format string, v ...interface{})
}

var _ Ilogger = &loggerMock{}

func (mockRecv *loggerMock) Debugf(format string, v ...interface{}) {
	mockRecv.DebugfFunc(format, v...)
}

func (mockRecv *loggerMock) Errorf(format string, v ...interface{}) {
	mockRecv.ErrorfFunc(format, v...)
}

func (mockRecv *loggerMock) Fatalf(format string, v ...interface{}) {
	mockRecv.FatalfFunc(format, v...)
}

func (mockRecv *loggerMock) Infof(format string, v ...interface{}) {
	mockRecv.InfofFunc(format, v...)
}

func (mockRecv *loggerMock) SetNotify(notify Notifier) {
	mockRecv.SetNotifyFunc(notify)
}

func (mockRecv *loggerMock) Tracef(format string, v ...interface{}) {
	mockRecv.TracefFunc(format, v...)
}

func (mockRecv *loggerMock) Warnf(format string, v ...interface{}) {
	mockRecv.WarnfFunc(format, v...)
}

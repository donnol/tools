package log

type LoggerMock struct {
	DebugfFunc func(format string, v ...any)

	ErrorfFunc func(format string, v ...any)

	FatalfFunc func(format string, v ...any)

	InfofFunc func(format string, v ...any)

	SetNotifyFunc func(notify Notifier)

	TracefFunc func(format string, v ...any)

	WarnfFunc func(format string, v ...any)
}

var _ Logger = &LoggerMock{}

func (mockRecv *LoggerMock) Debugf(format string, v ...any) {
	mockRecv.DebugfFunc(format, v...)
}

func (mockRecv *LoggerMock) Errorf(format string, v ...any) {
	mockRecv.ErrorfFunc(format, v...)
}

func (mockRecv *LoggerMock) Fatalf(format string, v ...any) {
	mockRecv.FatalfFunc(format, v...)
}

func (mockRecv *LoggerMock) Infof(format string, v ...any) {
	mockRecv.InfofFunc(format, v...)
}

func (mockRecv *LoggerMock) SetNotify(notify Notifier) {
	mockRecv.SetNotifyFunc(notify)
}

func (mockRecv *LoggerMock) Tracef(format string, v ...any) {
	mockRecv.TracefFunc(format, v...)
}

func (mockRecv *LoggerMock) Warnf(format string, v ...any) {
	mockRecv.WarnfFunc(format, v...)
}

type loggerMock struct {
	DebugfFunc func(format string, v ...any)

	ErrorfFunc func(format string, v ...any)

	FatalfFunc func(format string, v ...any)

	InfofFunc func(format string, v ...any)

	SetNotifyFunc func(notify Notifier)

	TracefFunc func(format string, v ...any)

	WarnfFunc func(format string, v ...any)
}

var _ Ilogger = &loggerMock{}

func (mockRecv *loggerMock) Debugf(format string, v ...any) {
	mockRecv.DebugfFunc(format, v...)
}

func (mockRecv *loggerMock) Errorf(format string, v ...any) {
	mockRecv.ErrorfFunc(format, v...)
}

func (mockRecv *loggerMock) Fatalf(format string, v ...any) {
	mockRecv.FatalfFunc(format, v...)
}

func (mockRecv *loggerMock) Infof(format string, v ...any) {
	mockRecv.InfofFunc(format, v...)
}

func (mockRecv *loggerMock) SetNotify(notify Notifier) {
	mockRecv.SetNotifyFunc(notify)
}

func (mockRecv *loggerMock) Tracef(format string, v ...any) {
	mockRecv.TracefFunc(format, v...)
}

func (mockRecv *loggerMock) Warnf(format string, v ...any) {
	mockRecv.WarnfFunc(format, v...)
}

type LevelMockMock struct {
	AllFunc func() []Level
}

var _ ILevelMock = &LevelMockMock{}

func (mockRecv *LevelMockMock) All() []Level {
	return mockRecv.AllFunc()
}

type loggerMockMock struct {
	DebugfFunc func(format string, v ...any)

	ErrorfFunc func(format string, v ...any)

	FatalfFunc func(format string, v ...any)

	InfofFunc func(format string, v ...any)

	SetNotifyFunc func(notify Notifier)

	TracefFunc func(format string, v ...any)

	WarnfFunc func(format string, v ...any)
}

var _ IloggerMock = &loggerMockMock{}

func (mockRecv *loggerMockMock) Debugf(format string, v ...any) {
	mockRecv.DebugfFunc(format, v...)
}

func (mockRecv *loggerMockMock) Errorf(format string, v ...any) {
	mockRecv.ErrorfFunc(format, v...)
}

func (mockRecv *loggerMockMock) Fatalf(format string, v ...any) {
	mockRecv.FatalfFunc(format, v...)
}

func (mockRecv *loggerMockMock) Infof(format string, v ...any) {
	mockRecv.InfofFunc(format, v...)
}

func (mockRecv *loggerMockMock) SetNotify(notify Notifier) {
	mockRecv.SetNotifyFunc(notify)
}

func (mockRecv *loggerMockMock) Tracef(format string, v ...any) {
	mockRecv.TracefFunc(format, v...)
}

func (mockRecv *loggerMockMock) Warnf(format string, v ...any) {
	mockRecv.WarnfFunc(format, v...)
}

type NotifierMockMock struct {
	LevelsFunc func() []Level

	NotifyFunc func(msg string)
}

var _ INotifierMock = &NotifierMockMock{}

func (mockRecv *NotifierMockMock) Levels() []Level {
	return mockRecv.LevelsFunc()
}

func (mockRecv *NotifierMockMock) Notify(msg string) {
	mockRecv.NotifyFunc(msg)
}

type LoggerMockMock struct {
	DebugfFunc func(format string, v ...any)

	ErrorfFunc func(format string, v ...any)

	FatalfFunc func(format string, v ...any)

	InfofFunc func(format string, v ...any)

	SetNotifyFunc func(notify Notifier)

	TracefFunc func(format string, v ...any)

	WarnfFunc func(format string, v ...any)
}

var _ ILoggerMock = &LoggerMockMock{}

func (mockRecv *LoggerMockMock) Debugf(format string, v ...any) {
	mockRecv.DebugfFunc(format, v...)
}

func (mockRecv *LoggerMockMock) Errorf(format string, v ...any) {
	mockRecv.ErrorfFunc(format, v...)
}

func (mockRecv *LoggerMockMock) Fatalf(format string, v ...any) {
	mockRecv.FatalfFunc(format, v...)
}

func (mockRecv *LoggerMockMock) Infof(format string, v ...any) {
	mockRecv.InfofFunc(format, v...)
}

func (mockRecv *LoggerMockMock) SetNotify(notify Notifier) {
	mockRecv.SetNotifyFunc(notify)
}

func (mockRecv *LoggerMockMock) Tracef(format string, v ...any) {
	mockRecv.TracefFunc(format, v...)
}

func (mockRecv *LoggerMockMock) Warnf(format string, v ...any) {
	mockRecv.WarnfFunc(format, v...)
}

type LevelMock struct {
	AllFunc func() []Level
}

var _ ILevel = &LevelMock{}

func (mockRecv *LevelMock) All() []Level {
	return mockRecv.AllFunc()
}

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

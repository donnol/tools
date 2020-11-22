package log

type LoggerMock struct {
	SetNotifyFunc func(notify Notifier)

	FatalfFunc func(format string, v ...interface{})

	ErrorfFunc func(format string, v ...interface{})

	WarnfFunc func(format string, v ...interface{})

	InfofFunc func(format string, v ...interface{})

	DebugfFunc func(format string, v ...interface{})

	TracefFunc func(format string, v ...interface{})
}

var _ Logger = &LoggerMock{}

func (*LoggerMock) SetNotify(notify Notifier) {
	panic("Need to be implement!")
}

func (*LoggerMock) Fatalf(format string, v ...interface{}) {
	panic("Need to be implement!")
}

func (*LoggerMock) Errorf(format string, v ...interface{}) {
	panic("Need to be implement!")
}

func (*LoggerMock) Warnf(format string, v ...interface{}) {
	panic("Need to be implement!")
}

func (*LoggerMock) Infof(format string, v ...interface{}) {
	panic("Need to be implement!")
}

func (*LoggerMock) Debugf(format string, v ...interface{}) {
	panic("Need to be implement!")
}

func (*LoggerMock) Tracef(format string, v ...interface{}) {
	panic("Need to be implement!")
}

type IloggerMock struct {
	DebugfFunc func(format string, v ...interface{})

	ErrorfFunc func(format string, v ...interface{})

	FatalfFunc func(format string, v ...interface{})

	InfofFunc func(format string, v ...interface{})

	SetNotifyFunc func(notify Notifier)

	TracefFunc func(format string, v ...interface{})

	WarnfFunc func(format string, v ...interface{})
}

var _ Ilogger = &IloggerMock{}

func (*IloggerMock) Debugf(format string, v ...interface{}) {
	panic("Need to be implement!")
}

func (*IloggerMock) Errorf(format string, v ...interface{}) {
	panic("Need to be implement!")
}

func (*IloggerMock) Fatalf(format string, v ...interface{}) {
	panic("Need to be implement!")
}

func (*IloggerMock) Infof(format string, v ...interface{}) {
	panic("Need to be implement!")
}

func (*IloggerMock) SetNotify(notify Notifier) {
	panic("Need to be implement!")
}

func (*IloggerMock) Tracef(format string, v ...interface{}) {
	panic("Need to be implement!")
}

func (*IloggerMock) Warnf(format string, v ...interface{}) {
	panic("Need to be implement!")
}

type ILevelMock struct {
	AllFunc func() []Level
}

var _ ILevel = &ILevelMock{}

func (*ILevelMock) All() []Level {
	panic("Need to be implement!")
}

type NotifierMock struct {
	LevelsFunc func() []Level

	NotifyFunc func(msg string)
}

var _ Notifier = &NotifierMock{}

func (*NotifierMock) Levels() []Level {
	panic("Need to be implement!")
}

func (*NotifierMock) Notify(msg string) {
	panic("Need to be implement!")
}

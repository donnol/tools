package log

type ILevel interface{ All() []Level }

type Ilogger interface {
	Debugf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	SetNotify(notify Notifier)
	Tracef(format string, v ...interface{})
	Warnf(format string, v ...interface{})
}

type ILevelMock interface{ All() []Level }

type IloggerMock interface {
	Debugf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	SetNotify(notify Notifier)
	Tracef(format string, v ...interface{})
	Warnf(format string, v ...interface{})
}

type INotifierMock interface {
	Levels() []Level
	Notify(msg string)
}

type ILoggerMock interface {
	Debugf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	SetNotify(notify Notifier)
	Tracef(format string, v ...interface{})
	Warnf(format string, v ...interface{})
}

package log

type ILevel interface{ All() []Level }

type Ilogger interface {
	Debugf(format string, v ...any)
	Errorf(format string, v ...any)
	Fatalf(format string, v ...any)
	Infof(format string, v ...any)
	SetNotify(notify Notifier)
	Tracef(format string, v ...any)
	Warnf(format string, v ...any)
}

type ILevelMock interface{ All() []Level }

type IloggerMock interface {
	Debugf(format string, v ...any)
	Errorf(format string, v ...any)
	Fatalf(format string, v ...any)
	Infof(format string, v ...any)
	SetNotify(notify Notifier)
	Tracef(format string, v ...any)
	Warnf(format string, v ...any)
}

type INotifierMock interface {
	Levels() []Level
	Notify(msg string)
}

type ILoggerMock interface {
	Debugf(format string, v ...any)
	Errorf(format string, v ...any)
	Fatalf(format string, v ...any)
	Infof(format string, v ...any)
	SetNotify(notify Notifier)
	Tracef(format string, v ...any)
	Warnf(format string, v ...any)
}

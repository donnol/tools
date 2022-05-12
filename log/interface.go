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

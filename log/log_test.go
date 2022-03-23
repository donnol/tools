package log

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestInfo(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	for _, cas := range []struct {
		value any
	}{
		{errors.New("Test Message")},
	} {
		Fatalf("%+v\n", cas.value)
		Errorf("%+v\n", cas.value)
		Warnf("%+v\n", cas.value)
		Infof("%+v\n", cas.value)
		Debugf("%+v\n", cas.value)
		Tracef("%+v\n", cas.value)
	}
}

func TestLogger(t *testing.T) {
	logger := New(os.Stdout, "[Haha]", log.LstdFlags|log.Lshortfile)

	for _, cas := range []struct {
		value any
	}{
		{errors.New("Test Message")},
	} {
		logger.Fatalf("%+v\n", cas.value)
		logger.Errorf("%+v\n", cas.value)
		logger.Warnf("%+v\n", cas.value)
		logger.Infof("%+v\n", cas.value)
		logger.Debugf("%+v\n", cas.value)
		logger.Tracef("%+v\n", cas.value)
	}
}

type notify struct{}

// Levels 级别
func (n *notify) Levels() []Level {
	return []Level{FatalLevel}
}

// Notify 通知
func (n *notify) Notify(msg string) {
	fmt.Printf("%s", msg)
}

func TestNotifyLogger(t *testing.T) {
	logger := New(os.Stdout, "[NOTIFY]", log.LstdFlags|log.Lshortfile)
	n := &notify{}
	logger.SetNotify(n)

	for _, cas := range []struct {
		value any
	}{
		{errors.New("Notify Message")},
	} {
		logger.Fatalf("%+v\n", cas.value)
		logger.Errorf("%+v\n", cas.value)
		logger.Warnf("%+v\n", cas.value)
		logger.Infof("%+v\n", cas.value)
		logger.Debugf("%+v\n", cas.value)
		logger.Tracef("%+v\n", cas.value)
	}
}

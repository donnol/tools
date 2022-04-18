package retry

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	logger = zap.S()
)

type Doer = func(context.Context) (canRetry bool, err error)

// DoWithTimes 如果执行f失败了，重试tryTimes次
func DoWithTimes(ctx context.Context, tryTimes int, f Doer) error {
	if f == nil {
		return errors.New("f is nil")
	}

	var err error
	for i := 1; i <= tryTimes; i++ {
		if canRetry, err1 := f(ctx); err1 != nil {
			logger.Errorf("do No.%d failed: %+v", i, err1)
			err = err1

			if !canRetry {
				return err
			}

			time.Sleep(time.Second * (1 << i))
			continue
		}

		err = nil
		break
	}

	return err
}

// DoWithDeadline 如果执行f失败了，在t时间之前重试
func DoWithDeadline(ctx context.Context, d time.Time, f Doer) error {
	if f == nil {
		return errors.New("f is nil")
	}

	var err error
	var i = 1
	for {
		if canRetry, err1 := f(ctx); err1 != nil {
			logger.Errorf("do No.%d failed: %+v", i, err1)
			err = err1

			if !canRetry {
				return err
			}

			now := time.Now()
			if now.After(d) {
				break
			}

			i++
			time.Sleep(time.Second * (1 << i))
			continue
		}

		err = nil
		break
	}

	return err
}

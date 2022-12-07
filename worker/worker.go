package worker

import (
	"errors"
	"log"
	"os"
	"sync"
	"time"

	utillog "github.com/donnol/tools/log"
	"github.com/smallnest/chanx"
)

const (
	errCount = 100

	defaultCount = 50
)

var logger = utillog.New(os.Stdout, "[Worker]", log.LstdFlags|log.Lshortfile)

// 错误
var (
	ErrWorkerIsStop = errors.New("Worker is stop")
	ErrNilJobDo     = errors.New("Job do field is nil")
)

// DefaultWorker 默认Wroker
var DefaultWorker = New(defaultCount)

func init() {
	DefaultWorker.Start()
}

// Job 工作
type Job struct {
	do Do // 方法

	timeout time.Duration // 超时时间

	errorHandler ErrorHandler // 错误处理方法
}

// Do 执行
type Do func() error

// ErrorHandler 错误处理方法
type ErrorHandler func(error)

// MakeJob 新建工作
func MakeJob(do Do, timeout time.Duration, eh ErrorHandler) Job {
	return Job{
		do:           do,
		timeout:      timeout,
		errorHandler: eh,
	}
}

func (job Job) run() error {
	if job.do == nil {
		return ErrNilJobDo
	}

	if err := job.do(); err != nil {
		if job.errorHandler != nil {
			job.errorHandler(err)
		} else {
			return err
		}
	}

	return nil
}

// Worker 工人
type Worker struct {
	// 所有管道都要有make, read, write, close操作
	limitChan chan struct{}             // 并发控制管道
	stopChan  chan struct{}             // 停止管道
	jobChan   *chanx.UnboundedChan[Job] // 工作管道
	errChan   chan error                // 错误管道

	wg   *sync.WaitGroup
	stop bool // 是否调用了Stop方法
}

// New 新建
func New(n int) *Worker {
	if n <= 0 {
		n = defaultCount
	}
	return &Worker{
		limitChan: make(chan struct{}, n),
		stopChan:  make(chan struct{}),
		jobChan:   chanx.NewUnboundedChan[Job](n),
		errChan:   make(chan error, errCount),
		wg:        new(sync.WaitGroup),
	}
}

// Start 开始
func (w *Worker) Start() {
	go w.handleError()
	go w.start()

	logger.Debugf("Start.\n")
}

func (w *Worker) start() {
	for {
		select {
		case job, ok := <-w.jobChan.Out: // 有工作
			if !ok {
				continue
			}

			w.do(job)

		case <-w.stopChan:
			w.close()
			return
		}
	}
}

func (w *Worker) do(job Job) {
	// 占据一个坑
	w.limitChan <- struct{}{}

	// 开始工作
	go func(job Job) {
		defer func() {
			if r := recover(); r != nil {
				logger.Fatalf("job: %+v\n", r)
			}

			// 释放一个坑
			<-w.limitChan
			w.wg.Done()
		}()

		// 执行
		if job.timeout > 0 {
			w.doWithTimeout(job)
		} else {
			if err := job.run(); err != nil {
				w.errChan <- err
			}
		}
	}(job)
}

func (w *Worker) doWithTimeout(job Job) {
	var retChan = make(chan struct{})
	defer close(retChan)

	// 执行
	go func(retChan chan struct{}) {
		defer func() {
			if r := recover(); r != nil {
				logger.Fatalf("job exec: %+v\n", r)
			}
		}()

		// FIXME:
		// 这里不能直接这样调，如果job.run()执行的时间很长，将不会在超时后停止
		// 正确的做法应该是传入一个stopper管道，用户端代码需要适时检查该管道，判断是否需要停止
		// 参照'github.com/eapache/go-resiliency'的deadline包
		if err := job.run(); err != nil {
			w.errChan <- err
		}
		retChan <- struct{}{}
	}(retChan)

	// 超时
	timer := time.NewTimer(job.timeout)
	select {
	case <-retChan:
		return
	case t := <-timer.C:
		logger.Fatalf("job timeout: %+v\n", t)
		return
	}
}

func (w *Worker) handleError() {
	for err := range w.errChan {
		logger.Errorf("err is %v\n", err)
	}
}

// Stop 停止
func (w *Worker) Stop() {
	w.stop = true

	w.wait()

	w.stopChan <- struct{}{}

	logger.Debugf("Stop.\n")
}

func (w *Worker) close() {
	// close管道时，有可能panic
	defer func() {
		if r := recover(); r != nil {
			logger.Fatalf("close: %+v\n", r)
		}
	}()

	close(w.stopChan)
	close(w.errChan)
	close(w.jobChan.In)
	close(w.limitChan)
}

func (w *Worker) wait() {
	// 等待所有工作完成
	w.wg.Wait()
}

// Push 添加
func (w *Worker) Push(job Job) error {
	if w.stop {
		return ErrWorkerIsStop
	}
	if job.do == nil {
		return ErrNilJobDo
	}

	w.jobChan.In <- job
	w.wg.Add(1)

	return nil
}

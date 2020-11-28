package worker

type WorkerMock struct {
	PushFunc func(job Job) error

	StartFunc func()

	StopFunc func()
}

var _ IWorker = &WorkerMock{}

func (mockRecv *WorkerMock) Push(job Job) error {
	return mockRecv.PushFunc(job)
}

func (mockRecv *WorkerMock) Start() {
	mockRecv.StartFunc()
}

func (mockRecv *WorkerMock) Stop() {
	mockRecv.StopFunc()
}

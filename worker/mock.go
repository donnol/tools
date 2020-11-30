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

type WorkerMockMock struct {
	PushFunc func(job Job) error

	StartFunc func()

	StopFunc func()
}

var _ IWorkerMock = &WorkerMockMock{}

func (mockRecv *WorkerMockMock) Push(job Job) error {
	return mockRecv.PushFunc(job)
}

func (mockRecv *WorkerMockMock) Start() {
	mockRecv.StartFunc()
}

func (mockRecv *WorkerMockMock) Stop() {
	mockRecv.StopFunc()
}

package worker

type IWorkerMock struct {
	PushFunc func(job Job) error

	StartFunc func()

	StopFunc func()
}

var _ IWorker = &IWorkerMock{}

func (*IWorkerMock) Push(job Job) error {
	panic("Need to be implement!")
}

func (*IWorkerMock) Start() {
	panic("Need to be implement!")
}

func (*IWorkerMock) Stop() {
	panic("Need to be implement!")
}

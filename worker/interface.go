package worker

type IWorkerMock interface {
	Push(job Job) error
	Start()
	Stop()
}

type IWorker interface {
	Push(job Job) error
	Start()
	Stop()
}

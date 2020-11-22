package worker

type IWorker interface {
	Push(job Job) error
	Start()
	Stop()
}

package service

type PingSrv interface {
	Ping() string
}

func NewPingSrv() PingSrv {
	return &pingSrvImpl{}
}

type pingSrvImpl struct {
}

func (impl *pingSrvImpl) Ping() string {
	return "pong"
}

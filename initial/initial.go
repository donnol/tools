// Package initial init project and data
package initial

type Initial interface {
	Init() error
}

type (
	Config struct {
		Type string // project | data
		kind string // if Type is data, choose config | db
		Name string
	}
)

func New(
	conf *Config,
) Initial {
	switch conf.Type {
	case "project":
		return &project{
			name: conf.Name,
		}
	case "data":
		return &data{
			kind: conf.kind,
		}
	}
	return nil
}

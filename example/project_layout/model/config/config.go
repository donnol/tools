package config

import "fmt"

type Config struct {
	Name string `toml:"name"` // 项目名称
	Env  Env    `toml:"env"`  // 运行环境

	Server Server `toml:"server"`
	Mysql  Mysql  `toml:"mysql"`
	Redis  Redis  `toml:"redis"`
	Kafka  Kafka  `toml:"kafka"`
	Email  Email  `toml:"email"`
}

func (conf *Config) String() string {
	if conf == nil {
		return "<nil>"
	}
	return fmt.Sprintf("\n=== config ===\n%+v\n=== end    ===\n", *conf)
}

type Server struct {
	Port int
}

type Mysql struct {
}

type Redis struct {
}

type Kafka struct {
}

type Email struct {
	Host     string   `json:"host" toml:"host" flag:"host"`
	Port     int      `json:"port" toml:"port" flag:"port"`
	Username string   `json:"username" toml:"username" flag:"username"`
	Password Password `json:"password" toml:"password" flag:"password"`
}

type Password string

type Source string

const (
	SourceDefault Source = "DEFAULT" // 默认配置
	SourceFile    Source = "FILE"    // 配置文件
	SourceEnv     Source = "ENV"     // 环境变量
	SourceFlag    Source = "FLAG"    // 命令行参数
)

// Setter 如果返回值非nil，则使用返回值覆盖conf；为nil时则对传入conf覆盖指定字段值
type Setter func(conf *Config) *Config

type Option struct {
	Source  Source
	Setters []Setter
}

// 优先级从低到高：默认配置 -> 配置文件 -> 环境变量 -> 命令行参数
// 也就是优先级高的(只要有值就)覆盖优先级低的
func New(opts ...Option) *Config {
	conf := defaultConf

	// setter重新排序
	fileSetters := make([]Setter, 0, len(opts))
	envSetters := make([]Setter, 0, len(opts))
	flagSetters := make([]Setter, 0, len(opts))
	for _, opt := range opts {
		switch opt.Source {
		case SourceFile:
			fileSetters = append(fileSetters, opt.Setters...)
		case SourceEnv:
			envSetters = append(envSetters, opt.Setters...)
		case SourceFlag:
			flagSetters = append(flagSetters, opt.Setters...)
		}
	}

	// 设置
	setters := make([]Setter, 0, len(opts))
	setters = append(setters, fileSetters...)
	setters = append(setters, envSetters...)
	setters = append(setters, flagSetters...)
	for _, setter := range setters {
		c := setter(conf)
		if c != nil {
			conf = c
		}
	}

	return conf
}

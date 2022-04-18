package config

import "os"

type Env string

const (
	EnvDev  Env = "DEV"
	EnvTest Env = "TEST"
	EnvPre  Env = "PRE"
	EnvPrd  Env = "PRD"
)

func (env Env) Valid() bool {
	switch env {
	case EnvDev, EnvTest, EnvPre, EnvPrd:
		return true
	}
	return false
}

const (
	EnvOSKey = "PROJECT_LAYOUT_ENV"
)

func EnvFromOS() (Env, bool) {
	v, ok := os.LookupEnv(EnvOSKey)
	if !ok || v == "" {
		return Env(""), false
	}
	env := Env(v)
	if !env.Valid() {
		panic("'PROJECT_LAYOUT_ENV'环境变量的值非法")
	}
	return env, true
}

func (env Env) String() string {
	return string(env)
}

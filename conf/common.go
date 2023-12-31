package conf

import "log"

var EnvConfig Config

func LoadEnv(path string) {
	Env, err := LoadConfig(path, "app", "env")
	if err != nil {
		log.Print("failed to change ")
		return
	}
	EnvConfig = Env
}

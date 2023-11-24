package conf

import "log"

var EnvConfig Config

func LoadEnv() {
	Env, err := LoadConfig(".")
	if err != nil {
		log.Print("failed to change ")
		return
	}
	EnvConfig = Env
}

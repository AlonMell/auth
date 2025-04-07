package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

func MustLoad(config any, path string) {
	if path == "" {
		panic("config file path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file is not exist: " + path)
	}

	if err := cleanenv.ReadConfig(path, config); err != nil {
		panic("load config fail: " + err.Error())
	}
}

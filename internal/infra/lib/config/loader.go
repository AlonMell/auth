package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

func MustLoad(config any) {
	path := fetchConfigPath()
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

func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}

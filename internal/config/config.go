package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"providerHub/internal/constants"
	"time"
)

type Config struct {
	Env        string        `yaml:"env" env:"ENV" env-required:"true"`
	TokenTTL   time.Duration `yaml:"token_ttl" env-default:"1h"`
	DB         `yaml:"db" env-required:"true"`
	HTTPServer `yaml:"http_server" env-required:"true"`
}

type DB struct {
	Port     string `yaml:"port" env-default:"5432"`
	Version  string `yaml:"version" env-default:"16"`
	Username string `yaml:"user" env-default:"postgres"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == constants.EmptyString {
		panic("config file path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file is not exist: " + path)
	}

	cfg := new(Config)
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		panic("load config fail: " + err.Error())
	}

	return cfg
}

func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config", constants.EmptyString, "path to config file")
	flag.Parse()

	if path == constants.EmptyString {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}

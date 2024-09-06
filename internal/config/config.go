package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-required:"true"`
	DB         `yaml:"db" env-required:"true"`
	HTTPServer `yaml:"http_server" env-required:"true"`
	JWT        `yaml:"jwt" env-required:"true"`
}

type DB struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-default:"postgres"`
	Password string `yaml:"password" env-default:"postgres"`
	Database string `yaml:"database" env-default:"postgres"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type JWT struct {
	AccessTTL  time.Duration `yaml:"access_ttl" env-default:"15m"`
	RefreshTTL time.Duration `yaml:"refresh_ttl" env-default:"168h"`
	Secret     string        `yaml:"secret" env-required:"true"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
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

	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}

package config

import "flag"

type Paths struct {
	Migrations string
	Config     string
}

func MustLoadFlags() Paths {
	var paths Paths
	flag.StringVar(&paths.Migrations, "migrations", "", "Path to the migrations directory")
	flag.StringVar(&paths.Config, "config", "", "Path to the configuration file")
	flag.Parse()

	if paths.Migrations == "" {
		panic("migrations path is required")
	}
	if paths.Config == "" {
		panic("config path is required")
	}

	return paths
}

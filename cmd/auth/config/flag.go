package config

import "flag"

type Paths struct {
	Config string
}

func MustLoadFlags() Paths {
	var paths Paths
	flag.StringVar(&paths.Config, "config", "", "Path to the configuration file")
	flag.Parse()

	if paths.Config == "" {
		panic("config path is required")
	}

	return paths
}

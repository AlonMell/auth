package main

import (
	"fmt"
	"providerHub/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Printf("%#v\n", *cfg)
}

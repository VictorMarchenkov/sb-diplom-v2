package main

import (
	"flag"
	"fmt"
	"os"

	"sb-diplom-v2/internal/app"
	"sb-diplom-v2/pkg/cfg"
)

func main() {
	var configFilePath string
	flag.StringVar(&configFilePath, "config", "config/config.json", "cofig. file path")
	flag.Parse()

	if configFilePath == "" {
		fmt.Println("config. file path is empty")
		os.Exit(1)
	}

	cfg, err := cfg.Init(configFilePath)
	if err != nil {
		fmt.Println("config. file path is empty")
		os.Exit(1)
	}

	app.Run(cfg)
}

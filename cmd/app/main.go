package main

import (
	"flag"
	"log"
	"os"

	"sb-diplom-v2/internal/app"
	"sb-diplom-v2/pkg/cfgPath"
)

// main -.
func main() {
	var configFilePath string
	flag.StringVar(&configFilePath, "config", "config/config.json", "config file path")
	flag.Parse()

	defer func() {
		if err := recover(); err != nil {
			log.Println("Panic recovered after error: %v", err)
		}
	}()
	if configFilePath == "" {
		log.Println("config. file path is empty")
		os.Exit(1)
	}

	appCfg, err := cfgPath.Init(configFilePath)
	if err != nil {
		log.Println("config file path is not valid", err)
		os.Exit(1)
	}
	app.Run(appCfg)
}

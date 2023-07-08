package main

import (
	"flag"
	"log"
	"os"

	"sb-diplom-v2/internal/app"
	"sb-diplom-v2/pkg/cfgPath"
	"sb-diplom-v2/pkg/logger"
)

// main -.
func main() {
	var configFilePath string
	flag.StringVar(&configFilePath, "config", "config/config.json", "config file path")
	flag.Parse()

	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic recovered after error: %v\n", err)
		}
	}()

	l := logger.New("main")
	l.Info("starting skillbox diploma")
	defer l.Info("exit")

	if configFilePath == "" {
		l.Error("config. file path is empty")
		os.Exit(1)
	}

	appCfg, err := cfgPath.Init(configFilePath)
	if err != nil {
		l.Error("config file path is not valid", err)
		os.Exit(1)
	}
	app.Run(appCfg)
}

package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cfg "sb-diplom-v2/internal"
	"sb-diplom-v2/internal/logger"
	"sb-diplom-v2/pkg"
	cfg2 "sb-diplom-v2/pkg/cfg"

	"github.com/go-co-op/gocron"
)

// Run presents the server logic
func Run(cfgRoot *cfg2.Root) {
	var cfg_ cfg.Config
	var resultT cfg.StatusResult
	var service = "skillbox diploma"

	defer func() {
		if err := recover(); err != nil {
			logger.Logger.Warnf("Panic recovered after error: %v", err)
		}
	}()
	start := fmt.Sprintf(":%d", cfgRoot.HTTPServer.Port)
	data, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	json.Unmarshal(data, &cfg_)

	mux := http.NewServeMux()
	s := gocron.NewScheduler(time.UTC)
	fs := http.FileServer(http.Dir("generator"))
	http.Handle("/", fs)
	info := http.FileServer(http.Dir("."))
	mux.Handle("/info/", info)

	s.Every("10s").Do(func() {
		logger.Logger.Println("files rereading")
		resultT.HandlerFiles(cfg_)
	})

	s.StartAsync()

	mux.Handle("/presentation/", info)
	mux.HandleFunc("/api", resultT.HandlerHTTP)

	panicWrapped := pkg.PanicMiddleware(mux)
	siteHandler := pkg.AccessLogMiddleware(panicWrapped)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s", start),
		Handler: siteHandler,
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Fatalf("listen error %v\n", err)
		}
	}()
	logger.Logger.Printf("%s starting on %d", service, cfgRoot.HTTPServer.Port)
	<-stop

	logger.Logger.Printf("%s shutting down ...", service)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Logger.Errorf("Server Shutdown Failed:%+v", err)
	}
}

package app

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cfg "sb-diplom-v2/internal"
	"sb-diplom-v2/pkg"
	cfg2 "sb-diplom-v2/pkg/cfgPath"
)

// Run presents the server logic
func Run(cfgRoot *cfg2.Root) {

	var resultT cfg.StatusResult
	var service = "skillbox diploma"

	start := fmt.Sprintf(":%d", cfgRoot.HTTPServer.Port)

	mux := http.NewServeMux()
	s := gocron.NewScheduler(time.UTC)
	fs := http.FileServer(http.Dir("generator"))
	http.Handle("/", fs)
	info := http.FileServer(http.Dir("."))
	mux.Handle("/info/", info)

	s.Every("10s").Do(func() {
		fmt.Println("files rereading")
		resultT.HandlerFiles(cfgRoot)
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
			fmt.Printf("listen error %v\n", err)
			os.Exit(1)
		}
	}()
	fmt.Printf("%s starting on %d", service, cfgRoot.HTTPServer.Port)
	<-stop

	fmt.Printf("%s shutting down ...", service)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Errorf("Server Shutdown Failed:%+v", err)
	}
}

package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cfg "sb-diplom-v2/internal"
	"sb-diplom-v2/pkg"
	cfg2 "sb-diplom-v2/pkg/cfgPath"
	"sb-diplom-v2/pkg/logger"

	"github.com/go-co-op/gocron"
)

// Run runs HTTP-server.
func Run(cfgRoot *cfg2.Root) {
	l := logger.New("http-server")
	l.Info("starting on %v", cfgRoot.HTTPServer.HostPort())
	defer l.Info("exit")

	s := gocron.NewScheduler(time.UTC)
	var resultT cfg.StatusResult
	if _, err := s.Every("10s").Do(func() {
		l.Info("files rereading")
		resultT.HandlerFiles(cfgRoot)
	}); err != nil {
		l.Error("scheduling error", err)
		return
	}
	s.StartAsync()

	mux := http.NewServeMux()
	http.Handle("/", http.FileServer(http.Dir("generator")))
	mux.Handle("/info/", http.FileServer(http.Dir(".")))
	mux.Handle("/presentation/", http.FileServer(http.Dir(".")))
	mux.HandleFunc("/api", resultT.HandlerHTTP)

	panicWrapped := pkg.PanicMiddleware(mux)
	siteHandler := pkg.AccessLogMiddleware(panicWrapped)

	srv := &http.Server{
		Addr:    cfgRoot.HTTPServer.HostPort(),
		Handler: siteHandler,
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Error("listen: %v", err)
			os.Exit(1)
		}
	}()

	<-stop

	l.Info("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		l.Error("shutdown: %v", err)
	}
}

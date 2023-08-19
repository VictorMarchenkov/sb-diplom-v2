package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"sb-diplom-v2/internal"
	"sb-diplom-v2/pkg"
	"sb-diplom-v2/pkg/configs"
	"sb-diplom-v2/pkg/logger"
)

// Run runs HTTP-server.
func Run(cfgRoot *configs.Root) {
	l := logger.New("http-server")
	l.Info("starting on %v", cfgRoot.HTTPServer.HostPort())
	defer l.Info("exit")

	// async file updater
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	resultT := internal.NewStatusResult(cfgRoot)

	go func(wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()

		ticker := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-ticker.C:
				l.Info("files reading")
				resultT.HandlerFiles(cfgRoot)
			case <-ctx.Done():
				l.Info("updater exited")
				return
			}
		}
	}(&wg)

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

	cancel()
	wg.Wait()

	l.Info("shutting down")

	if err := srv.Shutdown(ctx); err != nil {
		l.Error("shutdown: %v", err)
	}
}

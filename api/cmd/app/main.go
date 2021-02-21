package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/julienschmidt/httprouter"

	// Logger
	_loggerUcase "github.com/wascript3r/cryptopay/pkg/logger/usecase"

	// Query
	_queryHandler "github.com/wascript3r/scraper/api/pkg/query/delivery/http"
	_queryRepo "github.com/wascript3r/scraper/api/pkg/query/repository"
	_queryUcase "github.com/wascript3r/scraper/api/pkg/query/usecase"
)

const (
	AppLoggerPrefix = "[APP]"
)

var (
	WorkDir string
	Cfg     *Config
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var err error

	// Get working directory
	WorkDir, err = os.Getwd()
	if err != nil {
		fatalError(err)
	}

	// Parse config file
	cfgPath, err := getConfigPath()
	if err != nil {
		fatalError(err)
	}

	Cfg, err = parseConfig(filepath.Join(WorkDir, cfgPath))
	if err != nil {
		fatalError(err)
	}
}

func fatalError(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func main() {
	// Logging
	logFlags := 0
	if Cfg.Log.ShowTimestamp {
		logFlags = log.Ltime
	}
	logger := _loggerUcase.New(
		AppLoggerPrefix,
		log.New(os.Stdout, "", logFlags),
	)

	// Startup message
	logger.Info("... Starting app ...")

	// Query
	queryRepo := _queryRepo.NewMySQLRepo()
	queryUcase := _queryUcase.New(
		queryRepo,
		Cfg.Database.MySQL.QueryTimeout.Duration,
	)

	// Graceful shutdown
	stopSig := make(chan os.Signal, 1)
	signal.Notify(stopSig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	_, cancel := context.WithCancel(context.Background())

	// HTTP server
	httpRouter := httprouter.New()

	if Cfg.HTTP.EnablePprof {
		// pprof
		httpRouter.Handler(http.MethodGet, "/debug/pprof/*item", http.DefaultServeMux)
	}

	_queryHandler.NewHTTPHandler(httpRouter, queryUcase)

	httpServer := &http.Server{
		Addr:    ":" + Cfg.HTTP.Port,
		Handler: httpRouter,
	}

	// Graceful shutdown
	gracefulShutdown := func() {
		cancel()

		if err := httpServer.Shutdown(context.Background()); err != nil {
			logger.Error("Cannot shutdown HTTP server: %s", err)
		}

		logger.Info("... Exited ...")
		os.Exit(0)
	}

	go func() {
		<-stopSig
		logger.Info("... Received stop signal ...")
		gracefulShutdown()
	}()

	if err := httpServer.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			fmt.Println(err)
			gracefulShutdown()
		}
	}

	var c chan struct{}
	<-c
}

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

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"

	// Logger
	_loggerUcase "github.com/wascript3r/cryptopay/pkg/logger/usecase"

	// Query
	_queryHandler "github.com/wascript3r/scraper/api/pkg/query/delivery/http"
	_queryRepo "github.com/wascript3r/scraper/api/pkg/query/repository"
	_queryUcase "github.com/wascript3r/scraper/api/pkg/query/usecase"

	// Location
	_locationRepo "github.com/wascript3r/scraper/api/pkg/location/repository"

	// Photo
	_photoRepo "github.com/wascript3r/scraper/api/pkg/photo/repository"

	// Condition
	_conditionRepo "github.com/wascript3r/scraper/api/pkg/condition/repository"

	// Seller
	_sellerRepo "github.com/wascript3r/scraper/api/pkg/seller/repository"

	// Listing
	_listingHandler "github.com/wascript3r/scraper/api/pkg/listing/delivery/http"
	_listingRepo "github.com/wascript3r/scraper/api/pkg/listing/repository"
	_listingUcase "github.com/wascript3r/scraper/api/pkg/listing/usecase"
	_listingValidator "github.com/wascript3r/scraper/api/pkg/listing/validator"

	// Auth
	_authMid "github.com/wascript3r/scraper/api/pkg/auth/delivery/http/middleware"
	_authUcase "github.com/wascript3r/scraper/api/pkg/auth/usecase"
)

const (
	// Database
	DatabaseDriver = "mysql"

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

	// Database connection
	dbConn, err := openDatabase(DatabaseDriver, Cfg.Database.MySQL.DSN)
	if err != nil {
		fatalError(err)
	}

	// Query
	queryRepo := _queryRepo.NewMySQLRepo(dbConn)
	queryUcase := _queryUcase.New(
		queryRepo,
		Cfg.Database.MySQL.QueryTimeout.Duration,
	)

	// Location
	locationRepo := _locationRepo.NewMySQLRepo(dbConn)

	// Photo
	photoRepo := _photoRepo.NewMySQLRepo(dbConn)

	// Photo
	conditionRepo := _conditionRepo.NewMySQLRepo(dbConn)

	// Seller
	sellerRepo := _sellerRepo.NewMySQLRepo(dbConn)

	// Listing
	listingRepo := _listingRepo.NewMySQLRepo(dbConn)
	listingValidator := _listingValidator.New()
	listingUcase := _listingUcase.New(
		listingRepo,
		locationRepo,
		photoRepo,
		queryRepo,
		conditionRepo,
		sellerRepo,
		Cfg.Database.MySQL.QueryTimeout.Duration,

		listingValidator,
	)

	// Graceful shutdown
	stopSig := make(chan os.Signal, 1)
	signal.Notify(stopSig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	_, cancel := context.WithCancel(context.Background())

	// HTTP server
	httpRouter := httprouter.New()
	httpRouter.MethodNotAllowed = MethodNotAllowedHnd
	httpRouter.NotFound = NotFoundHnd

	if Cfg.HTTP.EnablePprof {
		// pprof
		httpRouter.Handler(http.MethodGet, "/debug/pprof/*item", http.DefaultServeMux)
	}

	_queryHandler.NewHTTPHandler(httpRouter, queryUcase)
	_listingHandler.NewHTTPHandler(httpRouter, listingUcase)

	authUcase := _authUcase.New(Cfg.HTTP.Auth.BearerToken)
	authMid := _authMid.NewHTTPMiddleware(httpRouter, authUcase)

	httpServer := &http.Server{
		Addr:    ":" + Cfg.HTTP.Port,
		Handler: authMid,
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

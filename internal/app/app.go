package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"syspodhot/config"
	"syspodhot/internal/app/service/climanager"
	"syspodhot/internal/app/service/server"
	"syspodhot/internal/pkg/db"
	"syspodhot/internal/pkg/logger"
	"syspodhot/internal/pkg/validator"
)

// Service describes app service started in parallel with other similar services.
type Service interface {
	// Started returns chan that signals that service's readiness.
	Started() <-chan struct{}
	// Exited returns chan that signals that service's exited.
	Exited() <-chan struct{}
	// Run starts service. Waits to context done to exit.
	Run(ctx context.Context, errChan chan<- error)
}

// App represents full application with a list of service.
type App struct {
	services []Service
}

// NewServe returns a new instance of App with HTTP-server.
func NewServe() (*App, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("create config: %w", err)
	}
	// setup logger
	logger.InitSlog(
		logger.WithLogLevel(cfg.Logging.LogLevel),
		logger.WithFormat(cfg.Logging.LogFormat),
		logger.WithLocale(cfg.Locale),
	)
	// print out config params in debug log mode
	slog.Debug("current config", "config", cfg)

	// init validator
	valid, err := validator.NewRuTagValidator()
	if err != nil {
		return nil, fmt.Errorf("create validator: %w", err)
	}
	// init sources
	dbInst, err := db.New(cfg.DB.DSN, cfg.DB.Name)
	if err != nil {
		return nil, fmt.Errorf("create db conn: %w", err)
	}
	// init services
	serverService, err := server.New(cfg, dbInst, valid)
	if err != nil {
		return nil, fmt.Errorf("init server service: %w", err)
	}

	slog.Info("init app", "status", "OK")
	return &App{
		[]Service{serverService},
	}, nil
}

// NewCLI returns a new instance of App with CLI-manager.
func NewCLI() (*App, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("create config: %w", err)
	}
	// setup logger (text format and error level)
	logger.InitSlog(logger.WithErrorLevel(), logger.WithLocale(cfg.Locale))

	// connect to DB
	dbInst, err := db.New(cfg.DB.DSN, cfg.DB.Name)
	if err != nil {
		return nil, fmt.Errorf("create db conn: %w", err)
	}

	// init cmd manager service
	manager, err := climanager.New(cfg, dbInst)
	if err != nil {
		return nil, fmt.Errorf("create cli manager service: %w", err)
	}

	return &App{
		services: []Service{manager},
	}, nil
}

// Run starts app. This function is blocking.
// It waits for os signal to gracefully shutdown.
func (a *App) Run() (appErr error) {
	// ctx for app services
	appContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// handle shutdown process signals
	quitSig := make(chan os.Signal, 1)
	signal.Notify(quitSig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	// wait groups to sync starting/exiting all services
	var wgStarted, wgExited sync.WaitGroup

	// start services
	serviceErr := make(chan error, 1)
	for _, service := range a.services {
		wgStarted.Go(func() {
			<-service.Started()
		})
		wgExited.Go(func() {
			<-service.Exited()
		})
		go service.Run(appContext, serviceErr)
	}

	wgStarted.Wait()
	select {
	case handledSignal := <-quitSig:
		cancel()
		slog.Warn("got os signal; shutdown app...", "signal", handledSignal.String())
	case err := <-serviceErr:
		cancel()
		appErr = fmt.Errorf("service: %w", err)
		slog.Warn("one of services fell down; shutdown app...")
	case <-appContext.Done():
		appErr = appContext.Err()
		slog.Warn("context canceled; shutdown app...")
	}
	wgExited.Wait()
	slog.Info("all services exited")
	return appErr
}

// Package server provides HTTP-server service.
package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/template/html/v2"

	"syspodhot/config"
	"syspodhot/internal/app/service/server/middleware"
	"syspodhot/internal/pkg/db"
	"syspodhot/internal/pkg/healthcheck"
	"syspodhot/internal/pkg/tls"
	"syspodhot/internal/pkg/validator"
)

// Server represents a HTTP-server.
type Server struct {
	cfg      *config.Config
	fiberApp *fiber.App

	started chan struct{}
	exited  chan struct{}
}

// New returns a new instance of Server.
func New(cfg *config.Config, dbStorage *db.DB, valid validator.Validator) (*Server, error) {
	// fiber init
	server := &Server{
		cfg: cfg,
		fiberApp: fiber.New(fiber.Config{
			Views:         html.New("./web/template", ".html"),
			ServerHeader:  "Syspodhot",
			StrictRouting: false,
			// ErrorHandler:  errhandler.CustomErrorHandler,
		}),
		started: make(chan struct{}),
		exited:  make(chan struct{}),
	}

	// set up base middlewares
	server.fiberApp.Use(middleware.Logger(
		middleware.LoggerWithLevel(cfg.Logging.LogLevel),
		middleware.LoggerWithFormat(cfg.Logging.LogFormat),
		middleware.LoggerWithLocale(cfg.Locale)))
	server.fiberApp.Use(middleware.Recover())
	// connect healthcheck
	hcInst := healthcheck.New([]healthcheck.Checking{dbStorage})
	server.fiberApp.Get("/health", adaptor.HTTPHandlerFunc(hcInst.HealthHanlder))
	// register all endpoints
	// TODO: uncomment
	// server.registerEndpointsV1(cfg, dbStorage, valid)

	return server, nil
}

// Run starts server and waits for
// context is done for gracefully shutdown server.
// This method is blocking.
func (s *Server) Run(ctx context.Context, errChan chan<- error) {
	slog.Info("start http-server service...")

	// start server
	go s.listenAndServe(errChan)
	defer close(s.exited)
	// wait for context
	<-ctx.Done()
	if err := s.fiberApp.ShutdownWithTimeout(s.cfg.Server.ShutdownTimeout); err != nil {
		errChan <- fmt.Errorf("server: shutdown: %w", err)
	}
}

// Started returns chan that signals the service's readiness.
func (s *Server) Started() <-chan struct{} {
	return s.started
}

// Exited returns chan that signals that service's exited.
func (s *Server) Exited() <-chan struct{} {
	return s.exited
}

// listenAndServe starts HTTP-server and serve incomming connections.
func (s *Server) listenAndServe(errChan chan<- error) {
	addr := "0.0.0.0:" + s.cfg.Server.Port

	// listen port
	listener, err := net.Listen("tcp", addr) //nolint:noctx // ListenConfig is not necessary
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		errChan <- fmt.Errorf("start server listener: %w", err)
		return
	}
	close(s.started)

	useTLS := s.cfg.SSL.Cert != "" && s.cfg.SSL.Key != ""
	slog.Info("listener started successfully",
		"addr", addr,
		"tls", useTLS)

	// with SSL
	if useTLS {
		listener, err = tls.WrapWithTLS(listener, s.cfg.SSL.Cert, s.cfg.SSL.Key)
		if err != nil {
			errChan <- fmt.Errorf("server: wrap listener with tls: %w", err)
			return
		}
	}
	if err := s.fiberApp.Listener(listener); err != nil {
		errChan <- fmt.Errorf("server: listen: %w", err)
	}
}

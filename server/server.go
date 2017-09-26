// Package server is a simple wrapper around std http.Server
// that applies graceful shutdown.
//
// Typical usage:
//  srv := server.New(addr, handler, logger)
//  go srv.Start()
//  srv.Wait()
//  srv.Shutdown()
//
// If you want to use custom http.Server:
//  httpServer := &http.Server{...}
//  srv := server.Wrap(srv, logger)
package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

// Server is a http server with graceful shutdown.
type Server struct {
	origin      *http.Server
	stopSignals chan os.Signal
	lg          logrus.StdLogger // TODO: better interface
}

// New returns a new Server.
func New(addr string, handler http.Handler, logger logrus.StdLogger) *Server {
	stopSignals := make(chan os.Signal, 1)
	signal.Notify(stopSignals, syscall.SIGINT, syscall.SIGTERM)

	return &Server{
		origin:      &http.Server{Addr: addr, Handler: handler},
		stopSignals: stopSignals,
		lg:          logger,
	}
}

// Wrap returns a new Server that wraps http.Server.
func Wrap(srv *http.Server, logger logrus.StdLogger) *Server {
	stopSignals := make(chan os.Signal, 1)
	signal.Notify(stopSignals, syscall.SIGINT, syscall.SIGTERM)

	return &Server{
		origin:      srv,
		stopSignals: stopSignals,
		lg:          logger,
	}
}

// Start makes server listen and serve.
// It blocks until server is stopped.
func (s *Server) Start() {
	s.lg.Printf("Start listening @ %s", s.origin.Addr)
	err := s.origin.ListenAndServe()
	if err != http.ErrServerClosed {
		s.lg.Fatal(err)
	}
	s.lg.Println("Server closed.")
}

// Wait blocks until SIGINT or SIGTERM is received.
// Stop() can be called to unblock manually.
func (s *Server) Wait() {
	<-s.stopSignals
}

// Stop unblocks waiting server.
func (s *Server) Stop() {
	s.stopSignals <- syscall.SIGTERM
}

// Shutdown tries to gracefully shutdown server.
func (s *Server) Shutdown() {
	s.lg.Println("Shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
	defer cancel()

	if err := s.origin.Shutdown(ctx); err != nil {
		s.lg.Printf("Server graceful shutdown failed: %s\n", err)
	} else {
		s.lg.Println("Server gracefully shut down.")
	}
}

const (
	gracefulTimeout = time.Second * 10 // TODO: make configurable
)

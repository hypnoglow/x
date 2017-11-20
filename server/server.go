// Package server is a simple wrapper around std http.Server
// that applies graceful shutdown.
//
// Typical usage:
//  srv := server.New(addr, handler)
//  go srv.Start()
//  srv.Wait()
//  srv.Shutdown()
//
// The example above stops the server only when a SIGINT is sent to the app.
// If you want to manually stop the server, just call Stop() when you need:
//  go func() {
//      time.Sleep(time.Second * 5)
//      srv.Stop()
//  }()
//
// If you want to use custom http.Server:
//  httpServer := &http.Server{...}
//  srv := server.Wrap(srv)
package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

// Server is a http server with graceful shutdown.
type Server struct {
	origin *http.Server
	log    io.Writer

	stopSignals chan os.Signal
	onceCloser  sync.Once
}

// Option for server.
type Option func(*Server)

// Log returns an option that sets server logger.
func Log(log io.Writer) Option {
	return func(s *Server) {
		s.log = log
	}
}

// New returns a new Server.
func New(addr string, handler http.Handler, opts ...Option) *Server {
	stopSignals := make(chan os.Signal, 1)
	signal.Notify(stopSignals, os.Interrupt)

	s := &Server{
		origin:      &http.Server{Addr: addr, Handler: handler},
		stopSignals: stopSignals,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Wrap returns a new Server that wraps http.Server.
func Wrap(srv *http.Server, opts ...Option) *Server {
	stopSignals := make(chan os.Signal, 1)
	signal.Notify(stopSignals, os.Interrupt)

	s := &Server{
		origin:      srv,
		stopSignals: stopSignals,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Start makes server listen and serve.
// It blocks until server is stopped.
func (s *Server) Start() {
	s.logMessage("Start listening @ %s", s.origin.Addr)
	err := s.origin.ListenAndServe()
	if err != http.ErrServerClosed {
		s.logMessage(err.Error())
		s.Stop() // just to ensure everything is cleaned.
		return
	}

	s.logMessage("Server closed.")
}

// Wait blocks until SIGINT or SIGTERM is received.
// Stop() can be called to unblock manually.
func (s *Server) Wait() {
	<-s.stopSignals
}

// Stop unblocks waiting server, closing its signal channel.
func (s *Server) Stop() {
	s.onceCloser.Do(func() {
		signal.Stop(s.stopSignals)
		close(s.stopSignals)
	})
}

// Shutdown tries to gracefully shutdown server.
func (s *Server) Shutdown() {
	s.logMessage("Shutdown server...")
	s.Stop() // in case shutdown is triggered by a signal from os.

	ctx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
	defer cancel()

	if err := s.origin.Shutdown(ctx); err != nil {
		s.logMessage("Server graceful shutdown failed: %s\n", err)
	} else {
		s.logMessage("Server gracefully shut down.")
	}
}

func (s *Server) logMessage(format string, args ...interface{}) {
	if s.log == nil {
		return
	}

	fmt.Fprintf(s.log, format, args...)
}

const (
	gracefulTimeout = time.Second * 10 // TODO: make configurable
)

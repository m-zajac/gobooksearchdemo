package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"
)

// Server handles app's http requests.
type Server struct {
	port    uint
	handler http.Handler
}

// NewServer creates new Server instance.
func NewServer(port uint, handler http.Handler) *Server {
	return &Server{
		port:    port,
		handler: handler,
	}
}

// Run runs the server. Waits SIGINT is received, then gracefully shutdowns.
// Blocks until shutdown is complete.
func (s *Server) Run() {
	srv := http.Server{
		Addr: fmt.Sprintf(":%d", s.port),

		// For timeouts explanation see: https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      300 * time.Second,
		IdleTimeout:       10 * time.Second,

		Handler: s.handler,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		logrus.Infof("starting http server, listening on :%d", s.port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("server returned error: %v", err)
		}
	}()

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		logrus.Errorf("http server shutdown returned error: %v", err)
	}
	logrus.Info("http server shut down")
}

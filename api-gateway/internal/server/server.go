package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    port,
			Handler: handler,
		},
	}
}

func (s *Server) Runserver() error {
	log.Printf("API Gateway is starting at %v", s.httpServer.Addr)

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v", s.httpServer.Addr, err)
		}
	}()

	<-stop

	log.Println("Shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Forced to shut down: %v", err)
	}
	log.Println("Server gracefully stopped")
	return nil
}

package server

import (
	"fmt"
	"forum/internal/config"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) ServerRun(cfg *config.Config, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:         ":" + cfg.Handler.Addr,
		Handler:      handler,
		ReadTimeout:  time.Duration(cfg.Handler.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Handler.WriteTimeout) * time.Second,
	}
	fmt.Printf("Server starting http://localhost:%s\n", cfg.Handler.Addr)
	return s.httpServer.ListenAndServe()
}

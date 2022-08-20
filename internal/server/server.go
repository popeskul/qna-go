// Package server implements HTTP server.
// You can use it to implement a web server.
// Example:
// 		func main() {
// 			r := mux.NewRouter()
// 			r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 				w.Write([]byte("Hello World!"))
// 			}).Methods("GET")
// 			srv := server.NewServer(server.Config{
// 				Addr: ":8080",
// 				Handler: r,
// 			})
// 			srv.Run()
// 		}

package server

import (
	"context"
	"net/http"
)

// Server is an HTTP server.
type Server struct {
	httpServer *http.Server
}

// NewServer creates a new server.
// Need to pass config to create a server.
// Note that the server is not started.
// You must call Run method to start the server.
func NewServer(cfg *http.Server) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           cfg.Addr,
			Handler:        cfg.Handler,
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.WriteTimeout,
			MaxHeaderBytes: cfg.MaxHeaderBytes,
		},
	}
}

// Run starts the server.
func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server without interrupting any active connections.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

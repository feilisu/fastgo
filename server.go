package fastgo

import (
	"log"
	"net/http"
	"time"
)

type Server struct {
	addr         string
	port         string
	readTimeout  time.Duration
	writeTimeout time.Duration
	errorLog     *log.Logger
	server       *http.Server
}

func NewServer() *Server {
	return &Server{
		addr:         "0.0.0.0",
		port:         "88",
		readTimeout:  5 * time.Second,
		writeTimeout: 5 * time.Second,
		errorLog:     GetLogger(),
	}
}

func (s *Server) Run(r *http.ServeMux) error {
	s.server = &http.Server{
		Addr:         s.addr + ":" + s.port,
		Handler:      r,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		ErrorLog:     s.errorLog,
	}
	return s.server.ListenAndServe()
}

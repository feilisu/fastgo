package fastgo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/pprof"
	"time"
)

type Server struct {
	Name         string        `json:"name"`
	Addr         string        `json:"addr"`
	Port         string        `json:"port"`
	ReadTimeout  time.Duration `json:"readTimeout"`
	WriteTimeout time.Duration `json:"writeTimeout"`
	ErrorLog     *log.Logger
	server       *http.Server
	Middlewares  []Middleware
	Router       Router
}

func DefaultServer() *Server {
	return &Server{
		Name:         "default",
		Addr:         "0.0.0.0",
		Port:         "8090",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
}

func (s *Server) Run(r *Router) error {

	elog := s.ErrorLog
	if elog == nil {
		elog = serverErrorLogger()
	}

	s.server = &http.Server{
		Addr:         s.Addr + ":" + s.Port,
		ReadTimeout:  s.ReadTimeout,
		WriteTimeout: s.WriteTimeout,
		ErrorLog:     elog,
		//BaseContext:  ServerBaseContext,
		ConnContext: ServerConnContext,
	}
	s.registerHandle(r)
	return s.server.ListenAndServe()
}

// registerHandle
func (s *Server) registerHandle(r *Router) {

	//serveMux := http.NewServeMux()
	serveMux := mux.NewRouter()
	s.initDebugRouter(serveMux)

	if r.entryMap != nil {
		for host, entryList := range r.entryMap {
			for _, entry := range entryList {
				log.Printf("%s %s %s", entry.methods, host, entry.path)

				if entry.host == DefaultHost {
					serveMux.Host(entry.host).Path(entry.path).Methods(entry.methods...).HandlerFunc(s.getHandler(entry))
				} else {
					serveMux.Path(entry.path).Methods(entry.methods...).HandlerFunc(s.getHandler(entry))
				}
			}
		}
	}

	s.server.Handler = serveMux
}

// initDebugRouter
func (s *Server) initDebugRouter(serveMux *mux.Router) {
	serveMux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	serveMux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	serveMux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	serveMux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
}

// middlewareHandler 中间件处理
func (s *Server) middlewareHandle(ctx *Context, ms []Middleware) error {
	for _, middleware := range ms {
		err := middleware.Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// getHandler 路由处理器
func (s *Server) getHandler(entry Entry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer LoggerSync()
		defer func() {
			err := recover()
			if err != nil {

				res := make(map[string]string)
				res["httpCode"] = "500"
				res["message"] = fmt.Sprintf("%v", err)
				marshal, _ := json.Marshal(res)
				w.WriteHeader(HttpStatus500)
				w.Header().Set("Content-Type", "application/json;charset=utf8")
				_, _ = w.Write(marshal)

			}
		}()

		ctx := context.WithValue(r.Context(), 1, 1)
		response := buildResponse(w)
		request, err := buildRequest(r)
		if err != nil {
			_ = response.Error(err)
			return
		}
		sctx := &Context{
			Context:  ctx,
			Request:  request,
			Response: response,
		}

		err = s.middlewareHandle(sctx, s.Middlewares)
		if err != nil {
			_ = sctx.Response.Json(err)
			return
		}

		err = s.middlewareHandle(sctx, entry.middlewares)
		if err != nil {
			_ = sctx.Response.Json(err)
			return
		}

		err = entry.handle.Handle(sctx)
		if err != nil {
			_ = sctx.Response.Json(err)
			return
		}
		return
	}
}

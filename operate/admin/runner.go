package admin

import (
	"net"
	"net/http"
	"net/http/pprof"
	"os"

	"golang.org/x/net/trace"

	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/http_server"
)

func Runner(port string, options ...func(*Server) error) (ifrit.Runner, error) {
	mux := defaultDebugEndpoints()

	return &Server{
		port:    port,
		handler: mux,
	}, nil
}

type Server struct {
	port    string
	handler http.Handler
}

func (s *Server) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	address := net.JoinHostPort("localhost", s.port)

	return http_server.New(address, s.handler).Run(signals, ready)
}

func defaultDebugEndpoints() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

	mux.HandleFunc("/debug/requests", func(w http.ResponseWriter, req *http.Request) {
		any, sensitive := trace.AuthRequest(req)
		if !any {
			http.Error(w, "not allowed", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		trace.Render(w, req, sensitive)
	})

	mux.HandleFunc("/debug/events", func(w http.ResponseWriter, req *http.Request) {
		any, sensitive := trace.AuthRequest(req)
		if !any {
			http.Error(w, "not allowed", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		trace.RenderEvents(w, req, sensitive)
	})

	return mux
}

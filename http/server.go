package http

import "net/http"

type Middleware func(http.Handler) http.Handler

type Server struct {
	middlewares []Middleware
	mux         *http.ServeMux
	srv         *http.Server
}

func NewServer() *Server {
	return &Server{
		mux: http.NewServeMux(),
	}
}

// Use add a midleware to be wrap over handlers.
func (s *Server) Use(m Middleware) {
	s.middlewares = append(s.middlewares, m)
}

func (s *Server) Handle(path string, h http.Handler) {
	// add all middleware over the handler.
	for i := range s.middlewares {
		h = s.middlewares[len(s.middlewares)-1-i](h)
	}
	s.mux.Handle(path, h)
}

func (s *Server) Listen(addr string) error {
	s.srv = &http.Server{
		Handler: s.mux,
		Addr:    addr,
	}
	return s.srv.ListenAndServe()
}

package webserver

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]http.Handler
	WebServerPort string
}

func NewWebServer(WebServerPort string) *WebServer {
	return &WebServer{
		WebServerPort: WebServerPort,
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]http.Handler),
	}
}

func (s *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	for path, handler := range s.Handlers {
		s.Router.Handle(path, handler)
	}

	if err := http.ListenAndServe(s.WebServerPort, s.Router); err != nil {
		panic(err.Error())
	}
}

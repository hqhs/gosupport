package app

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Serve starts http server
func (s *Server) Serve() {

}

// InitRoutes initializes url schema. Separate function argument
// for routes is used to escape bugs there server tries to init
// routes without provided chi.Mux
func (s *Server) InitRoutes(router *chi.Mux) {
	s.r = router
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(CORSMiddleware)

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	router.Route("/api", func(r chi.Router) {
		r.Route("v1", s.APIRoutesv1)
	})
}

// APIRoutesv1 ...
func (s *Server) APIRoutesv1(r chi.Router) {

}

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
func (s *Server) InitRoutes(router chi.Router) {
	s.r = router
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(CORSMiddleware)

	router.HandleFunc("/login", loginForm)
	router.HandleFunc("/signin", signInForm)
	router.HandleFunc("/reset-password", resetPasswordForm)

	router.HandleFunc("/email/reset/{token}", emailResetRedirect)
	router.Get("/email/signin/{token}", emailSignInRedirect)

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	router.Mount("/", s.authorizedOnlyRoutes())
}

func (s *Server) authorizedOnlyRoutes() chi.Router {
	r := chi.NewRouter()
	r.Use(s.authMiddleware)

	r.HandleFunc("/settings", settingsForm)
	r.Get("/file/{id}", fileProxy)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/broadcast", broadcastMessage)
		// user resource
		r.Get("/user/", apiUserList)
		r.Route("/user/{userID}", func(r chi.Router) {
			r.Use(s.userCtx)
			r.Get("/", apiGetUserInfo)
			r.Get("/messages", apiGetUserMessages)
			r.Post("/messages", apiSendMessage)
		})
	})
	// TODO websocket management should be here
	return r
}

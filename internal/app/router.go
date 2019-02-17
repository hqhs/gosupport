package app

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// InitRoutes initializes url schema. Separate function argument
// for routes is used to escape bugs there server tries to init
// routes without provided chi.Mux
func (s *Server) InitRoutes(router chi.Router, static string) {
	s.router = router

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(CORSMiddleware)

	// NOTE serving static files is done before authorization
	if len(static) > 0 {
		fileServer(router, "/static", http.Dir(static))
	}

	router.With(s.RenderTemplate).HandleFunc("/login", s.loginForm)
	router.With(s.RenderTemplate).HandleFunc("/signin", s.signInForm)
	router.With(s.RenderTemplate).HandleFunc("/reset-password", s.resetPasswordForm)

	router.HandleFunc("/email/reset/{token}", s.emailResetRedirect)
	router.Get("/email/signin/{token}", s.emailSignInRedirect)

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	router.Mount("/", s.authorizedOnlyRoutes())
}

func (s *Server) authorizedOnlyRoutes() chi.Router {
	r := chi.NewRouter()
	r.Use(s.authMiddleware)

	r.With(s.RenderTemplate).Get("/", s.renderChatTemplate)
	r.HandleFunc("/settings", settingsForm)
	// FileProxy is used for serving media from chat, e.g. Telegram Photos
	r.Get("/file/{id}", s.fileProxy)

	// FIXME api is not versioned
	r.Route("/api", func(r chi.Router) {
		r.Post("/broadcast", broadcastMessage)
		// user resource
		r.Get("/user/", s.apiListUsers)
		r.Route("/user/{userID}", func(r chi.Router) {
			// r.Use(s.userCtx) TODO is this really necessary?
			r.Get("/", apiGetUserInfo)
			r.Get("/messages", s.apiGetUserMessages)
			r.Post("/messages", s.apiSendMessage)
		})
	})
	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		// TODO fetch bot hash from somewhere
		ctx := r.Context()
		bot := ctx.Value(botKey).(string)
		s.logger.Log("bot", bot)
		hub := s.hubs[bot]
		serveWs(hub, w, r)
	})
	return r
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

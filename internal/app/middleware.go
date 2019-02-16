package app

import (
	"context"
	"net/http"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/auth0/go-jwt-middleware"
	"github.com/go-chi/render"
)

type contextKey int

const (
	botKey contextKey = iota
	templateErrorKey
)

// CORSMidlleware writer cors headers to requests
func CORSMiddleware(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	JWTErrorHandler := func(w http.ResponseWriter, r *http.Request, err string) {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Secret), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		Extractor: func(r *http.Request) (string, error) {
			authCookie, err := r.Cookie("Authorization")
			if err != nil {
				return "", err
			}
			return authCookie.Value, nil
		},
		ErrorHandler: JWTErrorHandler,
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := jwtMiddleware.CheckJWT(w, r); err == nil {
			// Since there's no rights management at all yet, this middleware just adds
			// random bot hash to user cookies if hash not yet set. Since generated
			// only based on token and counter, this approach would work in most cases.
			ctx := r.Context()
			// I'm not quite sure error here is possible, middleware would redirect
			// to login page if jwt is not in request (or default ctx key was changed)
			token, _ := ctx.Value("user").(*jwt.Token)
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				bot := claims["bot"]
				// ctx.Set(botKey, bot)
				c := context.WithValue(ctx, botKey, bot)
				next.ServeHTTP(w, r.WithContext(c))
			} else {
				// TODO render 500
			}
		}
	})
}

// userCtx is used to load an objects
// the URL parameters passed through as the request. In case
// the user could not be found, we stop here and return a 404.
func (s *Server) userCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (s *Server) RenderTemplate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}
		var tmpl string
		switch r.URL.String() {
		case "/login":
			tmpl = "login.tmpl"
		case "/signin":
			tmpl = "signin.tmpl"
		case "/reset-password":
			tmpl = "reset_password.tmpl"
		case "/":
			tmpl = "chats.tmpl"
		}
		if err := s.templator.Render(w, tmpl, nil); err != nil {
			s.logger.Log("err", err, "then", fmt.Sprintf("during rendering %s template", tmpl))
			render.Render(w, r, errInternal)
		}
		return
	})
}

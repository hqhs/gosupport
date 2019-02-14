package app

import (
	"net/http"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/auth0/go-jwt-middleware"
	"github.com/go-chi/render"
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
			// ctx := r.Context()
			// token, ok := ctx.Value("user").(*jwt.Token)
			// if claims, ok := token.Claims.(jwt.MapClaims); ok {

			// }
			next.ServeHTTP(w, r)
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

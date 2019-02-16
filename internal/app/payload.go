package app

import (
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/go-chi/render"
	"github.com/dgrijalva/jwt-go"
)

type CustomJWTClaims struct {
	CurrentBot string `json:"bot"`
	jwt.StandardClaims
}

type loginData struct {
	Email string
	Password string
}

func (l *loginData) Bind(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	l.Email = r.Form.Get("email")
	l.Password = r.Form.Get("password")
	// mode this to standard checking utilities
	if len(l.Password) < 12 {
		return fmt.Errorf("Minimal password length is 12 characters")
	}
	if len(l.Email) == 0 {
		return fmt.Errorf("Fill in all required fields")
	}
	if _, err := mail.ParseAddress(l.Email); err != nil {
		return fmt.Errorf("Provided email is not valid")
	}
	return nil
}


type signInData struct {
	Email     string
	Password1 string
	Password2 string
}

func (s *signInData) Bind(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	// TODO additional checks
	s.Email = r.Form.Get("email")
	s.Password1 = r.Form.Get("password1")
	s.Password2 = r.Form.Get("password2")
	if len(s.Password1) < 12 || len(s.Password2) < 12 {
		return fmt.Errorf("Minimal password length is 12 characters")
	}
	if len(s.Email) == 0 {
		return fmt.Errorf("Fill in all required fields")
	}
	if _, err := mail.ParseAddress(s.Email); err != nil {
		return fmt.Errorf("Provided email is not valid")
	}
	if strings.Compare(s.Password1, s.Password2) != 0 {
		return fmt.Errorf("Passwords does not match")
	}
	return nil
}

type userResponse struct {
	User
}

func (s *Server) newUserResponse(u User) userResponse {
	return userResponse{u}
}

func (u userResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) newUserListResponse(users []User) []render.Renderer {
	list := []render.Renderer{}
	for _, u := range users {
		list = append(list, s.newUserResponse(u))
	}
	return list
}

type messageResponse struct {
	Message
}

func (m messageResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) newMessageResponse(m Message) messageResponse {
	return messageResponse{m}
}

func (s *Server) newMessageListResponse(msgs []Message) []render.Renderer {
	list := []render.Renderer{}
	for _, m := range msgs {
		list = append(list, s.newMessageResponse(m))
	}
	return list
}

//--
// Error response payloads & renderers
//--

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type errResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *errResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func errInvalidRequest(err error) render.Renderer {
	return &errResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func errRender(err error) render.Renderer {
	return &errResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

var errInvalid = &errResponse{HTTPStatusCode: 400, StatusText: "Invalid request."}
var errNotFound = &errResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}
var errInternal = &errResponse{HTTPStatusCode: 500, StatusText: "Internal server error."}

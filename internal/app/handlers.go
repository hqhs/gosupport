package app

import (
	"net/http"

	"github.com/go-chi/render"
)

func (s *Server) loginForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// TODO: rendering could be middleware
		if err := s.templator.Render(w, "login.tmpl", nil); err != nil {
			s.logger.Log("err", err, "then", "during rendering login template")
			// TODO: support custom error pages
			render.Render(w, r, errInternal)
		}
		return
	}
}

func (s *Server) signInForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if err := s.templator.Render(w, "signin.tmpl", nil); err != nil {
			s.logger.Log("err", err, "then", "during rendering singing template")
			render.Render(w, r, errInternal)
		}
		return
	}
}

func (s *Server) resetPasswordForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if err := s.templator.Render(w, "signin.tmpl", nil); err != nil {
			s.logger.Log("err", err, "then", "during rendering singing template")
			render.Render(w, r, errInternal)
		}
		return
	}
}

func emailResetRedirect(w http.ResponseWriter, r *http.Request) {

}

func emailSignInRedirect(w http.ResponseWriter, r *http.Request) {

}

func settingsForm(w http.ResponseWriter, r *http.Request) {

}

func fileProxy(w http.ResponseWriter, r *http.Request) {

}

func broadcastMessage(w http.ResponseWriter, r *http.Request) {

}

func apiUserList(w http.ResponseWriter, r *http.Request) {

}

func apiGetUserInfo(w http.ResponseWriter, r *http.Request) {

}

func apiGetUserMessages(w http.ResponseWriter, r *http.Request) {

}

func apiSendMessage(w http.ResponseWriter, r *http.Request) {

}

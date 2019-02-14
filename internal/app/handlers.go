package app

import (
	"fmt"
	"net/http"

	// "github.com/go-chi/render"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

const defaultPasswordHashingCost = 10

func (s *Server) loginForm(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) signInForm(w http.ResponseWriter, r *http.Request) {
	data := &signInData{}
	response := make(map[string]interface{})
	if err := data.Bind(r); err != nil {
		s.logger.Log("error", err.Error(), "then", "during binding signInForm data")
		// render.Render(w, r, errInvalid)
		response["Error"] = err.Error()
		s.templator.Render(w, "signin.tmpl", response)
		return
	}
	// TODO check what admin with provided email does not exist in database
	// if not, create one
	admin := Admin{}
	s.DB.Where("email = ?", data.Email).First(&admin)
	if len(admin.HashedPassword) > 0 {
		response["Error"] = fmt.Errorf("Admin with provided email already exists")
		s.templator.Render(w, "signin.tmpl", response)
		return
	}
	admin = Admin{Email: data.Email, IsActive: true, EmailConfirmed: false}
	hash, _ := bcrypt.GenerateFromPassword([]byte(data.Password1), defaultPasswordHashingCost)
	admin.HashedPassword = string(hash)
	if err := s.DB.Create(&admin).Error; err != nil {
		s.logger.Log("msg", "Admin creation error", "err", err)
		response["Error"] = fmt.Errorf("Internal server error. Try again or contact administrators.")
		s.templator.Render(w, "signin.tmpl", response)
		return
	}
	// TODO save admin to database
	token, _ := generateRandomStringURLSafe(60)
	url := fmt.Sprintf("http://%s/email/signin/%s", s.domain, token) // FIXME add optional port
	// TODO send authorization email with url
	mail := AuthMail{
		data.Email,
		"Authorization letter for Support Dashboard",
		"Click button below to authenticate.",
		url,
		"Confirm Email",
	}
	err := s.mailer.SendAuthMail(mail)
	if err != nil {
		response["Error"] = err.Error()
		s.templator.Render(w, "signin.tmpl", response)
		return
	}
	response["Message"] = "Success! Email sent."
	s.templator.Render(w, "success.tmpl", response)
}

func (s *Server) resetPasswordForm(w http.ResponseWriter, r *http.Request) {
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

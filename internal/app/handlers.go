package app

import (
	"fmt"
	"net/http"

	// "github.com/go-chi/render"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
)

const defaultPasswordHashingCost = 10

func (s *Server) loginForm(w http.ResponseWriter, r *http.Request) {
	data := &loginData{}
	response := make(map[string]interface{})
	if err := data.Bind(r); err != nil {
		// TODO add error to context and render it in templates middleware
		s.logger.Log("error", err.Error(), "then", "during binding signInForm data")
		response["Error"] = err.Error()
		s.templator.Render(w, "login.tmpl", response)
		return
	}
	admin := Admin{}
	s.DB.Where("email = ?", data.Email).First(&admin)
	// TODO check user is active and email is confirmed
	if bcrypt.CompareHashAndPassword([]byte(admin.HashedPassword), []byte(data.Password)) != nil {
		response["Error"] = "Email/Passwor pair is wrong. Try again or reset password."
		s.templator.Render(w, "login.tmpl", response)
		return
	}
	// TODO set expiration date
	claims := jwt.StandardClaims{Issuer: data.Email}
	JWTToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := JWTToken.SignedString([]byte(s.Secret))
	if err != nil {
		s.logger.Log("error", err.Error(), "then", "during signing jwt token")
		response["Error"] = "Internal server error. Try again or contact administrators."
		s.templator.Render(w, "login.tmpl", response)
		return
	}
	authCookie := http.Cookie{
		Name:     "Authorization",
		Value:    token,
		Domain:   s.domain,
		MaxAge:   3600,  // in seconds
		HttpOnly: false, // preact need this one
		SameSite: http.SameSiteStrictMode,
		// Secure:   true,     // allow only https
	}
	http.SetCookie(w, &authCookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Server) signInForm(w http.ResponseWriter, r *http.Request) {
	data := &signInData{}
	response := make(map[string]interface{})
	if err := data.Bind(r); err != nil {
		s.logger.Log("error", err.Error(), "then", "during binding signInForm data")
		response["Error"] = err.Error()
		s.templator.Render(w, "signin.tmpl", response)
		return
	}
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
	token, _ := generateRandomStringURLSafe(60)
	url := fmt.Sprintf("http://%s/email/signin/%s", s.domain, token) // FIXME add optional port
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

func (s *Server) emailSignInRedirect(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) resetPasswordForm(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) emailResetRedirect(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) renderChatTemplate(w http.ResponseWriter, r *http.Request) { /* this template renders in middleware */ }

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

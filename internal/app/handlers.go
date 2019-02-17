package app

import (
	"fmt"
	"net/http"
	"time"
	"io/ioutil"
	"strconv"
	"database/sql"
	"encoding/json"

	"github.com/go-chi/render"
	"github.com/go-chi/chi"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
)

const defaultPasswordHashingCost = 10

func (s *Server) loginForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := &loginData{}
	if err := data.Bind(r); err != nil {
		s.renderError(w, "login.tmpl", err)
		return
	}
	admin := &Admin{}
	query := "SELECT hashed_password FROM admins WHERE email = $1"
	if err := s.DB.QueryRowContext(ctx, query, data.Email).Scan(&admin.HashedPassword); err != nil {
		s.renderError(w, "login.tmpl", fmt.Errorf("Email/Passwor pair is wrong. Try again or reset password"))
		return
	}
	// TODO check user is active and email is confirmed
	if nil != bcrypt.CompareHashAndPassword([]byte(admin.HashedPassword), []byte(data.Password)) {
		s.renderError(w, "login.tmpl", fmt.Errorf("Email/Passwor pair is wrong. Try again or reset password"))
		return
	}
	// TODO set expiration date
	// add first available bot to claims, auth middleware would add it to context
	// later, read middleware comment for more information.
	var bot string
	for k := range s.conns {
		bot = k
		break
	}
	claims := CustomJWTClaims{bot, jwt.StandardClaims{Issuer: data.Email}}
	JWTToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := JWTToken.SignedString([]byte(s.Secret))
	if err != nil {
		s.renderError(w, "login.tmpl", err)
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
	ctx := r.Context()
	data := &signInData{}
	if err := data.Bind(r); err != nil {
		s.renderError(w, "signin.tmpl", err)
		return
	}
	admin := &Admin{}
	query := "SELECT email FROM admins WHERE email = $1"
	if err := s.DB.QueryRowContext(ctx, query, data.Email).Scan(&admin.Email); err != sql.ErrNoRows {
		s.renderError(w, "signin.tmpl", fmt.Errorf("Admin with provided email already exists"))
		return
	}
	admin = NewAdmin(data.Email, data.Password1, false)
	if err := dbCreateAdmin(ctx, s.DB, admin); err != nil {
		s.renderError(w, "signin.tmpl", fmt.Errorf("Internal server error. Try again or contact administrators"))
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
	if err := s.mailer.SendAuthMail(mail); err != nil {
		s.renderError(w, "signin.tmpl", err)
		return
	}
	s.templator.Render(w, "success.tmpl", map[string]interface{}{"Message": "Success! Email sent."})
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

func (s *Server) fileProxy(w http.ResponseWriter, r *http.Request) {
	var fileID string
	if fileID = chi.URLParam(r, "id"); fileID == "" {
		render.Render(w, r, errNotFound)
		return
	}
	ctx := r.Context()
	botHash, ok := ctx.Value(botKey).(string);
	if !ok {
		s.logger.Log("err", "botKey is missing in context")
		render.Render(w, r, errInternal)
		return
	}
	// FIXME site should not be started until all bots are connected
	s.logger.Log("botHash", botHash)
	bot := s.bots[botHash]
	URL, err := bot.GetFileDirectURL(fileID)
	if err != nil {
		s.logger.Log("err", err.Error(), "then", "during getting direct url for file download")
		render.Render(w, r, errInternal)
		return
	}
	client := http.DefaultClient
	resp, err := client.Get(URL)
	if err != nil {
		s.logger.Log("err", err.Error(), "then", "during downloading file in proxy")
		return
	}
	// TODO: add cache
	defer resp.Body.Close()
	buffer, _ := ioutil.ReadAll(resp.Body)
	contentType := http.DetectContentType(buffer)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprint(resp.ContentLength))
	w.Write(buffer)
}

func broadcastMessage(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) apiListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// TODO fetch page from context, if none, show 0
	users, err := dbListUsers(ctx, s.DB, 0)
	if err != nil {
		s.logger.Log("err", err.Error())
		render.Render(w, r, errInternal)
		return
	}
	if err := render.RenderList(w, r, s.newUserListResponse(users)); err != nil {
		s.logger.Log("err", err.Error())
		render.Render(w, r, errInternal)
	}
}

func apiGetUserInfo(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) apiGetUserMessages(w http.ResponseWriter, r *http.Request) {
	var userID string
	if userID = chi.URLParam(r, "userID"); userID == "" {
		render.Render(w, r, errNotFound)
		return
	}
	var id int
	var err error
	if id, err = strconv.Atoi(userID); err != nil {
		render.Render(w, r, errNotFound)
		return
	}
	ctx := r.Context()
	// TODO pagination (afterloading)
	messages, err := dbListUserMessages(ctx, s.DB, id, 0)
	if err := render.RenderList(w, r, s.newMessageListResponse(messages)); err != nil {
		s.logger.Log("err", err.Error(), "then", "during rendering message list")
		render.Render(w, r, errInternal)
	}
}

func (s *Server) apiSendMessage(w http.ResponseWriter, r *http.Request) {
	data := &messagePayload{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, errInvalid)
		return
	}
	msg := data.Message
	msg.FromAdmin = true
	msg.CreatedAt = time.Now()
	msg.UpdatedAt = time.Now()
	msg.MessageID = 100
	msg.ChatID = int64(msg.UserID)
	ctx := r.Context()
	// TODO save to database
	// err := dbCreateMessage(ctx, s.DB, msg)

	// send threw bot to customer
	botHash, ok := ctx.Value(botKey).(string);
	if !ok {
		s.logger.Log("err", "botKey is missing in context")
		render.Render(w, r, errInternal)
		return
	}
	s.conns[botHash].Output <- msg
	payload, _ := json.Marshal(msg)
	s.hubs[botHash].broadcast <- payload
}

func (s *Server) renderError(w http.ResponseWriter, t string, err error) {
	// FIXME not all errors should be rendered
	s.logger.Log("error", err.Error(), "then", fmt.Sprintf("during rendering %s template", t))
	s.templator.Render(w, "login.tmpl", map[string]interface{}{"Error": err.Error()})
}

package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"

	"errors"
	"log"
	"net/http"
	"time"
)

var cannotLogin = errors.New("The email or password is incorrect")

type login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// Login takes ...
func Login(context *gin.Context) {
	var json login
	if err := context.ShouldBindJSON(&json); err != nil {
		context.AbortWithError(http.StatusOK, err)
	}
	log.Printf("Email: %v, Password: %v\n", json.Email, json.Password)

	// try to get user from database
	db := context.MustGet("db").(*mgo.Database)
	hDesker, err := FetchHelpdesker(db, json.Email)
	if err != nil {
		context.AbortWithError(http.StatusOK, cannotLogin)
	}

	// compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(hDesker.PasswordHash), []byte(json.Password))
	if err != nil {
		context.AbortWithError(http.StatusOK, cannotLogin)
	}

	// Token generation
	// NOTE I don't neew query store every time helpdekser do something,
	// I only need two things: email, cos I'm sure it's unique, and list of available bots
	// expirationTime := time.Now().Add(time.Hour).Unix()
	requestedTime := time.Now().Unix()
	validUntil := time.Now().Add(24 * time.Hour).Unix()
	refreshUntil := time.Now().Add(168 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":          json.Email,
		"nbf":           requestedTime,
		"valid_until":   validUntil,
		"refresh_until": refreshUntil,
		"bots":          hDesker.AvailableBots,
	})

	tokenString, err := token.SignedString([]byte(Conf.Secret))

	if err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
	}

	context.JSON(http.StatusOK, gin.H{
		"success":       "OK",
		"token":         tokenString,
		"nbf":           requestedTime,
		"valid_until":   validUntil,
		"refresh_until": refreshUntil,
		"bots":          hDesker.AvailableBots,
	})
}

type refresh struct {
	Token string `form:"token" json:"token" binding:"required"`
}

// Refresh return new token with updated expirationTime
func Refresh(context *gin.Context) {
	var json refresh
	if err := context.ShouldBindJSON(&json); err != nil {
		context.AbortWithError(http.StatusOK, err)
	}
	// Parse token
	unvalidatedToken := ParseJWTHeader(json.Token)

	token, claims, err := ValidateJWT(unvalidatedToken)

	if err != nil {
		context.AbortWithError(http.StatusOK, err)
	}

	// token is valid, check time intervals
	currentTime := time.Now()
	refreshUntil := time.Unix(int64(claims["refresh_until"].(float64)), 0)
	if currentTime.After(refreshUntil) {
		context.AbortWithError(http.StatusOK, errTokenExpired)
	}
	// everything ok, generate new token with updated claims
	claims["nbf"] = time.Now().Unix()
	claims["refresh_until"] = time.Now().Add(168 * time.Hour).Unix()
	claims["valid_until"] = time.Now().Add(24 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(Conf.Secret))

	if err != nil {
		context.AbortWithError(http.StatusOK, err)
	}

	context.JSON(http.StatusOK, gin.H{
		"success": "OK",
		"token":   tokenString,
	})

}

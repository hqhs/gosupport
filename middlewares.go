package main

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/globalsign/mgo"
	"log"
	"net/http"
	"strings"
	"time"
)

// DatabaseMiddleware adds database session to each gin.Context
func DatabaseMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		s := Session.Clone()

		defer s.Close()

		context.Set("db", s.DB(Conf.DBName))
		context.Next()
	}
}

// AuthMiddleware checks jwt token, redirect to /auth if there's no one,
// and somewhere else if it's outdated
func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		// Parse JWT header
		authHeader := context.GetHeader("Authorization")
		unvalidatedToken := ParseJWTHeader(authHeader)

		log.Printf("token: %v\n", unvalidatedToken)

		token, claims, err := ValidateJWT(unvalidatedToken)

		if err != nil {
			context.AbortWithError(http.StatusBadRequest, err)
		}

		log.Printf("claims: %v\n", claims)

		// validate timing
		// NOTE: jwt-go checked nbf in previous part
		currentTime := time.Now()
		refreshUntil := time.Unix(int64(claims["refresh_until"].(float64)), 0)
		validUntil := time.Unix(int64(claims["valid_until"].(float64)), 0)
		if currentTime.After(refreshUntil) { // currentTime after refreshUntil
			context.AbortWithError(http.StatusUnauthorized, errTokenExpired)
		} else if currentTime.After(validUntil) {
			context.AbortWithError(http.StatusUnauthorized, errRefreshToken)
		}

		context.Next()

	}
}

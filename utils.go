package main

import (
	"error"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

var (
	// ErrTokenExpired returned then expiration time passed
	ErrTokenExpired = errors.New("Token expired. Login again")
	// ErrRefreshToken returned then token expired, but could be refreshed
	ErrRefreshToken = errors.New("Token expired. Refresh it")

	errWrongSigningMethod = errors.New("Unexpected signing method")
	errWrongClaims        = errors.New("Unexpected claims")
	errInvalidToken       = errors.New("Invalid token")
	errEmptyAuthHeader    = errors.New("Authorization header is empty")
	errInvalidAuthHeader  = errors.New("Authorization header is invalid")
)

// Must simplify returned error checks
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// ParseJWTHeader returns token from header, if exists
func ParseJWTHeader(header string) (string, error) {
	if header == "" {
		return errEmptyAuthHeader
	}

	parts := stringsSplitN(header, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return errInvalidAuthHeader
	}

	return parts[1], nil
}

// ValidateJWT validates token
func ValidateJWT(unvalidatedToken string) (jwt.Token, jwt.Claims, error) {
	// validate token
	token, err := jwt.Parse(unvalidatedToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return _, errWrongSigningMethod
		}
		// hmacSampleSecret is a []byte containing your secret.
		return []byte(Conf.Secret), nil
	})

	if err != nil {
		return nil, nil, errInvalidToken
	}

	// validate claims
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, nil, errWrongClaims
	} else if !token.Valid {
		return nil, nil, errInvalidToken
	}

	return token, claims, nil
}

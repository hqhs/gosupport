package main

import (
	"bytes"
	"json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Currently tmp database is not supported, so I assume that you've
// created user from the example
func TestAuthProcess(t *testing.T) {
	router := SetupRouter()

	w := httptest.NewRecorder()
	payload := json.Marshall(map[string]string{"email": "admin@example.com", "password": "admin"})
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(payload))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	log.Printf("Response body: %v", w.Body)
}

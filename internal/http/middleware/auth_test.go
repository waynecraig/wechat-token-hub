package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestAuthMiddleware(t *testing.T) {
	// create a test server with the Auth middleware
	handler := Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	server := httptest.NewServer(handler)
	defer server.Close()

	// Create a new JWT token with the specified header and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud": "wechat-token-hub",
		"exp": time.Now().Add(time.Second * 300).Unix(),
	})
	token.Header["kid"] = "key1"
	os.Setenv("JWT_KEY_key1", "secret1")

	// Sign the token with the specified secret
	tokenString, err := token.SignedString([]byte("secret1"))

	// Handle error if any
	if err != nil {
		t.Errorf("Error creating JWT token: %v", err)
	}

	// create a request with a valid JWT token in the Authorization header
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+tokenString)

	// send the request and check the response status code
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	// create a request with an invalid JWT token in the Authorization header
	req, err = http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer invalid_token")

	// send the request and check the response status code
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status code %d, but got %d", http.StatusUnauthorized, resp.StatusCode)
	}

	// create a request without Authorization header
	req, err = http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// send the request and check the response status code
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status code %d, but got %d", http.StatusUnauthorized, resp.StatusCode)
	}

	// create a request with a other token in the Authorization header
	req, err = http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Other "+tokenString)

	// send the request and check the response status code
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status code %d, but got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

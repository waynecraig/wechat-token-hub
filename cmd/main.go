package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	// get the port from env, default to 8567
	port := os.Getenv("PORT")
	if port == "" {
		port = "8567"
	}

	// set up the http server
	http.HandleFunc("/access_token", handleAccessToken)
	http.HandleFunc("/ticket", handleTicket)
	log.Fatal(http.ListenAndServe(":"+port, authMiddleware(http.DefaultServeMux)))
}

// authMiddleware is a middleware function that checks the authorization header for a valid JWT token
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// check that the authorization header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		// get the token string
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// get the secret from env
			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				return nil, fmt.Errorf("JWT_SECRET environment variable not set")
			}

			// verify the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// return the secret as the key
			return []byte(secret), nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// check if the token is valid
		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// call the next handler
		next.ServeHTTP(w, r)
	})
}

// handleAccessToken handles requests to the /access_token path
func handleAccessToken(w http.ResponseWriter, r *http.Request) {
	// TODO: handle access token requests
	// check if the token is in the cache
	tokenCache := getTokenFromCache()
	if tokenCache != "" {
		// if the token is in the cache, return it
		w.Write([]byte(tokenCache))
		return
	}

	// if the token is not in the cache, make a request to get it
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", os.Getenv("APPID"), os.Getenv("APPSECRET"))
	resp, err := http.Get(url)
	if err != nil {
		// if the request fails, return an error message
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to get access token"))
		return
	}
	defer resp.Body.Close()

	// decode the response body
	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		// if the response body cannot be decoded, return an error message
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to decode response body"))
		return
	}

	if result.AccessToken == "" {
		// if the response contains an error message, return it
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(result.ErrMsg))
		return
	}

	// save the token to the cache
	saveTokenToCache(result.AccessToken, result.ExpiresIn)

	// return the token
	w.Write([]byte(result.AccessToken))
}

// handleTicket handles requests to the /ticket path
func handleTicket(w http.ResponseWriter, r *http.Request) {
	// TODO: handle ticket requests
}

// define a struct to hold the token and its expiration time
type tokenInfo struct {
	Token      string
	Expiration time.Time
}

// define a variable to hold the token cache
var tokenCache *tokenInfo

// getTokenFromCache retrieves the token from the cache
func getTokenFromCache() string {
	// check if the token is expired
	if tokenCache != nil && tokenCache.Expiration.After(time.Now()) {
		// return the token cache
		return tokenCache.Token
	}

	// if the token is expired or not set, return an empty string
	return ""
}

// saveTokenToCache saves the token to the cache
func saveTokenToCache(token string, expiresIn int) {
	// calculate the expiration time
	expiration := time.Now().Add(time.Duration(expiresIn) * time.Second)

	// set the token cache to the new token and its expiration time
	tokenCache = &tokenInfo{
		Token:      token,
		Expiration: expiration,
	}
}

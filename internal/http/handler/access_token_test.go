package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/waynecraig/wechat-token-hub/internal/cache"
)

func TestAccessToken(t *testing.T) {
	handler := http.HandlerFunc(AccessToken)

	// Set the expect result to cache
	cache.SaveCacheItem("access_token", "token1", 3600)

	// Create a new request without a rotate_token query parameter
	req, err := http.NewRequest("GET", "/access_token", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the AccessToken function with the request and response recorder
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect
	if rr.Body.String() != "token1" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), "token1")
	}

	// Create a new request with a rotate_token query parameter
	req, err = http.NewRequest("GET", "/access_token?rotate_token=token1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr = httptest.NewRecorder()

	// Call the AccessToken function with the request and response recorder
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

}

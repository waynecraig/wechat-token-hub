package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOnlyGet(t *testing.T) {
	// Create a new request with a non-GET method
	req, err := http.NewRequest(http.MethodPost, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new response recorder
	rr := httptest.NewRecorder()

	// Create a new handler that wraps the OnlyGet middleware
	handler := OnlyGet(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Serve the request through the handler
	handler.ServeHTTP(rr, req)

	// Check that the response status code is 405 (Method Not Allowed)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}

	// Create a new request with a GET method
	req, err = http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Reset the response recorder
	rr = httptest.NewRecorder()

	// Serve the request through the handler again
	handler.ServeHTTP(rr, req)

	// Check that the response status code is 200 (OK)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

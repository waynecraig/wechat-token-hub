package middleware

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)

	// Create a mock handler
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create a mock request and response
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp := httptest.NewRecorder()

	// Call the Logger middleware with the mock handler
	loggerHandler := Logger(mockHandler)
	loggerHandler.ServeHTTP(resp, req)

	// Assert that the log output matches the expected format
	expectedLog := fmt.Sprintf("%s %s %s", http.MethodGet, "/test", http.StatusText(http.StatusOK))
	if !strings.Contains(logOutput.String(), expectedLog) {
		t.Errorf("Expected log output to contain '%s', but got '%s'", expectedLog, logOutput.String())
	}
}

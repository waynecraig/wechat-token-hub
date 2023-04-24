package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/waynecraig/wechat-token-hub/internal/cache"
)

func TestTicket(t *testing.T) {
	handler := http.HandlerFunc(Ticket)

	// Define test cases
	testCases := []struct {
		name           string
		ticketType     string
		rotateTicket   string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Get jsapi ticket",
			ticketType:     "jsapi",
			expectedStatus: http.StatusOK,
			expectedBody:   "ticket1",
		},
		{
			name:           "Get wx_card ticket",
			ticketType:     "wx_card",
			expectedStatus: http.StatusOK,
			expectedBody:   "ticket2",
		},
		{
			name:           "Get wx_card ticket with rotate ticket not match",
			ticketType:     "wx_card",
			rotateTicket:   "ticket1",
			expectedStatus: http.StatusOK,
			expectedBody:   "ticket2",
		},
		{
			name:           "Get wx_card ticket with rotate ticket match",
			ticketType:     "wx_card",
			rotateTicket:   "ticket2",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	// Set the expect result to cache
	cache.SaveCacheItem("ticket_jsapi", "ticket1", 3600)
	cache.SaveCacheItem("ticket_wx_card", "ticket2", 3600)

	// Loop through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new request
			url := fmt.Sprintf("/ticket?type=%s", tc.ticketType)
			if tc.rotateTicket != "" {
				url += fmt.Sprintf("&rotate_ticket=%s", tc.rotateTicket)
			}
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Call the handler function with the request and response recorder
			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}

			// Check the response body is what we expect
			if tc.expectedBody != "" && rr.Body.String() != tc.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tc.expectedBody)
			}
		})
	}
}

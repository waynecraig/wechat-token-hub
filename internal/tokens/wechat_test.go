package tokens

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func mockWechatServer(t *testing.T) *httptest.Server {
	// create a new http server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if the request path is "/cgi-bin/access-token"
		if r.URL.Path == "/cgi-bin/token" {
			// check if the query parameter "secret" is "secret1"
			if r.URL.Query().Get("secret") == "secret1" {
				// if "secret" is "secret1", return {"access_token":"token1","expires_in":7200}
				w.Write([]byte(`{"access_token":"token1","expires_in":7200}`))
			} else {
				// if "secret" is not "secret1", return {"errcode":40001,"errmsg":"invalid appsecret"}
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"errcode":40001,"errmsg":"invalid appsecret"}`))
			}
			// check if the request path is "/cgi-bin/ticket/getticket"
		} else if r.URL.Path == "/cgi-bin/ticket/getticket" {
			// special type, to test error branch
			if r.URL.Query().Get("type") == "not_support" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"errcode":40001,"errmsg":"invalid access token"}`))
			} else {
				// check if the query parameter "access_token" is "token1"
				if r.URL.Query().Get("access_token") == "token1" {
					// if "access_token" is "token1", return {ticket:"ticket1", expires_in:7200}
					w.Write([]byte(`{"ticket": "ticket1", "expires_in": 7200}`))
				} else {
					// if "access_token" is not "token1", return {"errcode":40001,"errmsg":"invalid access token"}
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"errcode":40001,"errmsg":"invalid access token"}`))
				}
			}
		} else {
			// if the request path is not "/cgi-bin/access-token" or "/cgi-bin/ticket", return 404
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 not found"))
		}
	}))
	// close the server when the test is finished
	t.Cleanup(func() {
		server.Close()
	})

	return server
}

func TestGetAccessToken(t *testing.T) {
	server := mockWechatServer(t)
	os.Setenv("WECHAT_API_ROOT", server.URL)

	tests := []struct {
		name        string
		rotateToken string
		envVars     map[string]string
		accessToken string
		wantErr     bool
	}{
		{
			name:        "secret error",
			rotateToken: "",
			envVars: map[string]string{
				"APPID":     "app1",
				"APPSECRET": "secret2",
			},
			accessToken: "",
			wantErr:     true,
		},
		{
			name:        "normal",
			rotateToken: "",
			envVars: map[string]string{
				"APPID":     "app1",
				"APPSECRET": "secret1",
			},
			accessToken: "token1",
			wantErr:     false,
		},
		{
			name:        "from cache",
			rotateToken: "",
			envVars: map[string]string{
				"APPID":     "app1",
				"APPSECRET": "secret1",
			},
			accessToken: "token1",
			wantErr:     false,
		},
		{
			name:        "roate token",
			rotateToken: "token1",
			envVars: map[string]string{
				"APPID":     "app1",
				"APPSECRET": "secret1",
			},
			accessToken: "token1",
			wantErr:     false,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			// Call function under test
			token, err := GetAccessToken(tt.rotateToken)

			// Check result
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccessToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if token != tt.accessToken {
				t.Errorf("Expect accessToken = %s, got %s", tt.accessToken, token)
			}

			// Clean up environment variables
			for k := range tt.envVars {
				os.Unsetenv(k)
			}
		})
	}

	os.Unsetenv("WECHAT_API_ROOT")
}

func TestGetTicket(t *testing.T) {
	server := mockWechatServer(t)
	os.Setenv("WECHAT_API_ROOT", server.URL)

	tests := []struct {
		name         string
		ticketType   string
		rotateTicket string
		envVars      map[string]string
		ticket       string
		wantErr      bool
	}{
		{
			name:         "type error",
			ticketType:   "not_support",
			rotateTicket: "",
			envVars: map[string]string{
				"APPID":     "app1",
				"APPSECRET": "secret1",
			},
			ticket:  "",
			wantErr: true,
		},
		{
			name:         "normal",
			ticketType:   "jsapi",
			rotateTicket: "",
			envVars: map[string]string{
				"APPID":     "app1",
				"APPSECRET": "secret1",
			},
			ticket:  "ticket1",
			wantErr: false,
		},
		{
			name:         "from cache",
			ticketType:   "jsapi",
			rotateTicket: "",
			envVars: map[string]string{
				"APPID":     "app1",
				"APPSECRET": "secret1",
			},
			ticket:  "ticket1",
			wantErr: false,
		},
		{
			name:         "rotate ticket",
			ticketType:   "jsapi",
			rotateTicket: "ticket1",
			envVars: map[string]string{
				"APPID":     "app1",
				"APPSECRET": "secret1",
			},
			ticket:  "ticket1",
			wantErr: false,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			// Call function under test
			ticket, err := GetTicket(tt.ticketType, tt.rotateTicket)

			// Check result
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTicket() error = %v, wantErr %v", err, tt.wantErr)
			}
			if ticket != tt.ticket {
				t.Errorf("Expect ticket = %s, got %s", tt.ticket, ticket)
			}

			// Clean up environment variables
			for k := range tt.envVars {
				os.Unsetenv(k)
			}
		})
	}

	os.Unsetenv("WECHAT_API_ROOT")
}

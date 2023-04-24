package auth

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestVerifyJwtToken(t *testing.T) {
	// Create a new JWT token with the specified header and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud": "wechat-token-hub",
		"exp": time.Now().Add(time.Second * 300).Unix(),
	})
	token.Header["kid"] = "key1"

	// Sign the token with the specified secret
	tokenString, err := token.SignedString([]byte("secret1"))

	// Handle error if any
	if err != nil {
		t.Errorf("Error creating JWT token: %v", err)
	}

	// Create an invalid token with wrong method
	invalidToken1 := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"aud": "wechat-token-hub",
		"exp": time.Now().Add(time.Second * 300).Unix(),
	})
	invalidToken1.Header["kid"] = "key1"

	// Sign the invalid token with the specified secret
	invalidTokenString1, err := invalidToken1.SignedString([]byte("secret1"))

	// Handle error if any
	if err != nil {
		t.Errorf("Error creating JWT token: %v", err)
	}

	// Create an invalid token with wrong audiend
	invalidToken2 := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud": "other service",
		"exp": time.Now().Add(time.Second * 300).Unix(),
	})
	invalidToken2.Header["kid"] = "key1"

	// Sign the invalid token with the specified secret
	invalidTokenString2, err := invalidToken2.SignedString([]byte("secret1"))

	// Handle error if any
	if err != nil {
		t.Errorf("Error creating JWT token: %v", err)
	}

	// Create an invalid token expired
	invalidToken3 := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud": "wechat-token-hub",
		"exp": time.Now().Add(time.Second * -300).Unix(),
	})
	invalidToken3.Header["kid"] = "key1"

	// Sign the invalid token with the specified secret
	invalidTokenString3, err := invalidToken3.SignedString([]byte("secret1"))

	// Handle error if any
	if err != nil {
		t.Errorf("Error creating JWT token: %v", err)
	}

	// Create an invalid token without kid
	invalidToken4 := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud": "wechat-token-hub",
		"exp": time.Now().Add(time.Second * 300).Unix(),
	})

	// Sign the invalid token with the specified secret
	invalidTokenString4, err := invalidToken4.SignedString([]byte("secret1"))

	// Handle error if any
	if err != nil {
		t.Errorf("Error creating JWT token: %v", err)
	}

	// Create an invalid token with invalid signature
	invalidToken5 := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud": "wechat-token-hub",
		"exp": time.Now().Add(time.Second * 300).Unix(),
	})
	invalidToken5.Header["kid"] = "key1"

	// Sign the invalid token with the specified secret
	invalidTokenString5, err := invalidToken5.SignedString([]byte("secret2"))

	// Handle error if any
	if err != nil {
		t.Errorf("Error creating JWT token: %v", err)
	}

	// Use the created token string in the test case
	tests := []struct {
		name        string
		tokenString string
		envVars     map[string]string
		wantErr     bool
	}{
		{
			name:        "valid token",
			tokenString: tokenString,
			envVars: map[string]string{
				"JWT_KEY_key1": "secret1",
			},
			wantErr: false,
		},
		{
			name:        "invalid token using HS512",
			tokenString: invalidTokenString1,
			envVars: map[string]string{
				"JWT_KEY_key1": "secret1",
			},
			wantErr: true,
		},
		{
			name:        "invalid token with invalid audience",
			tokenString: invalidTokenString2,
			envVars: map[string]string{
				"JWT_KEY_key1": "secret1",
			},
			wantErr: true,
		},
		{
			name:        "invalid token expired",
			tokenString: invalidTokenString3,
			envVars: map[string]string{
				"JWT_KEY_key1": "secret1",
			},
			wantErr: true,
		},
		{
			name:        "invalid token without kid",
			tokenString: invalidTokenString4,
			envVars: map[string]string{
				"JWT_KEY_key1": "secret1",
			},
			wantErr: true,
		},
		{
			name:        "invalid token with invalid signature",
			tokenString: invalidTokenString5,
			envVars: map[string]string{
				"JWT_KEY_key1": "secret1",
			},
			wantErr: true,
		},
		{
			name:        "missing env var",
			tokenString: tokenString,
			envVars:     map[string]string{},
			wantErr:     true,
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
			err := VerifyJwtToken(tt.tokenString)

			// Check result
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyJwtToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Clean up environment variables
			for k := range tt.envVars {
				os.Unsetenv(k)
			}
		})
	}
}

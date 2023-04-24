package auth

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// a function that parses and validates a JWT token
func VerifyJwtToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// get key id from token header
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid is not a string")
		}

		// get the secret from env
		secret := os.Getenv("JWT_KEY_" + kid)
		if secret == "" {
			return nil, fmt.Errorf("JWT_KEY_%s environment variable not set", kid)
		}

		// return the secret as the key
		return []byte(secret), nil
	},
		jwt.WithAudience("wechat-token-hub"),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
	if err != nil {
		return err
	}

	// check if the token is valid
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

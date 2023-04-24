package handler

import (
	"net/http"

	"github.com/waynecraig/wechat-token-hub/internal/tokens"
)

// AccessToken handles requests to the /access_token path
func AccessToken(w http.ResponseWriter, r *http.Request) {
	rotateToken := r.URL.Query().Get("rotate_token")
	accessToken, err := tokens.GetAccessToken(rotateToken)
	if err != nil {
		// return error if get access token fail
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	// return the access token
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(accessToken))
}

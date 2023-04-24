package tokens

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/waynecraig/wechat-token-hub/internal/cache"
)

func GetAccessToken(rotateToken string) (string, error) {
	// check if the access token is in the cache
	accessToken := cache.GetCacheItem("access_token")
	if accessToken != "" {
		// check if the rotate token is the same as the cached access token
		if rotateToken != "" && rotateToken == accessToken {
			// if the rotate token is the same as the cached access token, request a new one
			accessToken = ""
		} else {
			// if the access token is in the cache and not expired, return it
			return accessToken, nil
		}
	}

	// if the access token is not in the cache or needs to be refreshed, make a request to get it
	accessToken, err := retrieveAccessToken()
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func GetTicket(ticketType string, rotateTicket string) (string, error) {
	// check if the ticket is in the cache
	ticket := cache.GetCacheItem("ticket_" + ticketType)
	if ticket != "" {
		// check if the rotate ticket is the same as the cached ticket
		if rotateTicket != "" && rotateTicket == ticket {
			// if the rotate ticket is the same as the cached ticket, request a new one
			ticket = ""
		} else {
			// if the ticket is in the cache and not expired, return it
			return ticket, nil
		}
	}

	// if the ticket is not in the cache or needs to be refreshed, make a request to get it
	accessToken, err := GetAccessToken("")
	if err != nil {
		return "", err
	}
	ticket, err = retrieveTicket(accessToken, ticketType)

	// if get the 40001 error code, means the access token is expired, rotate it.
	if err != nil && strings.Contains(err.Error(), "40001") {
		accessToken, err = GetAccessToken(accessToken)
		if err != nil {
			return "", err
		}
		ticket, err = retrieveTicket(accessToken, ticketType)
	}

	if err != nil {
		return "", err
	}
	return ticket, nil
}

// request the wechat API to get a new access token
func retrieveAccessToken() (string, error) {
	url := fmt.Sprintf("%s/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", os.Getenv("WECHAT_API_ROOT"), os.Getenv("APPID"), os.Getenv("APPSECRET"))
	resp, err := http.Get(url)
	if err != nil {
		return "", err
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
		return "", err
	}

	if result.AccessToken == "" {
		return "", fmt.Errorf("fetch access token fail, code: %d, message: %s", result.ErrCode, result.ErrMsg)
	}

	// save the access token to the cache
	cache.SaveCacheItem("access_token", result.AccessToken, result.ExpiresIn)

	return result.AccessToken, nil
}

// request the wechat API to get a new ticket
func retrieveTicket(accessToken string, ticketType string) (string, error) {
	url := fmt.Sprintf("%s/cgi-bin/ticket/getticket?access_token=%s&type=%s", os.Getenv("WECHAT_API_ROOT"), accessToken, ticketType)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// decode the response body
	var result struct {
		Ticket    string `json:"ticket"`
		ExpiresIn int    `json:"expires_in"`
		ErrCode   int    `json:"errcode"`
		ErrMsg    string `json:"errmsg"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	if result.Ticket == "" {
		return "", fmt.Errorf("fetch ticket fail, code: %d, message: %s", result.ErrCode, result.ErrMsg)
	}

	// save the ticket to the cache
	cache.SaveCacheItem("ticket_"+ticketType, result.Ticket, result.ExpiresIn)

	return result.Ticket, nil
}

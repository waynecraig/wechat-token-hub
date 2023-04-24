package main

import (
	"log"
	"net/http"
	"os"

	"github.com/waynecraig/wechat-token-hub/internal/http/handler"
	mw "github.com/waynecraig/wechat-token-hub/internal/http/middleware"
)

func main() {
	// get the port from env, default to 8567
	port := os.Getenv("PORT")
	if port == "" {
		port = "8567"
	}
	log.Printf("use port %s", port)

	// set up the http server
	http.HandleFunc("/access_token", handler.AccessToken)
	http.HandleFunc("/ticket", handler.Ticket)
	log.Fatal(http.ListenAndServe(":"+port, mw.Logger(mw.OnlyGet(mw.Auth(http.DefaultServeMux)))))
}

package main

import (
	"net/http"

	"github.com/go-acd/token-server/goapp"
)

func init() {
	http.HandleFunc("/", goapp.RootHandler)
	http.HandleFunc("/oauth", goapp.OauthHandler)
	http.HandleFunc("/refresh", goapp.RefreshHandler)
}

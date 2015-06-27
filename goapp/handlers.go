package goapp

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"google.golang.org/appengine"
)

var (
	oauthConfig = oauth2.Config{
		// clientID is declared in credentials.go file which is encrypted.
		ClientID: clientID,
		// clientSecret is declared in credentials.go file which is encrypted.
		ClientSecret: clientSecret,
		RedirectURL:  "https://go-acd.appspot.com/oauth",
		Scopes:       []string{"clouddrive:read", "clouddrive:write"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.amazon.com/ap/oa",
			TokenURL: "https://api.amazon.com/auth/o2/token",
		},
	}

	indexTmpl = template.Must(template.ParseFiles("index.html"))
)

// RootHandler is the root page handler
func RootHandler(w http.ResponseWriter, r *http.Request) {
	authURL := oauthConfig.AuthCodeURL("go-acd", oauth2.AccessTypeOffline)
	indexTmpl.Execute(w, authURL)
}

// OauthHandler ...
func OauthHandler(w http.ResponseWriter, r *http.Request) {
	if state := r.URL.Query().Get("state"); state != "go-acd" {
		log.Print("error getting the state from the request")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	token, err := oauthConfig.Exchange(appengine.NewContext(r), r.URL.Query().Get("code"))
	if err != nil {
		log.Printf("error exchanging the code for an access token: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Disposition", "attachment; filename=acd-token.json")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(token); err != nil {
		log.Printf("error encoding the token in JSON: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// RefreshHandler ...
func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var token oauth2.Token
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		log.Printf("error decoding the token from the request: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tokenSource := oauthConfig.TokenSource(appengine.NewContext(r), &token)
	newToken, err := tokenSource.Token()
	if err != nil {
		log.Printf("error fetching a new token: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(newToken); err != nil {
		log.Printf("error encoding the token in JSON: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

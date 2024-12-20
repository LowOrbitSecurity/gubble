package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/savioxavier/termlink"
	"golang.org/x/oauth2"
)

var TICK string = ("[" + color.GreenString("+") + ("]"))

// Retrieve a token, saves the token, then returns the generated client.
func GetClient(config *oauth2.Config) *http.Client {
	tok := getTokenFromWeb(config)
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	// Create a channel to receive the authorization code
	codeCh := make(chan string)

	go func() {

		// Start a local web server to handle the highirect
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			authCode := r.URL.Query().Get("code")
			codeCh <- authCode
			fmt.Fprintln(w, "Authorization code received. You can close this window.")
		})

		server := &http.Server{
			Addr:              ":80",
			ReadHeaderTimeout: 3 * time.Second,
		}

		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Web server error: %v", err)
		}

	}()

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println(TICK, termlink.ColorLink("Click here to authenticate to GCP 🔗", authURL, "italic blue"))

	// Wait for the authorization code from the web server
	authCode := <-codeCh

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

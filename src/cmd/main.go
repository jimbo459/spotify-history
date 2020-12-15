package main

import (
	"fmt"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
	"time"
)

var (
	redirectUrl = "http://localhost:3000/callback"
	state = "abc123"
	auth = spotify.NewAuthenticator(redirectUrl, spotify.ScopeUserReadRecentlyPlayed)
	ch = make(chan *spotify.Client)
)

func main() {
	// Start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %v\n", r.URL.String())
	})
	go http.ListenAndServe(":3000", nil)

	url := auth.AuthURL(state)
	fmt.Printf("Please navigate to this URL to authenticate: %v\n", url)

	// need to understand this better...
	client := <-ch

	for x := 0; x < 30; x++ {
		lastPlayed, err := client.PlayerRecentlyPlayed()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Attempt %v: Recently Played %v\n", x, lastPlayed[0])
		time.Sleep(2 * time.Minute)
	}

}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	token, err := auth.Token(state, r)
	if err != nil {
		http.Error(w,"Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w,r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	client := auth.NewClient(token)
	fmt.Fprintf(w,"Login complete")
	ch <- &client
}
package handlers

import (
	"fmt"
	spotify "github.com/zmb3/spotify"
	"log"
	"net/http"
)

var (
	Ch = make(chan *spotify.Client)
	Auth = spotify.NewAuthenticator(redirectUrl, spotify.ScopeUserReadRecentlyPlayed)
	redirectUrl = "http://localhost:3000/callback"
	State = "abc123"
)

func CallBackHandler(w http.ResponseWriter, r *http.Request) {
		token, err := Auth.Token(State, r)
		if err != nil {
			http.Error(w,"Couldn't get token", http.StatusForbidden)
			log.Fatal(err)
		}
		if st := r.FormValue("state"); st != State {
			http.NotFound(w,r)
			log.Fatalf("State mismatch: %s != %s\n", st, State)
		}


		client := Auth.NewClient(token)
		fmt.Fprintf(w,"Login complete")
		Ch <- &client
}
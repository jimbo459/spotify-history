package main

import (
	"fmt"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
	"os"
)

const redirectURI="http://localhost:3000/callback"

type JClient spotify.Client

var(
	auth = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadPrivate)
	ch = make(chan *spotify.Client)
	state = "abasd123"
)

func completeAuth(w http.ResponseWriter, r * http.Request) {
	token, err := auth.Token(state, r)


	if err != nil {
		http.Error(w, "Could not get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w,r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	client := auth.NewClient(token)
	fmt.Fprintf(w, "Login Completed!")

	ch <-&client

}

func main() {

	http.HandleFunc("/callback", completeAuth)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})

	go http.ListenAndServe(":3000", nil)

	auth.SetAuthInfo(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))
	url := auth.AuthURL(state)

	fmt.Println("Please log into Spotify by visiting the following page in your browser: %v", url)

	JClient := <-ch

	user, err := JClient.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success. Logged in as:", user.DisplayName)

	history, err := JClient.RecentlyPlayed()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(history)

}

func (c *JClient) RecentlyPlayed(opt *Options) (*RecentlyPlayed, error) {
	spotifyURL := c.baseURL + "me/player/recently-played"

	var result RecentlyPlayed

	err := c.get(spotifyURL, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}



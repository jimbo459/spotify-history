package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
  "time"

	"github.com/zmb3/spotify"
)

const redirectURI = "http://localhost:3000/callback"

var (
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadRecentlyPlayed)
	ch    = make(chan *spotify.Client)
	state = "abasd123"
)

type Track struct {
  trackID string
  trackTitle string
  playedAt time.Time
}


func completeAuth(w http.ResponseWriter, r *http.Request) {
	token, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Could not get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	client := auth.NewClient(token)
	fmt.Fprintf(w, "Login Completed!")

	ch <- &client

}

func main() {
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})

	go http.ListenAndServe(":3000", nil)

	auth.SetAuthInfo(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))
	url := auth.AuthURL(state)

	fmt.Println("Please log into Spotify by visiting the following page in your browser:", url)

	client := <-ch

	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success. Logged in as:", user.DisplayName)

	options := spotify.RecentlyPlayedOptions{
		Limit:         50,
	}
	recentlyPlayed, err := client.PlayerRecentlyPlayedOpt(options)
	if err != nil {
		log.Fatal(err)
	}

  track := Track{
    trackID: string(recentlyPlayed[0].Track.ID),
    trackTitle: recentlyPlayed[0].Track.Name,
	playedAt: recentlyPlayed[0].PlayedAt,
  }

	fmt.Printf("Track Object=%v, Raw PlayedAt Output:=%s", track, recentlyPlayed[0].PlayedAt)
}

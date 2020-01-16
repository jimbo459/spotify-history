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

type playEntry struct {
	trackID  string
	playedAt time.Time
}

var playHistory []playEntry
var trackLibrary []spotify.SimpleTrack

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

	recentlyPlayed := getRecentlyPlayed(client)
	lastPlayed := playHistory[len(playHistory)-1].playedAt

	for _, trackReturn := range recentlyPlayed {
		if trackReturn.PlayedAt > lastPlayed {
			tmpTrack := playEntry{
				trackID:  string(trackReturn.Track.ID),
				playedAt: trackReturn.PlayedAt,
			}
			playHistory = append(playHistory, tmpTrack)
		}
	}

	fmt.Printf("Play history: %v", playHistory)

}

func getRecentlyPlayed(client *spotify.Client) []spotify.RecentlyPlayedItem {
	options := &spotify.RecentlyPlayedOptions{Limit: 50}

	recentlyPlayed, err := client.PlayerRecentlyPlayedOpt(options)
	if err != nil {
		log.Fatal(err)
	}

	return recentlyPlayed

}

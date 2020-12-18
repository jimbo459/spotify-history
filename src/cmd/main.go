package main

import (
	"database/sql"
	"fmt"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	"strings"
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

	lastPlayed, err := client.PlayerRecentlyPlayed()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Successfully got last played tracks")

	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/test_db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlstmt, er := db.Prepare("INSERT INTO play_history(played_at, track_name, track_id, artist_name, artist_id) VALUES (?,?,?,?,?)")
	if er != nil {
		log.Fatal(err)
	}

	for _,track := range lastPlayed{
		var artist []string
		var artistId []string

		for _, tempArtist := range track.Track.Artists {
			artist = append(artist, tempArtist.Name)
			artistId = append(artistId, string(tempArtist.ID))
		}
		_,err = sqlstmt.Exec(track.PlayedAt, track.Track.Name, track.Track.ID, strings.Join(artist, ","), strings.Join(artistId, ","))
		if err != nil {
			log.Fatal(err)
		}
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
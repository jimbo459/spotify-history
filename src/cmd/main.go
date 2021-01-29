package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jimbo459/spotify-history/src/internal/handlers"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)


func main() {
	//Get env vars for db
	dbUser := os.Getenv("SQL_USER")
	dbPassword := os.Getenv("SQL_PASSWORD")
	if len(dbUser) == 0 || len(dbPassword) == 0 {
		fmt.Printf("Error SQL_USER & SQL_PASSWORD must be set before start\n")
		os.Exit(1)
	}

	// Start an HTTP server
	http.HandleFunc("/callback", handlers.CallBackHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %v\n", r.URL.String())
	})

	go http.ListenAndServe(":3000", nil)

	url := handlers.Auth.AuthURL(handlers.State)
	fmt.Printf("Please navigate to this URL to authenticate: %v\n", url)

	// need to understand this better...
	client := <-handlers.Ch

	recentlyPlayedOpt := &spotify.RecentlyPlayedOptions{
		Limit:         50,
		AfterEpochMs:  0,
		BeforeEpochMs: 0,
	}
	for {
		lastPlayed, err := client.PlayerRecentlyPlayedOpt(recentlyPlayedOpt)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Got play history at: %v", time.Now().String())

		db, err := sql.Open("mysql", dbUser + ":" + dbPassword +"@tcp(127.0.0.1:3306)/test_db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		sqlstmt, er := db.Prepare("INSERT INTO play_history(played_at, track_name, track_id, artist_name, artist_id) VALUES (?,?,?,?,?)")
		if er != nil {
			log.Fatal(err)
		}

		err = writeHistory(sqlstmt, lastPlayed)
		if err != nil {
			log.Fatal(err)
		}

		err = db.Close()
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(90 * time.Minute)
	}
}

func writeHistory(statement *sql.Stmt, lastPlayed []spotify.RecentlyPlayedItem) error {
	for _,track := range lastPlayed{
		var artist []string
		var artistId []string

		for _, tempArtist := range track.Track.Artists {
			artist = append(artist, tempArtist.Name)
			artistId = append(artistId, string(tempArtist.ID))
		}
		_, err := statement.Exec(track.PlayedAt, track.Track.Name, track.Track.ID, strings.Join(artist, ","), strings.Join(artistId, ","))
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}
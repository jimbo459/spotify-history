package handlers

import (
	"fmt"
	"log"
	"net/http"
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

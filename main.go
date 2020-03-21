package main

import (
	"fmt"
	"log"
	"net/http"

	"./spotify"
)

const (
	ScopeUserReadPrivate  = "user-read-private"
	ScopeUserFollowRead   = "user-follow-read"
	ScopeUserFollowModify = "user-follow-modify"
)

// redirectURI is the OAuth reditect URI for the application
const RedirectURI = "http://localhost:8080/callback"

// Instance of an error type, auth Client and the state
var (
	err     error
	auth    = spotify.NewAuthenticator(RedirectURI, ScopeUserReadPrivate, ScopeUserFollowRead, ScopeUserFollowModify)
	channel = make(chan *spotify.Client)
	state   = "abc123"
)

func main() {

	// Calls to the OAuth
	http.HandleFunc("/callback", completeAuthorization)

	// Register the handle function with the 'Get User's Followed Artist' pattern
	http.HandleFunc("/me/following?type=artist", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	go func() {
		url := auth.AuthURL(state)
		fmt.Println("Please log in to Spotify : ", url)

		//wait for the auth to complete
		client := <-channel

		// use the client to make calls that require authorization
		user, err := client.CurrentUser()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("You are logged in as : ", user.DisplayName)

		artists, err := client.FollowedList(50, "")
		if err != nil {
			log.Fatal(err)
		}

		for i, a := range artists.Artists {
			fmt.Println(i, a.Name)
		}
	}()

	http.ListenAndServe(":8080", nil)
}

func completeAuthorization(w http.ResponseWriter, r *http.Request) {

	// Get the token using the new authenticator
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token.", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch : %s != %s\n", st, state)
	}

	// if the state returned is good and we received a token
	// Use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	channel <- &client
}

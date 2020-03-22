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
	// RedirectURI is the OAuth reditect URI for the application
	RedirectURI = "http://localhost:8080/callback"
)

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
	//http.HandleFunc("/me/following?type=artist", func(w http.ResponseWriter, r *http.Request) {})
	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	go http.ListenAndServe(":8080", nil)

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify : ", url)

	//wait for the auth to complete
	client := <-channel

	GetCurrentUser(client)

	// Get the list of all the artists followed
	artists := GetFollowedArtists(client)
	//GetArtistAlbums(artists)
	PrintFollowedArtists(artists)

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

func GetCurrentUser(client *spotify.Client) {
	// use the client to make calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as : ", user.DisplayName)
}

func GetFollowedArtists(client *spotify.Client) []spotify.Artist {
	var lastArtistID = ""
	var artistList = []spotify.Artist{}
	for {
		// Get the a list of followed artists
		artists, err := client.FollowedList(50, lastArtistID)
		artistList = append(artistList, artists.Artists...)
		if err != nil {
			log.Fatal(err)
		}
		// If there is no other pages, break out of the loop
		if artists.CursorBasedObj.Next == "" {
			break
		}
		lastArtistID = artistList[len(artistList)-1].ID
	}
	return artistList
}

func PrintFollowedArtists(artists []spotify.Artist) {
	for i, a := range artists {
		fmt.Println(i, a.Name)
	}
}

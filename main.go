package main

import (
	"fmt"
	"log"
	"net/http"

	"./spotify"
)

const ScopeUserReadPrivate = "user-read-private"

// redirectURI is the OAuth reditect URI for the application
const RedirectURI = "http://localhost:8080/callback"

// Instance of an error type, auth Client and the state
var (
	err     error
	auth    = spotify.NewAuthenticator(RedirectURI, ScopeUserReadPrivate)
	channel = make(chan *spotify.Client)
	state   = "abc123"
)

func main() {

	http.HandleFunc("/callback", completeAuthorization)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for : ", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

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

	// Creating a connection to DB
	//Config.DB, err = gorm.Open("mysql", "root:mypassword@/followers")

	/*if err != nil {
		fmt.Println("Status: ", err)
	}*/

	// The defer statement defers the execution of a func until the surroundings function returs
	//defer Config.DB.Close()

	// run the migrations: follower struct
	//Config.DB.AutoMigrate(&Models.Follower{})

	// setup routes
	//r := Routes.SetupRouter()

	// running
	//r.Run(":8080")

}

func completeAuthorization(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token.", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch : %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	channel <- &client
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"./spotify"
)

const (
	ScopeUserReadPrivate  = "user-read-private"
	ScopeUserFollowRead   = "user-follow-read"
	ScopeUserFollowModify = "user-follow-modify"
	// RedirectURI is the OAuth reditect URI for the application
	RedirectURI = "http://localhost:8080/callback"
	layoutISO   = "2006-01-02"
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
	followedArtists := GetFollowedArtists(client)

	/*latestReleases := */
	GetFollowedArtistsLatest(followedArtists, client)

	//AddLatestReleasesToPlaylist()
	//PrintFollowedArtists(artists)

}

func GetFollowedArtistsLatest(followedArtists []spotify.Artist, client *spotify.Client) /*[]*spotify.SimplifiedAlbumObject*/ {
	//var newReleases = []*spotify.SimplifiedAlbumObject{}

	artistsAlbums := GetFollowedArtistAlbums(client, followedArtists)
	for _, aa := range artistsAlbums {
		// Check if albums released less than a month ago
		/*monthlyReleases :=*/
		GetMonthyReleases(aa)
		//newReleases = append(newReleases, monthlyReleases)
		// if so add to newReleases slice
		// TODO : Should have a bd or keep track of the slice containing the latest releases
		// so that if I run the code two days in a row I only get the ones that released in the past day

	}

	//return newReleases
}

// GetMonthyReleases : Check if albums released less than a month ago
// TODO: Should implement a way to fetch a sclice (newReleases) to only return
// the ones that were not been fetch allready this month from another use of the function
func GetMonthyReleases(artistAlbum *spotify.SimplifiedAlbumObject) /*spotify.SimplifiedAlbumObject*/ {
	formReleaseDate, _ := time.Parse(layoutISO, artistAlbum.ReleaseDate)
	timeDiff := time.Since(formReleaseDate)
	if timeDiff.Hours() < 730 {
		fmt.Println(artistAlbum.Name, artistAlbum.ReleaseDate)
	}

}

// GetArtistAlbums : Get all the albums of artists
func GetFollowedArtistAlbums(client *spotify.Client, artists []spotify.Artist) []*spotify.SimplifiedAlbumObject {
	var allAlbums = []*spotify.SimplifiedAlbumObject{}
	for _, a := range artists {
		result, err := client.GetArtistAlbums(a.ID)
		if err != nil {
			log.Fatal(err)
		}
		allAlbums = append(allAlbums, result...)
		//PrintArtistWithAlbums(a, result)
	}
	return allAlbums
}

func PrintArtistWithAlbums(a spotify.Artist, albums []*spotify.SimplifiedAlbumObject) {
	fmt.Println("Artist : " + a.Name)
	for i, a := range albums {
		fmt.Println(i, a.Name)
	}
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

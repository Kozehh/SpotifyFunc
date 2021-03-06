package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Kozehh/SpotifyFunc/spotify"
	"github.com/Kozehh/SpotifyFunc/spotify/models"
)

const (
	ScopeUserReadPrivate       = "user-read-private"
	ScopeUserFollowRead        = "user-follow-read"
	ScopeUserFollowModify      = "user-follow-modify"
	ScopePlaylistModifyPrivate = "playlist-modify-private"
	// RedirectURI is the OAuth reditect URI for the application
	RedirectURI = "http://localhost:8080/callback"
	layoutISO   = "2006-01-02"
)

// Instance of an error type, auth Client and the state
var (
	err     error
	auth    = spotify.NewAuthenticator(RedirectURI, ScopeUserReadPrivate, ScopeUserFollowRead, ScopeUserFollowModify, ScopePlaylistModifyPrivate)
	channel = make(chan *models.Client)
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

	latestReleasedAlbum := GetFollowedArtistsLatest(followedArtists, client)

	AddLatestReleasesToPlaylist(latestReleasedAlbum, client)
	//PrintFollowedArtists(artists)

}

func AddLatestReleasesToPlaylist(latestReleasedAlbum []*models.SimplifiedAlbumObject, c *models.Client) {
	var newReleasedTracks = []*models.Track{}
	for _, l := range latestReleasedAlbum {
		tracks, err := c.GetAlbumTracks(l.ID, 50)
		if err != nil {
			log.Fatal(err)
		}
		newReleasedTracks = append(newReleasedTracks, tracks...)
	}
	c.AddLatestToPlaylist(newReleasedTracks)
}

func GetFollowedArtistsLatest(followedArtists []models.Artist, client *models.Client) []*models.SimplifiedAlbumObject {
	var newReleases = []*models.SimplifiedAlbumObject{}

	artistsAlbums := GetFollowedArtistAlbums(client, followedArtists)
	for _, aa := range artistsAlbums {
		// Check if albums released less than a month ago
		if isNew := GetMonthyReleases(aa); isNew {
			// if so add to newReleases slice
			// TODO : Should have a bd or keep track of the slice containing the latest releases
			// so that if I run the code two days in a row I only get the ones that released in the past day
			newReleases = append(newReleases, aa)
		}
	}
	return newReleases
}

// GetMonthyReleases : Check if albums released less than a month ago
// TODO: Should implement a way to fetch a sclice (newReleases) to only return
// the ones that were not been fetch allready this month from another use of the function
func GetMonthyReleases(artistAlbum *models.SimplifiedAlbumObject) bool {
	formReleaseDate, _ := time.Parse(layoutISO, artistAlbum.ReleaseDate)
	timeDiff := time.Since(formReleaseDate)
	if timeDiff.Hours() < 730 {
		return true
	}
	return false
}

// GetArtistAlbums : Get all the albums of artists
func GetFollowedArtistAlbums(client *models.Client, artists []models.Artist) []*models.SimplifiedAlbumObject {
	var allAlbums = []*models.SimplifiedAlbumObject{}
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

func PrintArtistWithAlbums(a models.Artist, albums []*models.SimplifiedAlbumObject) {
	fmt.Println("Artist : " + a.Name)
	for i, a := range albums {
		fmt.Println(i, a.Name)
	}
}

func GetCurrentUser(client *models.Client) {
	// use the client to make calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as : ", user.DisplayName)
}

func GetFollowedArtists(client *models.Client) []models.Artist {
	var lastArtistID = ""
	var artistList = []models.Artist{}
	for {
		// Get the a list of followed artists
		artists, err := client.GetFollowedArtists(50, lastArtistID)
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

func PrintFollowedArtists(artists []models.Artist) {
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

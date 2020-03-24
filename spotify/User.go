package spotify

import (
	"net/url"
	"strconv"
	"strings"
)

// //////////////////////////////////////////////////////////////////////////// //
// --------------------------------  STRUCTS  -------------------------------- //
// ////////////////////////////////////////////////////////////////////////// //

// Cursor : Key used to find the next page of items
type Cursor struct {
	After string `json:"after"`
}

// User contains the basic, publicly available information about a Spotify user.
type User struct {
	// The name displayed on the user's profile.
	// Note: Spotify currently fails to populate
	// this field when querying for a playlist.
	DisplayName string `json:"display_name"`
	// Known public external URLs for the user.
	ExternalURLs map[string]string `json:"external_urls"`
	// Information about followers of the user.
	Followers Followers `json:"followers"`
	// A link to the Web API endpoint for this user.
	Endpoint string `json:"href"`
	// The Spotify user ID for the user.
	ID string `json:"id"`
	// The user's profile image.
	Images []Image `json:"images"`
	// The Spotify URI for the user.
	URI URI `json:"uri"`
}

// PrivateUser contains additional information about a user.
// This data is private and requires user authentication.
type PrivateUser struct {
	User
	// The country of the user, as set in the user's account profile.
	// An ISO 3166-1 alpha-2 country code.  This field is only available when the
	// current user has granted acess to the ScopeUserReadPrivate scope.
	Country string `json:"country"`
	// The user's email address, as entered by the user when creating their account.
	// Note: this email is UNVERIFIED - there is no proof that it actually
	// belongs to the user.  This field is only available when the current user
	// has granted access to the ScopeUserReadEmail scope.
	Email string `json:"email"`
	// The user's Spotify subscription level: "premium", "free", etc.
	// The subscription level "open" can be considered the same as "free".
	// This field is only available when the current user has granted access to
	// the ScopeUserReadPrivate scope.
	Product string `json:"product"`
	// The user's date of birth, in the format 'YYYY-MM-DD'.  You can use
	// the DateLayout constant to convert this to a time.Time value.
	// This field is only available when the current user has granted
	// access to the ScopeUserReadBirthdate scope.
	Birthdate string `json:"birthdate"`
}

// Artist : Contains information about an artist
type Artist struct {
	// Name of the artist
	Name string `json:"name"`
	// The popularity of the artist
	Popularity int `json:"popularity"`
	// The object type "artist"
	Type string `json:"type"`
	// Known public external URLs for the artist.
	ExternalURLs map[string]string `json:"external_urls"`
	// Information about followers of the artist.
	Followers Followers `json:"followers"`
	// A list of the genres the artist is associated with
	Genres []string `json:"genres"`
	// A link to the Web API endpoint for this artist.
	Endpoint string `json:"href"`
	// The Spotify ID for the artist
	ID string `json:"id"`
	// Images of the artist
	Images []Image `json:"images"`
	// The Spotify URI for the artist.
	URI URI `json:"uri"`
}

// CursorBasedObj : aka cursor-based paging object is a container for a set of objects
// In this case, it is used for a set of artists
type CursorBasedObj struct {
	// A link to the Web API endpoint returning the full result of the request
	Link string `json:"href"`

	// A link to the Web API endpoint returning the full result of the request
	Limit int `json:"limit"`

	// A link to the Web API endpoint returning the full result of the request
	Next string `json:"next"`

	// A link to the Web API endpoint returning the full result of the request
	Cursor Cursor `json:"cursors"`

	// A link to the Web API endpoint returning the full result of the request
	Total int `json:"total"`
}

// FullArtistCursorPage : Is the full object returned by the API Endpoint '/v1/me/following?type=artist'
type FullArtistCursorPage struct {
	CursorBasedObj
	Artists []Artist `json:"items"`
}

// SimplifiedAlbumObject : Is the full object returned by the API Endpoint '/v1/artists/{id}/albums'
type SimplifiedAlbumObject struct {
	// Compare to AlbumType this field represents relationship between the artist and the album
	AlbumGroup string `json:"album_group"`
	// The type of the album: one of “album”, “single”, or “compilation”
	AlbumType string `json:"album_type"`
	// The artists of the album
	Artists              []Artist `json:"artists"`
	ReleaseDate          string   `json:"release_date"`
	ReleaseDatePrecision string   `json:"release_date_precision"`
	// The markets in which the album is available
	AvailableMarkets []string `json:"available_markets"`
	// Known external URLs for this album
	ExternalURLs map[string]string `json:"external_urls"`
	// A link to the Web API endpoint providing full details of the album
	Endpoint string `json:"href"`
	// The Spotify ID for the album
	ID string `json:"id"`
	// The cover art for the album
	Images []Image `json:"images"`
	// The Spotify URI for the album.
	URI URI `json:"uri"`
	// The object type "album"
	Type string `json:"type"`
	// Name of the album
	Name string `json:"name"`
}

// ////////////////////////////////////////////////////////////////////////////// //
// --------------------------------  FUNCTIONS  -------------------------------- //
// //////////////////////////////////////////////////////////////////////////// //

// CurrentUser :
func (c *Client) CurrentUser() (*PrivateUser, error) {
	var result PrivateUser

	err := c.get(c.baseURL+"me", &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// FollowedList :
// Arg :
// 		(1) - The maximum number of items to return. Default: 20 / Min: 1 / Max: 50 | *Put -1 to use default
// 		(2) - The last artist ID retrieved from the previous request
// Return : A pointer of the FullArtistCursorPage object recieved from the endpoint call
func (c *Client) FollowedList(limit int, after string) (*FullArtistCursorPage, error) {
	funcURL := c.baseURL + "me/following"

	// Set query parameters
	v := url.Values{}
	v.Set("type", "artist")

	if limit != -1 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if after != "" {
		v.Set("after", after)
	}
	if params := v.Encode(); params != "" {
		funcURL += "?" + params
	}

	// The return value of the API Endpoint is a 'FillArtistCursorPage'
	// which it's key is 'artists'
	var result struct {
		A FullArtistCursorPage `json:"artists"`
	}

	err := c.get(funcURL, &result)
	if err != nil {
		return nil, err
	}

	return &result.A, nil
}

func (c *Client) GetArtistAlbums(id string) ([]*SimplifiedAlbumObject, error) {
	//funcURL := c.baseURL + "v1/artists/{id}/albums"

	// Set query parameters
	v := url.Values{}
	/*

		v.Set("include_groups", "album,single,appears_on")
		v.Set("limit", "")
		v.Set("offset", "")
	*/

	funcURL := c.baseURL + "artists/id/albums"
	funcURL = strings.Replace(funcURL, "id", id, -1)
	if params := v.Encode(); params != "" {
		funcURL += "?" + params
	}

	var a struct {
		Albums []*SimplifiedAlbumObject `json:"items"`
	}

	err := c.get(funcURL, &a)
	if err != nil {
		return nil, err
	}
	return a.Albums, nil
}

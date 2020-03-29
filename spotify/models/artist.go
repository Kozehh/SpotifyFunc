package models

import (
	"net/url"
	"strings"
)

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
	URI string `json:"uri"`
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

// ////////////////////////////////////////////////////////////////////////////// //
// --------------------------------  FUNCTIONS  -------------------------------- //
// //////////////////////////////////////////////////////////////////////////// //

// GetArtistAlbums :
func (c *Client) GetArtistAlbums(id string) ([]*SimplifiedAlbumObject, error) {
	// Set query parameters
	v := url.Values{}
	v.Set("include_groups", "album,single")

	funcURL := c.BaseURL + "artists/id/albums"
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

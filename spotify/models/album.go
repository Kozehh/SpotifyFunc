package models

import (
	"net/url"
	"strconv"
	"strings"
)

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
	URI string `json:"uri"`
	// The object type "album"
	Type string `json:"type"`
	// Name of the album
	Name string `json:"name"`
}

// ////////////////////////////////////////////////////////////////////////////// //
// --------------------------------  FUNCTIONS  -------------------------------- //
// //////////////////////////////////////////////////////////////////////////// //

func (c *Client) GetAlbumTracks(albumID string, limit int) ([]*Track, error) {
	funcURL := c.BaseURL + "albums/id/tracks"
	funcURL = strings.Replace(funcURL, "id", albumID, -1)

	// Set query parameters
	v := url.Values{}
	v.Set("type", "artist")

	if limit != -1 {
		v.Set("limit", strconv.Itoa(limit))
	}

	var res struct {
		Tracks []*Track `json:"items"`
	}

	err := c.get(funcURL, &res)
	if err != nil {
		return nil, err
	}
	return res.Tracks, nil
}

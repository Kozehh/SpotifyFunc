package models

type Track struct {
	Artists []Artist `json:"artists"`
	// The markets in which the album is available
	AvailableMarkets []string `json:"available_markets"`
	DiscNumber       int      `json:"disc_number"`
	Duration         int      `json:"duration_ms"`
	Explicit         bool     `json:"explicit"`
	// Known external URLs for this track
	ExternalURLs map[string]string `json:"external_urls"`
	// A link to the Web API endpoint providing full details of the album
	Endpoint string `json:"href"`
	// The Spotify ID for the album
	ID         string      `json:"id"`
	IsPlayable bool        `json:"is_playable"`
	LinkedFrom LinkedTrack `json:"linked_from"`
	Name       string      `json:"name"`
	PreviewURL string      `json:"preview_url"`
	TrackNum   int         `json:"track_number"`
	// The Spotify URI for the album.
	URI string `json:"uri"`
	// The object type "album"
	Type string `json:"type"`
}

type LinkedTrack struct {
	// Known external URLs for this track
	ExternalURLs map[string]string `json:"external_urls"`
	// A link to the Web API endpoint providing full details of the album
	Endpoint string `json:"href"`
	// The Spotify ID for the album
	ID string `json:"id"`
	// The Spotify URI for the album.
	URI string `json:"uri"`
	// The object type "album"
	Type string `json:"type"`
}

// ////////////////////////////////////////////////////////////////////////////// //
// --------------------------------  FUNCTIONS  -------------------------------- //
// //////////////////////////////////////////////////////////////////////////// //

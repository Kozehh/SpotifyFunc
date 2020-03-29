package models

import (
	"net/url"
	"strconv"
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
	URI string `json:"uri"`
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

// ////////////////////////////////////////////////////////////////////////////// //
// --------------------------------  FUNCTIONS  -------------------------------- //
// //////////////////////////////////////////////////////////////////////////// //

// CurrentUser :
func (c *Client) CurrentUser() (*PrivateUser, error) {
	var result PrivateUser

	err := c.get(c.BaseURL+"me", &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetFollowedArtists :
// Arg :
// 		(1) - The maximum number of items to return. Default: 20 / Min: 1 / Max: 50 | *Put -1 to use default
// 		(2) - The last artist ID retrieved from the previous request
// Return : A pointer of the FullArtistCursorPage object recieved from the endpoint call
func (c *Client) GetFollowedArtists(limit int, after string) (*FullArtistCursorPage, error) {
	funcURL := c.BaseURL + "me/following"

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

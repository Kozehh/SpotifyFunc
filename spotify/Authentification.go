package spotify

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"os"

	"github.com/Kozehh/SpotifyFunc/spotify/models"
	"golang.org/x/oauth2"
)

const (
	// AuthURL is the URL to Spotify Accounts Service's OAuth2 endpoint.
	AuthURL = "https://accounts.spotify.com/authorize"
	// TokenURL is the URL to the Spotify Accounts Service's OAuth2
	// token endpoint.
	TokenURL = "https://accounts.spotify.com/api/token"
)

// Authenticator is a struct containing a http context and OAuth2 configurations
// Config describes a typical 3-legged OAuth2 flow, with both the
// client application information and the server's endpoint URLs.
type Authenticator struct {
	config  *oauth2.Config
	context context.Context
}

// NewAuthenticator : Returns new spotify authentificator
func NewAuthenticator(redirectURL string, scopes ...string) Authenticator {
	cfg := &oauth2.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  AuthURL,
			TokenURL: TokenURL,
		},
	}

	// disable HTTP/2 for DefaultClient, see: https://github.com/zmb3/spotify/issues/20
	tr := &http.Transport{
		TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{},
	}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{Transport: tr})
	return Authenticator{
		config:  cfg,
		context: ctx,
	}
}

// AuthURL : Calls OAuth2 method 'AuthCodeURL' with the current authenticator configs
// Returns: A URL to OAuth 2.0 provider's consent page
// that asks for permissions for the required scopes explicitly.
func (a Authenticator) AuthURL(state string) string {
	return a.config.AuthCodeURL(state)
}

// Token : Maps the internal token in the authenticator context with the OAuth2 token
func (a Authenticator) Token(state string, req *http.Request) (*oauth2.Token, error) {
	values := req.URL.Query()
	if e := values.Get("error"); e != "" {
		return nil, errors.New("spotify: auth failed - " + e)
	}
	code := values.Get("code")
	if code == "" {
		return nil, errors.New("spotify: didn't get access code")
	}
	actualState := values.Get("state")
	if actualState != state {
		return nil, errors.New("spotify: redirect state parameter doesn't match")
	}

	// if there was no errors or mismatches, (link/refresh/map) the internal token with the OAuth Token
	return a.config.Exchange(a.context, code)
}

// NewClient : Creates a Client that will use the specified access token for its API requests.
func (a Authenticator) NewClient(token *oauth2.Token) models.Client {
	// Create a new http client using the token and current context
	client := a.config.Client(a.context, token)
	// The app client object is now the new one created
	return models.Client{
		Http:    client,
		BaseURL: models.BaseAddress,
	}
}

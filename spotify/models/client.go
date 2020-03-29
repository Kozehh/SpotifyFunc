package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

///////////////////////////// ********* CONSTANTS ********* /////////////////////////////////

const (
	// DateLayout can be used with time.Parse to create time.Time values
	// from Spotify date strings.  For example, PrivateUser.Birthdate
	// uses this format.
	DateLayout = "2006-01-02"

	// TimestampLayout can be used with time.Parse to create time.Time
	// values from SpotifyTimestamp strings.  It is an ISO 8601 UTC timestamp
	// with a zero offset.  For example, PlaylistTrack's AddedAt field uses
	// this format.
	TimestampLayout = "2006-01-02T15:04:05Z"

	// defaultRetryDurationS helps us fix an apparent server bug whereby we will
	// be told to retry but not be given a wait-interval.
	defaultRetryDuration = time.Second * 5

	// rateLimitExceededStatusCode is the code that the server returns when our
	// request frequency is too high.
	rateLimitExceededStatusCode = 429

	// BaseAddress is the base address for all the the spotify API Endpoints
	BaseAddress = "https://api.spotify.com/v1/"
)

////////////////////////////////////////////////////////////////////////////////////////////

///////////////////////////// ********* STRUCTS ********* /////////////////////////////////

// Client is a client for working with the Spotify Web API.
// To create an authenticated client, use the `Authenticator.NewClient` method.
type Client struct {
	Http    *http.Client
	BaseURL string

	AutoRetry bool
}

// Error : Represents an error returned by the Spotify Web API.
type Error struct {
	// A short description of the error.
	Message string `json:"message"`
	// The HTTP status code.
	Status int `json:"status"`
}

////////////////////////////////////////////////////////////////////////////////////////////

/////////////////////////// ********* FUNCTIONS ********* ///////////////////////////////

// Return the response
func (c *Client) get(url string, result interface{}) error {
	for {
		resp, err := c.Http.Get(url)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode == rateLimitExceededStatusCode && c.AutoRetry {
			time.Sleep(retryDuration(resp))
			continue
		}
		if resp.StatusCode == http.StatusNoContent {
			return nil
		}
		if resp.StatusCode != http.StatusOK {
			return c.decodeError(resp)
		}

		err = json.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			return err
		}

		break
	}

	return nil
}

// `execute` executes a non-GET request. `needsStatus` describes other HTTP
// status codes that will be treated as success. Note that we allow all 200s
// even if there are additional success codes that represent success.
func (c *Client) execute(req *http.Request, result interface{}, needsStatus ...int) error {
	for {
		resp, err := c.Http.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if c.AutoRetry && shouldRetry(resp.StatusCode) {
			time.Sleep(retryDuration(resp))
			continue
		}
		if resp.StatusCode == http.StatusNoContent {
			return nil
		}
		if (resp.StatusCode >= 300 ||
			resp.StatusCode < 200) &&
			isFailure(resp.StatusCode, needsStatus) {
			return c.decodeError(resp)
		}

		if result != nil {
			if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
				return err
			}
		}
		break
	}
	return nil
}

// decodeError decodes an Error from an io.Reader.
func (c *Client) decodeError(resp *http.Response) error {
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(responseBody) == 0 {
		return fmt.Errorf("spotify: HTTP %d: %s (body empty)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	buf := bytes.NewBuffer(responseBody)

	var e struct {
		E Error `json:"error"`
	}
	err = json.NewDecoder(buf).Decode(&e)
	if err != nil {
		return fmt.Errorf("spotify: couldn't decode error: (%d) [%s]", len(responseBody), responseBody)
	}

	if e.E.Message == "" {
		// Some errors will result in there being a useful status-code but an
		// empty message, which will confuse the user (who only has access to
		// the message and not the code). An example of this is when we send
		// some of the arguments directly in the HTTP query and the URL ends-up
		// being too long.

		e.E.Message = fmt.Sprintf("spotify: unexpected HTTP %d: %s (empty error)",
			resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return e.E
}

// shouldRetry determines whether the status code indicates that the
// previous operation should be retried at a later time
func shouldRetry(status int) bool {
	return status == http.StatusAccepted || status == http.StatusTooManyRequests
}

// isFailure determines whether the code indicates failure
func isFailure(code int, validCodes []int) bool {
	for _, item := range validCodes {
		if item == code {
			return false
		}
	}
	return true
}

func retryDuration(resp *http.Response) time.Duration {
	raw := resp.Header.Get("Retry-After")
	if raw == "" {
		return defaultRetryDuration
	}
	seconds, err := strconv.ParseInt(raw, 10, 32)
	if err != nil {
		return defaultRetryDuration
	}
	return time.Duration(seconds) * time.Second
}

func (e Error) Error() string {
	return e.Message
}

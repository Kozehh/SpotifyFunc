package models

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

const _playlistNewNewID = "7bs6NLtazYeLmf7uqYs9he"

// ////////////////////////////////////////////////////////////////////////////// //
// --------------------------------  FUNCTIONS  -------------------------------- //
// //////////////////////////////////////////////////////////////////////////// //

func (c *Client) AddLatestToPlaylist(tracks []*Track) {
	funcURL := c.BaseURL + "playlists/{playlist_id}/tracks"
	funcURL = strings.Replace(funcURL, "{playlist_id}", _playlistNewNewID, -1)

	var reqTracks = []string{}
	for i := 1; i <= len(tracks); i++ {
		reqTracks = append(reqTracks, tracks[i-1].URI)
		if (i%100) == 0 || i == (len(tracks)) {
			AddToPlaylist(reqTracks, funcURL, c)
			reqTracks = nil
		}
	}

}

func AddToPlaylist(tracks []string, funcURL string, c *Client) {
	m := make(map[string]interface{})
	m["uris"] = tracks
	body, err := json.Marshal(m)
	os.Stdout.Write(body)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", funcURL, bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	result := struct {
		SnapshotID string `json:"snapshot_id"`
	}{}

	err = c.execute(req, &result, http.StatusCreated)
	if err != nil {
		log.Fatal(err)
	}
}

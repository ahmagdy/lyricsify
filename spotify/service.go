package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	config "github.com/Ahmad-Magdy/lyricsify/config"

	"github.com/Ahmad-Magdy/lyricsify/types"
)

// SpotifyService Service to communicate with spotify
type SpotifyService struct {
	spotifyAPIUrl string
	config        *config.Config
}

// New create a new instance of SpotifyService
func New(config *config.Config) *SpotifyService {
	return &SpotifyService{"https://api.spotify.com/v1/", config}
}

// getSongsList To get Me Songs list
func (spotifyService *SpotifyService) getSongsList(ctx context.Context, reqURL string) (response types.MeTrackResponse, err error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return types.MeTrackResponse{}, err
	}
	spotifyToken := spotifyService.config.SpotifyToken
	if spotifyToken == "" {
		return types.MeTrackResponse{}, errors.New("spotify token is not set")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", spotifyToken))
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return types.MeTrackResponse{}, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return types.MeTrackResponse{}, err
	}

	if res.StatusCode != 200 {
		return types.MeTrackResponse{}, fmt.Errorf("Request with URL %v exit with code %v and text %v", res.Request.URL, res.StatusCode, string(body))
	}
	trackResponse := types.MeTrackResponse{}
	json.Unmarshal(body, &trackResponse)
	return trackResponse, nil
}

// GetAllSongs to get all liked songs from spotify Me list return a map of string and string, the key is the song name and the value is the artists name
func (spotifyService *SpotifyService) GetAllLikedSongs(ctx context.Context) (map[string]string, error) {
	songs := make(map[string]string)
	reqURL := fmt.Sprintf("%vme/tracks", spotifyService.spotifyAPIUrl)
	for {
		anon, err := spotifyService.getSongsList(ctx, reqURL)
		if err != nil {
			return nil, err
		}
		for _, y := range anon.Items {
			songs[y.Track.Name] = spotifyService.getArtistsName(y.Track.Artists)
		}
		if len(anon.Next) == 0 {
			break
		}
		reqURL = anon.Next
	}
	return songs, nil
}

func (spotifyService *SpotifyService) getArtistsName(artistList []types.Artist) string {
	var artistsName []string
	for _, item := range artistList {
		artistsName = append(artistsName, item.Name)
	}
	return strings.Join(artistsName, ",")
}

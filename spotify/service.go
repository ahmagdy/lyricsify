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

const _spotifyBaseURL = "https://api.spotify.com/v1"

var ErrSpotifyTokenNotSet = errors.New("spotify token is not set")

// Service Service to communicate with spotify
type Service struct {
	spotifyAPIUrl string
	config        *config.Config
}

// New create a new instance of Service
func New(config *config.Config) *Service {
	return &Service{_spotifyBaseURL, config}
}

// GetAllSongs to get all liked songs from spotify Me list return a map of string and string, the key is the song name and the value is the artists name
func (s *Service) GetAllLikedSongs(ctx context.Context) (map[string]string, error) {
	songs := make(map[string]string)
	reqURL := fmt.Sprintf("%v/me/tracks", s.spotifyAPIUrl)
	for {
		trackResponse, err := s.songsList(ctx, reqURL)
		if err != nil {
			return nil, err
		}

		for _, trackRes := range trackResponse.Items {
			songs[trackRes.Track.Name] = s.artistsName(trackRes.Track.Artists)
		}
		if len(trackResponse.Next) == 0 {
			break
		}

		reqURL = trackResponse.Next
	}

	return songs, nil
}

// songsList To get Me Songs list
func (s *Service) songsList(ctx context.Context, reqURL string) (response types.MeTrackResponse, err error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return types.MeTrackResponse{}, err
	}

	spotifyToken := s.config.SpotifyToken
	if spotifyToken == "" {
		return types.MeTrackResponse{}, ErrSpotifyTokenNotSet
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
		return types.MeTrackResponse{}, fmt.Errorf("request with URL %v exit with code %v and text %v", res.Request.URL, res.StatusCode, string(body))
	}
	trackResponse := types.MeTrackResponse{}
	json.Unmarshal(body, &trackResponse)
	return trackResponse, nil
}

func (s *Service) artistsName(artistList []types.Artist) string {
	var artistsName []string
	for _, item := range artistList {
		artistsName = append(artistsName, item.Name)
	}
	return strings.Join(artistsName, ",")
}

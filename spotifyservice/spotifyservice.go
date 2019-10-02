package spotifyservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/Ahmad-Magdy/lyricsify/models"
)

// SpotifyService Service to communicate with spotify
type SpotifyService struct {
	spotifyAPIUrl string
}

// New create a new instance of SpotifyService
func New() *SpotifyService {
	return &SpotifyService{"https://api.spotify.com/v1/"}
}

// getSongsList To get Me Songs list
func (spotifyService *SpotifyService) getSongsList(ctx context.Context, reqURL string) (response models.MeTrackResponse, err error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return models.MeTrackResponse{}, err
	}
	spotifyToken := os.Getenv("SPOTIFY_TOKEN")
	if spotifyToken == "" {
		return models.MeTrackResponse{}, errors.New("SPOTIFY_TOKEN environment variable is not found.")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", spotifyToken))
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.MeTrackResponse{}, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return models.MeTrackResponse{}, err
	}

	if res.StatusCode != 200 {
		return models.MeTrackResponse{}, fmt.Errorf("Request with URL %v exit with code %v and text %v", res.Request.URL, res.StatusCode, string(body))
	}
	trackResponse := models.MeTrackResponse{}
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

func (spotifyService *SpotifyService) getArtistsName(artistList []models.Artist) string {
	var artistsName []string
	for _, item := range artistList {
		artistsName = append(artistsName, item.Name)
	}
	return strings.Join(artistsName, ",")
}

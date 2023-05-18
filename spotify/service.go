package spotify

import (
	"context"
	"errors"
	"fmt"
	"strings"

	config "github.com/ahmagdy/lyricsify/config"
	"github.com/zmb3/spotify/v2"
)

var ErrSpotifyTokenNotSet = errors.New("spotify token is not set")

// Service Service to communicate with spotify
type Service struct {
	spotifyClient *spotify.Client
	config        *config.Config
}

// New create a new instance of Service
// TODO: replace client with client interface
func New(config *config.Config, spotifyClient *spotify.Client) *Service {
	return &Service{spotifyClient, config}
}

// AllLikedSongs to get all liked songs from spotify Me list return a map of string and string, the key is the song name and the value is the artists name
func (s *Service) AllLikedSongs(ctx context.Context) (map[string]string, error) {
	songToArtists := make(map[string]string)

	userTracks, err := s.spotifyClient.CurrentUsersTracks(ctx)
	if err != nil {
		return nil, fmt.Errorf("spotifyClient.CurrentUsersTracks: %w", err)
	}

	for err == nil {
		for _, track := range userTracks.Tracks {
			songToArtists[track.Name] = s.artistsName(track.Artists)
		}

		err = s.spotifyClient.NextPage(ctx, userTracks)
	}

	if !errors.Is(err, spotify.ErrNoMorePages) {
		return nil, fmt.Errorf("spotifyClient.NextPage: %w", err)
	}

	return songToArtists, nil
}

func (s *Service) artistsName(artists []spotify.SimpleArtist) string {
	var artistsName []string
	for _, artist := range artists {
		artistsName = append(artistsName, artist.Name)
	}

	return strings.Join(artistsName, ",")
}

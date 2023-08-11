package scraper

import (
	"context"
	"fmt"
	"net/http"

	config "github.com/ahmagdy/lyricsify/config"
	mxm "github.com/milindmadhukar/go-musixmatch"
	"github.com/milindmadhukar/go-musixmatch/params"
	"go.uber.org/zap"
)

type muxixmatchService struct {
	logger *zap.Logger
	client *mxm.Client
}

func newMuxixmatchService(logger *zap.Logger, config *config.Config) Service {
	client := mxm.New(config.MusixmatchToken, http.DefaultClient)

	return &muxixmatchService{
		client: client,
		logger: logger,
	}
}

func (s *muxixmatchService) FindLyrics(ctx context.Context, songName string, artists string) (string, error) {
	// TODO: Consider using (QueryTrackArtist: Any word in the song title or artist name) as well
	tracks, err := s.client.SearchTrack(ctx, params.QueryTrack(songName), params.QueryArtist(artists), params.HasLyrics(true))
	if err != nil {
		return "", fmt.Errorf("client.SearchTrack: %w", err)
	}
	if len(tracks) == 0 {
		return "", fmt.Errorf("no tracks found for song (%v)", songName)
	}

	for _, track := range tracks {
		lyrics, err := s.client.GetTrackLyrics(ctx, params.TrackID(track.ID))
		if err != nil {
			return "", fmt.Errorf("client.GetLyrics: %w", err)
		}

		if lyrics.Body != "" {
			return lyrics.Body, nil
		}

		s.logger.Warn("No lyrics found for track", zap.String("track", track.Name), zap.Int("track_id", track.ID))
	}

	return "", fmt.Errorf("could not extrac lyrics for the songe (%v)", songName)
}

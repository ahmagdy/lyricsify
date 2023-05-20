package main

import (
	"context"
	"log"
	"time"

	"github.com/ahmagdy/lyricsify"
	"github.com/ahmagdy/lyricsify/config"
	"github.com/sourcegraph/conc/pool"
	"go.uber.org/zap"
)

type contextKey string

const (
	_ctxTimeout = 2 * time.Minute
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), _ctxTimeout)
	defer cancel()

	logger, _ := zap.NewProduction()
	cfg, _ := config.New()
	authServer := lyricsify.NewAuthServer(logger, cfg)
	authServer.Start()
	authServer.WaitForAuthToBeCompleted()

	if err := authServer.Verify(ctx); err != nil {
		// TODO: handle
	}

	svc, err := lyricsify.New(ctx, authServer.SpotifyClient())
	if err != nil {
		log.Fatal(err)
	}
	if err := loadSongs(ctx, svc, logger); err != nil {
		log.Fatal(err)
	}
}

func loadSongs(ctx context.Context, s *lyricsify.Service, logger *zap.Logger) error {
	songToArtists, err := s.LoadSongs(ctx)
	if err != nil {
		return err
	}

	p := pool.New().WithErrors()

	for song, artists := range songToArtists {
		song := song
		artists := artists
		p.Go(func() error {
			ctxKey := contextKey(song)
			ctx := context.WithValue(ctx, ctxKey, song)

			logger.Info("Checking song", zap.String("songName", song), zap.String("artists", artists))

			isExist, err := s.HasLyrics(ctx, song)
			if err != nil {
				return err
			}

			if isExist {
				logger.Info("Skipping song as it already exist in the datastore", zap.String("songName", song))
				return nil
			}

			lyrics, err := s.Fetch(ctx, song, artists)
			if err != nil {
				return err
			}

			return s.Save(ctx, song, lyrics)
		})
	}

	if err := p.Wait(); err != nil {
		// fail open for now
		// return err
		logger.Error("received the following error when loading titles", zap.Error(err))
	}

	logger.Info("done")

	searchResults, err := s.Search(ctx, "Did you work real hard")
	if err != nil {
		return err
	}

	for _, result := range searchResults {
		logger.Info("search results", zap.String("candidateTitle", result.Title))
	}

	return nil
}

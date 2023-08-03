package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
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
	logger, _ := zap.NewProduction()

	if err := run(logger); err != nil {
		logger.Fatal("failed to run the example", zap.Error(err))
	}
}
func run(logger *zap.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), _ctxTimeout)
	defer cancel()

	cfg, _ := config.New()
	authServer := lyricsify.NewAuthServer(logger, cfg)
	logger.Info("starting the server")
	if err := authServer.Start(); err != nil {
		logger.Error("failed to run", zap.Error(err))
	}

	authServer.WaitForAuthToBeCompleted()

	if err := authServer.Verify(ctx); err != nil {
		return fmt.Errorf("failed to verify: %w", err)
	}

	svc, err := lyricsify.New(ctx, authServer.SpotifyClient())
	if err != nil {
		return fmt.Errorf("failed to inti svc: %w", err)
	}

	if err := loadSongs(ctx, svc, logger); err != nil {
		return fmt.Errorf("failed to load songs: %w", err)
	}

	logger.Info("done")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (kill -2)
	<-stop

	authServer.Stop(ctx)

	return nil
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

	return nil
}

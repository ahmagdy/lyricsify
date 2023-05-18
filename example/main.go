package main

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/ahmagdy/lyricsify"
	"github.com/ahmagdy/lyricsify/config"
	"github.com/ahmagdy/lyricsify/spotify"
	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"
)

type contextKey string

var wg sync.WaitGroup

const (
	_ctxTimeout = 2 * time.Minute
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), _ctxTimeout)
	defer cancel()

	logger, _ := zap.NewProduction()
	cfg, _ := config.New()
	authServer := spotify.NewAuthServer(logger, cfg)
	authServer.Start()
	authServer.WaitForAuthToBeCompleted()

	if err := authServer.Verify(ctx); err != nil {
		// TODO: handle
	}

	svc, err := lyricsify.New(ctx, authServer.SpotifyClient())
	if err != nil {
		log.Fatal(err)
	}
	if err := loadSongs(ctx, svc); err != nil {
		log.Fatal(err)
	}
}

func loadSongs(ctx context.Context, l *lyricsify.Service) error {
	songsMap, err := l.LoadSongs(ctx)
	if err != nil {
		return err
	}

	var combinedErr error

	for song, artists := range songsMap {
		wg.Add(1)
		go func(song string, artists string) {
			defer wg.Done()

			ctxKey := contextKey(song)
			ctx := context.WithValue(ctx, ctxKey, song)

			log.Println(song, artists)

			isExist, err := l.HasLyrics(ctx, song)
			if err != nil {
				combinedErr = multierror.Append(combinedErr, err)
				return
			}
			if isExist {
				log.Printf("Skipping song %v since it's already exist in the datastore", song)
				return
			}

			lyrics, err := l.Fetch(ctx, song, artists)
			if err != nil {
				combinedErr = multierror.Append(combinedErr, err)
				return
			}

			err = l.Save(ctx, song, lyrics)
			combinedErr = multierror.Append(combinedErr, err)
		}(song, artists)
	}

	wg.Wait()
	if combinedErr != nil {
		return combinedErr
	}
	log.Println(strings.Repeat("-", 5), "Done! ", strings.Repeat("-", 5))

	searchResults, err := l.Search(ctx, "Did you work real hard")
	if err != nil {
		return err
	}

	for _, result := range searchResults {
		log.Println(result.Title)
	}

	return nil
}

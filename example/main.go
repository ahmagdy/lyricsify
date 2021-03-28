package main

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Ahmad-Magdy/lyricsify"
	"github.com/hashicorp/go-multierror"
)

type contextKey string

var wg sync.WaitGroup

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	svc, err := lyricsify.InitializeLyricsify(ctx)
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

			lyrics, err := l.FetchLyrics(ctx, song, artists)
			if err != nil {
				combinedErr = multierror.Append(combinedErr, err)
				return
			}

			err = l.SaveLyrics(ctx, song, lyrics)
			combinedErr = multierror.Append(combinedErr, err)
		}(song, artists)
	}

	wg.Wait()
	if combinedErr != nil {
		return combinedErr
	}
	log.Println(strings.Repeat("-", 5), "Done! ", strings.Repeat("-", 5))

	searchResults, err := l.SearchByText(ctx, "Did you work real hard")
	if err != nil {
		return err
	}

	for _, result := range searchResults {
		log.Println(result.Title)
	}

	return nil
}

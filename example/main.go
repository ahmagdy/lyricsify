package main

import (
	"context"
	"github.com/Ahmad-Magdy/lyricsify/errorhandler"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Ahmad-Magdy/lyricsify"
)

type contextKey string

var wg sync.WaitGroup

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	lyricsify := lyricsify.InitializeLyricsify(ctx)

	songsMap, err := lyricsify.LoadSongs(ctx)
	errorhandler.HandleError("Main LoadSongs: ", err)

	for song, artists := range songsMap {
		wg.Add(1)
		go func(song string, artists string) {
			defer wg.Done()
			ctxKey := contextKey(song)
			ctx := context.WithValue(ctx, ctxKey, song)

			log.Println(song, artists)


			isExist, err := lyricsify.IsLyricsExist(ctx, song)
			errorhandler.HandleError("Main IsLyricsExist: ", err)
			if isExist {
				log.Printf("Skipping song %v since it's already exist in the datastore", song)
				return
			}
			lyrics, err := lyricsify.FetchLyrics(ctx, song, artists)
			if err != nil {
				log.Println(err.Error())
				return
			}
			err = lyricsify.SaveLyrics(ctx, song, lyrics)
			errorhandler.HandleError("Main SaveLyrics: ", err)
		}(song, artists)
	}

	wg.Wait()
	log.Println(strings.Repeat("-", 5), "Done! ", strings.Repeat("-", 5))

	searchResults, err := lyricsify.SearchByText(ctx, "Did you work real hard")
	errorhandler.HandleError("search here ",err)
	for _, result := range searchResults {
		log.Println(result.Title)
	}
}

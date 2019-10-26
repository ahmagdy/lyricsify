package main

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Ahmad-Magdy/lyricsify"
	"github.com/Ahmad-Magdy/lyricsify/errorhandler"
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
			log.Println(song, artists)
			defer wg.Done()
			ctxKey := contextKey(song)
			ctx := context.WithValue(ctx, ctxKey, song)
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

}

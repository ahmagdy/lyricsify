package lyricsify

import (
	"context"

	"github.com/Ahmad-Magdy/lyricsify/search"
	lyricsscraping "github.com/Ahmad-Magdy/lyricsify/scraping"
	"github.com/Ahmad-Magdy/lyricsify/spotify"
)

// Lyricsify Main package service
type Lyricsify struct {
	spotifyService *spotify.SpotifyService
	scraper        *lyricsscraping.LyricsScrapingService
	elasticClient  *search.LyricsSearchService
}

// New To create a new instance of Lyricsify
func New(spotifyService *spotify.SpotifyService, scraper *lyricsscraping.LyricsScrapingService, elasticClient *search.LyricsSearchService) *Lyricsify {
	return &Lyricsify{spotifyService, scraper, elasticClient}
}

// LoadSongs To load all songs from "LikedSongs" section in spotify
func (lyricsService *Lyricsify) LoadSongs(ctx context.Context) (songsMap map[string]string, err error) {
	allSongs, err := lyricsService.spotifyService.GetAllLikedSongs(ctx)
	return allSongs, err
}

// FetchLyrics To fetch song lyrics from the scraper
func (lyricsService *Lyricsify) FetchLyrics(ctx context.Context, songName string, artists string) (lyrics string, err error) {
	lyricsContent, err := lyricsService.scraper.GetLyricsForSong(ctx, songName, artists)
	return lyricsContent, err
}

// SaveLyrics to save lyrics in a datastore in this case elasticsearch
func (lyricsService *Lyricsify) SaveLyrics(ctx context.Context, title string, lyrics string) (err error) {
	err = lyricsService.elasticClient.Create(ctx, title, lyrics)
	return err
}

// SearchByText to search in the saved songs by text
func (lyricsService *Lyricsify) SearchByText(ctx context.Context, text string) (res []search.LyricsBody, err error) {
	results, err := lyricsService.elasticClient.Search(ctx, text)
	return results, err
}

func(lyricsService *Lyricsify) IsLyricsExist(ctx context.Context, title string) (bool, error){
	id, err:= 	lyricsService.elasticClient.GetItemID(ctx, title)
	if err != nil{
		return false, err
	}
	return id != "", nil
}
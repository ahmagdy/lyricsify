package lyricsify

import (
	"context"

	"github.com/Ahmad-Magdy/lyricsify/elasticclient"
	lyricsscraping "github.com/Ahmad-Magdy/lyricsify/scraping"
	"github.com/Ahmad-Magdy/lyricsify/spotifyservice"
)

type Lyricsify struct {
	spotifyService *spotifyservice.SpotifyService
	scraper        *lyricsscraping.LyricsScrapingService
	elasticClient  *elasticclient.LyricsSearchService
}

func NewLyricsifyService(spotifyService *spotifyservice.SpotifyService, scraper *lyricsscraping.LyricsScrapingService, elasticClient *elasticclient.LyricsSearchService) *Lyricsify {
	return &Lyricsify{spotifyService, scraper, elasticClient}
}

func (lyricsService *Lyricsify) LoadSongs(ctx context.Context) (songsMap map[string]string, err error) {
	allSongs, err := lyricsService.spotifyService.GetAllLikedSongs(ctx)
	return allSongs, err
}

func (lyricsService *Lyricsify) FetchLyrics(ctx context.Context, songName string, artists string) (lyrics string, err error) {
	lyricsContent, err := lyricsService.scraper.GetLyricsForSong(ctx, songName, artists)
	return lyricsContent, err
}

func (lyricsService *Lyricsify) SaveLyrics(ctx context.Context, title string, lyrics string) (err error) {
	err = lyricsService.elasticClient.Create(ctx, title, lyrics)
	return err
}

func (lyricsService *Lyricsify) SearchByText(ctx context.Context, text string) (res []elasticclient.LyricsBody, err error) {
	results, err := lyricsService.elasticClient.Search(ctx, text)
	return results, err
}

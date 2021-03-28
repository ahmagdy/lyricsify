package lyricsify

import (
	"context"

	scraper "github.com/Ahmad-Magdy/lyricsify/scraper"
	"github.com/Ahmad-Magdy/lyricsify/search"
	"github.com/Ahmad-Magdy/lyricsify/spotify"
)

// Service Main package service
type Service struct {
	spotifyService *spotify.Service
	scraper        *scraper.Service
	elasticClient  *search.Service
}

// New To create a new instance of Service
func New(spotifyService *spotify.Service, scraper *scraper.Service, elasticClient *search.Service) *Service {
	return &Service{spotifyService, scraper, elasticClient}
}

// LoadSongs To load all songs from "LikedSongs" section in spotify
func (s *Service) LoadSongs(ctx context.Context) (songsMap map[string]string, err error) {
	allSongs, err := s.spotifyService.AllLikedSongs(ctx)
	return allSongs, err
}

// FetchLyrics To fetch song lyrics from the scraper
func (s *Service) FetchLyrics(ctx context.Context, songName string, artists string) (lyrics string, err error) {
	lyricsContent, err := s.scraper.Lyrics(ctx, songName, artists)
	return lyricsContent, err
}

// SaveLyrics to save lyrics in a datastore in this case elasticsearch
func (s *Service) SaveLyrics(ctx context.Context, title string, lyrics string) (err error) {
	err = s.elasticClient.Create(ctx, title, lyrics)
	return err
}

// SearchByText to search in the saved songs by text
func (s *Service) SearchByText(ctx context.Context, text string) (res []search.LyricsBody, err error) {
	results, err := s.elasticClient.Search(ctx, text)
	return results, err
}

func (s *Service) HasLyrics(ctx context.Context, title string) (bool, error) {
	id, err := s.elasticClient.GetItemID(ctx, title)
	if err != nil {
		return false, err
	}
	return id != "", nil
}

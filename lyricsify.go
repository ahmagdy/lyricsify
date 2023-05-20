package lyricsify

import (
	"context"

	scraper "github.com/ahmagdy/lyricsify/internal/scraper"
	"github.com/ahmagdy/lyricsify/internal/search"
	"github.com/ahmagdy/lyricsify/internal/spotify"
)

// Service Main package service
type Service struct {
	spotifyService *spotify.Service
	scraper        *scraper.Service
	elasticClient  *search.Service
}

// New To create a new instance of Service
func new(spotifyService *spotify.Service, scraper *scraper.Service, elasticClient *search.Service) *Service {
	return &Service{spotifyService, scraper, elasticClient}
}

// LoadSongs To load all songs from "LikedSongs" section in spotify
func (s *Service) LoadSongs(ctx context.Context) (map[string]string, error) {
	allSongs, err := s.spotifyService.AllLikedSongs(ctx)
	return allSongs, err
}

// Fetch To fetch song lyrics from the scraper
func (s *Service) Fetch(ctx context.Context, songName string, artists string) (string, error) {
	lyricsContent, err := s.scraper.FindLyrics(ctx, songName, artists)
	return lyricsContent, err
}

// Save to save lyrics in a datastore in this case elasticsearch
func (s *Service) Save(ctx context.Context, title string, lyrics string) error {
	err := s.elasticClient.Create(ctx, title, lyrics)
	return err
}

// Search to search in the saved songs by text
func (s *Service) Search(ctx context.Context, text string) ([]search.LyricsBody, error) {
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

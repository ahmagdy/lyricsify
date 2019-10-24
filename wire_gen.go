// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package lyricsify

import (
	"context"
	"github.com/Ahmad-Magdy/lyricsify/elasticclient"
	"github.com/Ahmad-Magdy/lyricsify/internal"
	"github.com/Ahmad-Magdy/lyricsify/scraping"
	"github.com/Ahmad-Magdy/lyricsify/spotifyservice"
)

// Injectors from container.go:

func CreateSomething(ctx context.Context) *Lyricsify {
	spotifyService := spotifyservice.New()
	lyricsScrapingService := scrapping.New()
	configConfig := config.NewConfig()
	lyricsSearchService := elasticclient.New(ctx, configConfig)
	lyricsify := NewLyricsifyService(spotifyService, lyricsScrapingService, lyricsSearchService)
	return lyricsify
}
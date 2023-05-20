// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package lyricsify

import (
	"context"
	"github.com/ahmagdy/lyricsify/config"
	"github.com/ahmagdy/lyricsify/scraper"
	"github.com/ahmagdy/lyricsify/search"
	spotify2 "github.com/ahmagdy/lyricsify/spotify"
	"github.com/olivere/elastic/v7"
	"github.com/zmb3/spotify/v2"
	"go.uber.org/zap"
)

// Injectors from container.go:

func New(ctx context.Context, spotifyClient *spotify.Client) (*Service, error) {
	configConfig, err := config.New()
	if err != nil {
		return nil, err
	}
	service := spotify2.New(configConfig, spotifyClient)
	logger, err := createLogger()
	if err != nil {
		return nil, err
	}
	scraperService, err := scraper.New(configConfig, logger)
	if err != nil {
		return nil, err
	}
	client, err := createElasticClient()
	if err != nil {
		return nil, err
	}
	searchService, err := search.New(ctx, configConfig, client, logger)
	if err != nil {
		return nil, err
	}
	lyricsifyService := new(service, scraperService, searchService)
	return lyricsifyService, nil
}

// container.go:

func createLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}

func createElasticClient() (*elastic.Client, error) {
	return elastic.NewClient(elastic.SetSniff(false), elastic.SetRetrier(elastic.NewBackoffRetrier(elastic.NewSimpleBackoff(100, 200))))
}

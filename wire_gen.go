// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package lyricsify

import (
	"context"
	"github.com/Ahmad-Magdy/lyricsify/config"
	"github.com/Ahmad-Magdy/lyricsify/scraper"
	"github.com/Ahmad-Magdy/lyricsify/search"
	"github.com/Ahmad-Magdy/lyricsify/spotify"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

// Injectors from container.go:

func InitializeLyricsify(ctx context.Context) (*Service, error) {
	configConfig, err := config.New()
	if err != nil {
		return nil, err
	}
	service := spotify.New(configConfig)
	logger, err := createLogger()
	if err != nil {
		return nil, err
	}
	scraperService := scraper.New(configConfig, logger)
	client, err := createElasticClient()
	if err != nil {
		return nil, err
	}
	searchService, err := search.New(ctx, configConfig, client, logger)
	if err != nil {
		return nil, err
	}
	lyricsifyService := New(service, scraperService, searchService)
	return lyricsifyService, nil
}

// container.go:

func createLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}

func createElasticClient() (*elastic.Client, error) {
	return elastic.NewClient(elastic.SetSniff(false), elastic.SetRetrier(elastic.NewBackoffRetrier(elastic.NewSimpleBackoff(100, 200))))
}

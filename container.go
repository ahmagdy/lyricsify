//go:build wireinject
// +build wireinject

package main

import (
	"context"

	config "github.com/ahmagdy/lyricsify/config"
	scrapping "github.com/ahmagdy/lyricsify/scraper"
	"github.com/ahmagdy/lyricsify/search"
	"github.com/ahmagdy/lyricsify/spotify"
	"github.com/google/wire"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

func New(ctx context.Context) (*Service, error) {
	wire.Build(
		createLogger,
		createElasticClient,
		config.New,
		spotify.New,
		scrapping.New,
		search.New,
		new)
	return &Service{}, nil
}

func createLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}

func createElasticClient() (*elastic.Client, error) {
	return elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetRetrier(elastic.NewBackoffRetrier(elastic.NewSimpleBackoff(100, 200))),
	)
}

//go:build wireinject
// +build wireinject

package lyricsify

import (
	"context"

	config "github.com/ahmagdy/lyricsify/config"
	"github.com/ahmagdy/lyricsify/internal/scraper"
	"github.com/ahmagdy/lyricsify/internal/search"
	spotifyService "github.com/ahmagdy/lyricsify/internal/spotify"
	"github.com/google/wire"
	"github.com/olivere/elastic/v7"
	"github.com/zmb3/spotify/v2"
	"go.uber.org/zap"
)

func New(ctx context.Context, spotifyClient *spotify.Client) (*Service, error) {
	wire.Build(
		createLogger,
		createElasticClient,
		config.New,
		spotifyService.New,
		scraper.New,
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

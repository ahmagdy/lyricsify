//+build wireinject

package lyricsify

import (
	"context"

	config "github.com/Ahmad-Magdy/lyricsify/config"
	scrapping "github.com/Ahmad-Magdy/lyricsify/scraping"
	"github.com/Ahmad-Magdy/lyricsify/search"
	"github.com/Ahmad-Magdy/lyricsify/spotify"
	"github.com/google/wire"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

func InitializeLyricsify(ctx context.Context) (*Lyricsify, error) {
	wire.Build(
		createLogger,
		createElasticClient,
		config.New,
		spotify.New,
		scrapping.New,
		search.New,
		New)
	return &Lyricsify{}, nil
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

//+build wireinject

package lyricsify

import (
	"context"

	"github.com/Ahmad-Magdy/lyricsify/search"
	config "github.com/Ahmad-Magdy/lyricsify/config"

	scrapping "github.com/Ahmad-Magdy/lyricsify/scraping"
	"github.com/Ahmad-Magdy/lyricsify/spotify"
	"github.com/google/wire"
)

func InitializeLyricsify(ctx context.Context) *Lyricsify {
	wire.Build(
		config.NewConfig,
		spotify.New,
		scrapping.New,
		search.New,
		New)
	return &Lyricsify{}
}

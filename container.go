//+build wireinject

package lyricsify

import (
	"context"
	"github.com/Ahmad-Magdy/lyricsify/elasticclient"
	config "github.com/Ahmad-Magdy/lyricsify/internal"

	scrapping "github.com/Ahmad-Magdy/lyricsify/scraping"
	"github.com/Ahmad-Magdy/lyricsify/spotifyservice"
	"github.com/google/wire"
)

func CreateSomething(ctx context.Context) *Lyricsify {
	//elasticService, err := elasticclient.New(ctx,"lyrics")
	//if err != nil{
	//	log.Fatalf("Error")
	//}
	//		wire.InterfaceValue(elasticclient.LyricsSearchService{}, elasticclient.NewX(ctx,"lyrics")),
	wire.Build(
		config.NewConfig,
		spotifyservice.New,
		scrapping.New,
		elasticclient.New,
		NewLyricsifyService)
	return &Lyricsify{}
}

//+build wireinject

package lyricsify

import (
	"context"

	"github.com/Ahmad-Magdy/lyricsify/elasticclient"
	scrapping "github.com/Ahmad-Magdy/lyricsify/scraping"
	"github.com/Ahmad-Magdy/lyricsify/spotifyservice"
	"github.com/google/wire"
)

func CreateSomething(ctx context.Context) *Lyricsify {
	//elasticService, err := elasticclient.New(ctx,"lyrics")
	//if err != nil{
	//	log.Fatalf("Error")
	//}
	wire.Build(
		spotifyservice.New,
		scrapping.New,
		elasticclient.NewOne,
		NewLyricsifyService)
	return &Lyricsify{}
}

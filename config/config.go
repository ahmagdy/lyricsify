package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	LyricsIndexName string
	SpotifyID       string
	SpotifySecret   string
	GeniusToken     string
	MusixmatchToken string
	LyricsProvider  string
}

func New() (*Config, error) {
	viperConfig := viper.New()
	viperConfig.SetConfigType("yaml")
	viperConfig.AddConfigPath("$HOME/Documents/lyricsify")

	err := viperConfig.ReadInConfig()
	if err != nil {
		return nil, err
	}
	viperConfig.AutomaticEnv()
	viperConfig.SetDefault("LYRICS_INDEX_NAME", "lyrics")
	viperConfig.SetDefault("LYRICS_PROVIDER", "genius")

	return &Config{
		LyricsIndexName: viperConfig.GetString("LYRICS_INDEX_NAME"),
		SpotifyID:       viperConfig.GetString("SPOTIFY_ID"),
		SpotifySecret:   viperConfig.GetString("SPOTIFY_SECRET"),
		GeniusToken:     viperConfig.GetString("GENIUS_TOKEN"),
		MusixmatchToken: viperConfig.GetString("MUSIXMATCH_TOKEN"),
		LyricsProvider:  viperConfig.GetString("LYRICS_PROVIDER"),
	}, nil
}

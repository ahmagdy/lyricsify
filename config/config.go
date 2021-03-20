package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	LyricsIndexName string
	SpotifyToken    string
	GeniusToken     string
	GeniusBaseURL   string
}

func New() (*Config, error) {
	viperConfig := viper.New()
	viperConfig.SetConfigType("yaml")
	viperConfig.AddConfigPath("$HOME/Documents")

	err := viperConfig.ReadInConfig()
	if err != nil {
		return nil, err
	}
	viperConfig.AutomaticEnv()
	viperConfig.SetDefault("LYRICS_INDEX_NAME", "lyrics")

	return &Config{
		LyricsIndexName: viperConfig.GetString("LYRICS_INDEX_NAME"),
		SpotifyToken:    viperConfig.GetString("SPOTIFY_TOKEN"),
		GeniusToken:     viperConfig.GetString("GENIUS_TOKEN"),
		GeniusBaseURL:   viperConfig.GetString("GENIUS_BASE_URL"),
	}, nil
}

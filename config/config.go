package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	LyricsIndexName string
	SpotifyToken    string
	GeniusToken     string
	GeniusBaseURL   string
}

func NewConfig() *Config {
	viperConfig := viper.New()
	viperConfig.SetConfigType("yaml")
	viperConfig.AddConfigPath("$HOME/Documents")

	err := viperConfig.ReadInConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}
	viperConfig.AutomaticEnv()
	viperConfig.SetDefault("LYRICS_INDEX_NAME", "lyrics")

	return &Config{
		LyricsIndexName: viperConfig.GetString("LYRICS_INDEX_NAME"),
		SpotifyToken:    viperConfig.GetString("SPOTIFY_TOKEN"),
		GeniusToken:     viperConfig.GetString("GENIUS_TOKEN"),
		GeniusBaseURL:   viperConfig.GetString("GENIUS_BASE_URL"),
	}
}

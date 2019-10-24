package config

type Config struct{
	LyricsIndexName string ``
}

func NewConfig() *Config {
	return &Config{
		LyricsIndexName: "SDsd"}
}
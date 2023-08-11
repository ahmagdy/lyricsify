package scraper

import (
	"fmt"

	config "github.com/ahmagdy/lyricsify/config"
	"go.uber.org/zap"
)

type provider int

const (
	unknownProvider provider = iota
	geniusProvider
	musixmatchProvider
)

var providers = map[string]provider{
	"genius":     geniusProvider,
	"musixmatch": musixmatchProvider,
}

// New creates a new instance of lyrics scraper service
func New(logger *zap.Logger, config *config.Config) (Service, error) {
	provider := getProvider(config)
	if provider == unknownProvider {
		return nil, fmt.Errorf("unknown provider (%v)", config.LyricsProvider)
	}

	logger.Info("using lyrics provider", zap.String("provider", config.LyricsProvider))

	if provider == geniusProvider {
		return newGeniusService(logger, config)
	}

	return newMuxixmatchService(logger, config), nil
}

func getProvider(config *config.Config) provider {
	return providers[config.LyricsProvider]
}

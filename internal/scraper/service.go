package scraper

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	config "github.com/ahmagdy/lyricsify/config"
	"github.com/ahmagdy/lyricsify/internal/types"
	"go.uber.org/zap"
)

var (
	ErrGeniusTokenNotSet = errors.New("genius token is not set")
	ErrRequestFailed     = errors.New("request failed with non-success code")
	errZeroHits          = errors.New("received zero hits")
	errNoMatchingTag     = errors.New("could not find a matching tag for the lyrics content")
	_baseURL             = "https://api.genius.com/search"
)

// Service a service to scrap song lyrics from the internet
type Service struct {
	logger            *zap.Logger
	config            *config.Config
	geniusAccessToken string
}

// New creates a new instance of lyrics scraper service
func New(config *config.Config, logger *zap.Logger) (*Service, error) {
	geniusAccessToken := config.GeniusToken
	if geniusAccessToken == "" {
		return nil, ErrGeniusTokenNotSet
	}

	return &Service{
		logger:            logger,
		config:            config,
		geniusAccessToken: geniusAccessToken,
	}, nil
}

// FindLyrics Get song lyrics
func (s *Service) FindLyrics(ctx context.Context, songName string, artists string) (string, error) {
	lyricsURL, err := s.fetchSongLyricsResults(ctx, songName, artists)
	if err != nil {
		return "", fmt.Errorf("couldn't find lyrics for song (%v): %w", songName, err)
	}

	s.logger.Debug("Calling URL", zap.String("url", lyricsURL))

	res, err := http.Get(lyricsURL)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	var lyrics string
	doc.Find(`div[class^="Lyrics__Container-"]`).Each(func(i int, s *goquery.Selection) {
		lyrics = s.Text()
	})

	if lyrics == "" {
		return "", errNoMatchingTag
	}

	return lyrics, nil
}

// fetchSongLyricsResults returns the URLs lyrics search results that matches the song.
func (s *Service) fetchSongLyricsResults(ctx context.Context, songName string, artists string) (string, error) {
	req, err := http.NewRequest("GET", _baseURL, nil)
	if err != nil {
		return "", err
	}

	queryParams := req.URL.Query()
	queryParams.Add("q", fmt.Sprintf("%s %s", songName, artists))
	req.URL.RawQuery = queryParams.Encode()
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", s.geniusAccessToken))
	req = req.WithContext(ctx)

	s.logger.Info("calling URL for song lyrics", zap.String("URL", req.URL.String()))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		s.logger.Error("songLyricsResults: request exited with non-success code",
			zap.Any("response", res),
			zap.String("url", res.Request.URL.String()),
			zap.Int("statusCode", res.StatusCode),
		)
		return "", ErrRequestFailed
	}

	var geniusResponse types.GeniusResponse
	err = json.NewDecoder(res.Body).Decode(&geniusResponse)
	if err != nil {
		return "", fmt.Errorf("json.NewDecoder: %w", err)
	}

	if len(geniusResponse.Response.Hits) == 0 {
		return "", errZeroHits
	}

	matchingURL, ok := s.findSongLyricsURLInResponse(artists, geniusResponse)

	if !ok {
		return "", fmt.Errorf("could not find matching lyrics for the song (%s)", songName)
	}

	return matchingURL, nil
}

func (s *Service) findSongLyricsURLInResponse(artists string, geniusResponse types.GeniusResponse) (string, bool /* ok */) {
	splitArtists := strings.Split(artists, ",")

	for _, hitItem := range geniusResponse.Response.Hits {
		for _, artist := range splitArtists {
			if !strings.Contains(hitItem.Result.PrimaryArtist.Name, artist) {
				continue
			}
			s.logger.Debug("found artist as part ", zap.String("artist", artist), zap.String("geniusArtist", hitItem.Result.PrimaryArtist.Name))
			return hitItem.Result.URL, true
		}

	}
	return "", false
}

func LoadCSV() ([][]string, error) {
	file, err := os.Open("../results.csv")
	if err != nil {
		return nil, fmt.Errorf("load csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '|'

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("reader.ReadAll: %w", err)
	}

	return records, nil
}

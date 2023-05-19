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
	"github.com/ahmagdy/lyricsify/types"
	"go.uber.org/zap"
)

var (
	ErrGeniusTokenNotSet = errors.New("genius token is not set")
	ErrRequestFailed     = errors.New("request failed with non-success code")
	_baseURL             = "https://api.genius.com/search"
)

// Service a service to scrap song lyrics from the internet
type Service struct {
	logger *zap.Logger
	config *config.Config
}

// New creates a new instance of lyrics scraper service
func New(config *config.Config, logger *zap.Logger) *Service {
	return &Service{
		logger: logger,
		config: config,
	}
}

// Lyrics Get song lyrics
func (s *Service) Lyrics(ctx context.Context, songName string, artists string) (string, error) {
	songInfo, err := s.songLyricsResults(ctx, songName, artists)
	if err != nil {
		return "", fmt.Errorf("couldn't find lyriccs for song (%v): %w", songName, err)
	}
	if songInfo.Type == "" {
		return "", fmt.Errorf("couldn't find lyriccs for song (%v)", songName)
	}

	s.logger.Debug("Calling URL", zap.String("url", songInfo.Result.URL))

	res, err := http.Get(songInfo.Result.URL)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	var lyrics string
	doc.Find("div.lyrics").Each(func(i int, s *goquery.Selection) {
		lyrics = s.Text()
	})

	return lyrics, nil
}

// songLyricsResults Search for a song lyrics and get the results list of the search. It doesn't contain the actual lyrics
func (s *Service) songLyricsResults(ctx context.Context, songName string, artists string) (*types.SearchResult, error) {
	geniusAccessToken := s.config.GeniusToken
	if geniusAccessToken == "" {
		return &types.SearchResult{}, ErrGeniusTokenNotSet
	}

	req, err := http.NewRequest("GET", _baseURL, nil)
	if err != nil {
		return nil, err
	}

	queryParams := req.URL.Query()
	queryParams.Add("q", fmt.Sprintf("%s %s", songName, artists))
	req.URL.RawQuery = queryParams.Encode()
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", geniusAccessToken))
	req = req.WithContext(ctx)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &types.SearchResult{}, err
	}
	if res.StatusCode != 200 {
		s.logger.Error("songLyricsResults: request exited with non-success code",
			zap.String("url", res.Request.URL.String()), zap.Int("statusCode", res.StatusCode))
		return &types.SearchResult{}, ErrRequestFailed
	}

	var geniusResponse types.GeniusResponse
	err = json.NewDecoder(res.Body).Decode(&geniusResponse)
	if err != nil {
		return &types.SearchResult{}, err
	}

	var songSearchResult types.SearchResult
	singersList := strings.Split(artists, ",")
	breakOuterLoop := false
	for _, hitItem := range geniusResponse.Response.Hits {
		for _, singer := range singersList {
			if strings.Contains(hitItem.Result.PrimaryArtist.Name, singer) {
				s.logger.Debug("found singer as part ", zap.String("singer", singer), zap.String("geniusArtist", hitItem.Result.PrimaryArtist.Name))
				songSearchResult = hitItem
				breakOuterLoop = true
				break
			}
		}
		if breakOuterLoop {
			break
		}

	}

	return &songSearchResult, nil
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

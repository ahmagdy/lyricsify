package scrapping

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Ahmad-Magdy/lyricsify/errorhandler"
	"github.com/Ahmad-Magdy/lyricsify/models"
	"github.com/PuerkitoBio/goquery"
)

// LyricsScrapingService Service to get song Lyrics from the internet
type LyricsScrapingService struct {
	baseSearchURL string
}

// New Create
func New() *LyricsScrapingService {
	return &LyricsScrapingService{"https://api.genius.com/search"}
}

func (songService *LyricsScrapingService) getSongLyricsResults(ctx context.Context, songName string, artists string) (searchResults models.SearchResult, err error) {
	geniusAccessToken := os.Getenv("GENIUS_TOKEN")
	if geniusAccessToken == "" {
		return models.SearchResult{}, errors.New("GENIUS_TOKEN environment variable is not found.")
	}
	req, _ := http.NewRequest("GET", songService.baseSearchURL, nil)
	queryParams := req.URL.Query()
	queryParams.Add("q", songName+" "+artists)
	req.URL.RawQuery = queryParams.Encode()
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", geniusAccessToken))
	req = req.WithContext(ctx)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.SearchResult{}, err
	}
	if res.StatusCode != 200 {
		errText := fmt.Sprintf("getSongLyricsResults: Request with URL %v exit with code %v", res.Request.URL, res.StatusCode)
		return models.SearchResult{}, errors.New(errText)
	}

	var geniusResponse models.GeniusResponse
	err = json.NewDecoder(res.Body).Decode(&geniusResponse)
	if err != nil {
		return models.SearchResult{}, err
	}

	var songSearchResult models.SearchResult
	singersList := strings.Split(artists, ",")
	breakOuterLoop := false
	for _, hitItem := range geniusResponse.Response.Hits {
		for _, singer := range singersList {
			if strings.Contains(hitItem.Result.PrimaryArtist.Name, singer) {
				log.Println("Found: " + singer + " as part of " + hitItem.Result.PrimaryArtist.Name)
				songSearchResult = hitItem
				breakOuterLoop = true
				break
			}
		}
		if breakOuterLoop {
			break
		}

	}
	log.Println(songSearchResult.Result.PrimaryArtist.Name)
	log.Println(songSearchResult.Result.URL)
	return songSearchResult, nil
}

// GetLyricsForSong Get song lyrics
func (songService *LyricsScrapingService) GetLyricsForSong(ctx context.Context, songName string, artists string) (lyricsText string, err error) {
	songInfo, err := songService.getSongLyricsResults(ctx, songName, artists)
	if songInfo.Type == "" {
		err = fmt.Errorf("Couldn't find scraping for song %v", songName)
		return "", err
	}
	log.Printf("Calling URL: %v", songInfo.Result.URL)
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
	doc.Find("div.scraping").Each(func(i int, s *goquery.Selection) {
		lyrics = s.Text()
	})
	return lyrics, nil
}

func LoadCSV() [][]string {
	file, err := os.Open("../results.csv")
	errorhandler.HandleError(err)
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = '|'
	records, err := reader.ReadAll()
	errorhandler.HandleError(err)
	fmt.Println(records[1][0], records[1][1])
	return records
}

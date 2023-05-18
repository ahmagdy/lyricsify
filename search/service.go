package search

import (
	"context"
	"encoding/json"

	"github.com/ahmagdy/lyricsify/config"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

const _mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
			"properties":{
				"title": {
					"type":"text"
				},
				"content": {
					"type":"text"
				}
			}
	}
} `

// LyricsBody lyrics inserted body
type LyricsBody struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Service a service to create lyrics documents and search across all lyrics
type Service struct {
	logger    *zap.Logger
	esClient  *elastic.Client
	indexName string
}

// New to create instance of search service
func New(ctx context.Context, config *config.Config, client *elastic.Client, logger *zap.Logger) (*Service, error) {
	esClient := &Service{logger, client, config.LyricsIndexName}
	// TODO [ahmed]: This shouldn't be placed here
	err := esClient.createIndexIfNotExist(ctx)
	if err != nil {
		return nil, err
	}
	return esClient, err
}

// Create Create ES Document.
func (s *Service) Create(ctx context.Context, title string, content string) error {
	res, err := s.esClient.Index().
		Index(s.indexName).
		BodyJson(LyricsBody{title, content}).
		Do(ctx)
	if err != nil {
		return err
	}

	s.logger.Info("Created item with ID", zap.String("id", res.Id))
	_, err = s.esClient.Flush().Index(s.indexName).Do(ctx)
	return err
}

// Update ES Document.
func (s *Service) Update(ctx context.Context, id string, title string, content string) (err error) {

	res, err := s.esClient.Update().
		Index(s.indexName).
		Id(id).
		Script(elastic.NewScript("ctx._source.content = params.content; ctx._source.title = params.title").
			Param("title", title).
			Param("content", content)).
		Upsert(map[string]interface{}{}).
		Do(ctx)

	if err != nil {
		return err
	}

	s.logger.Info("Updated item with ID", zap.String("id", res.Id))

	_, err = s.esClient.Flush().Index(s.indexName).Do(ctx)
	return err
}

// GetItemID Get Item ID
func (s *Service) GetItemID(ctx context.Context, title string) (id string, err error) {
	matchQuery := elastic.NewMatchQuery("title", title)

	searchResult, err := s.esClient.Search().
		Index(s.indexName).
		Query(matchQuery).
		Pretty(true).
		Do(ctx)
	if err != nil {
		return "", err
	}

	s.logger.Info("Item search match", zap.Int64("hits", searchResult.TotalHits()))

	if len(searchResult.Hits.Hits) == 0 {
		return "", nil
	}

	firstHit := searchResult.Hits.Hits[0]
	return firstHit.Id, nil
}

// Search to search for song lyrics across ES Index.
func (s *Service) Search(ctx context.Context, text string) (lyrics []LyricsBody, err error) {
	//termQuery := elastic.NewTermQuery("content", "All Around The World")
	matchQuery := elastic.NewMultiMatchQuery(text, "title", "content").Type("phrase_prefix")

	searchResult, err := s.esClient.Search().
		Index(s.indexName).
		Query(matchQuery).
		Pretty(true).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	var lyricsResults []LyricsBody
	for _, hit := range searchResult.Hits.Hits {
		var hitLyricsBody LyricsBody
		err := json.Unmarshal(hit.Source, &hitLyricsBody)
		if err != nil {
			return nil, err
		}
		lyricsResults = append(lyricsResults, hitLyricsBody)
		s.logger.Info("search result row", zap.String("type", hit.Type), zap.String("id", hit.Id), zap.String("title", hitLyricsBody.Title))
	}

	s.logger.Info("result", zap.Int64("totalHits", searchResult.Hits.TotalHits.Value))
	return lyricsResults, nil
}

// CreateIndexIfNotExist To check ES Index and create it if it doesn't exist.
func (s *Service) createIndexIfNotExist(ctx context.Context) error {
	exists, err := s.esClient.IndexExists(s.indexName).Do(ctx)
	if err != nil {
		return err
	}
	if !exists {
		s.logger.Info("index doesn't exist, creating a new one", zap.String("indexNAme", s.indexName))
		if _, err := s.esClient.CreateIndex(s.indexName).Body(_mapping).Do(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) deleteByIndex(ctx context.Context, itemID string) (err error) {
	res, err := s.esClient.Delete().Index(s.indexName).Id(itemID).Do(ctx)
	if err != nil {
		return err
	}
	s.logger.Info("delete response", zap.String("response", res.Result))
	return nil
}

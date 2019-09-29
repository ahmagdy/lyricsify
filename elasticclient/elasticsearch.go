package elasticclient

import (
	"context"
	"encoding/json"
	"log"

	"github.com/olivere/elastic/v7"
)

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
			"properties":{
				"title": {
					"type":"keyword"
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

// LyricsSearchService Service to Create Lyrics documents and search across all lyrics
type LyricsSearchService struct {
	esClient  *elastic.Client
	indexName string
}

// New to create instance of this service, TOBE discussed
func New(ctx context.Context, indexName string) (*LyricsSearchService, error) {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}
	info, code, err := client.Ping("http://localhost:9200").Do(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	esClient := &LyricsSearchService{client, indexName}
	err = esClient.createIndexIfNotExist(ctx)
	if err != nil {
		return nil, err
	}
	return esClient, nil
}

// CreateIndexIfNotExist To check ES Index and create it if it doesn't exist.
func (els *LyricsSearchService) createIndexIfNotExist(ctx context.Context) (err error) {
	exists, err := els.esClient.IndexExists(els.indexName).Do(ctx)
	if err != nil {
		return err
	}
	if !exists {
		log.Printf("Index %v doesn't exist, creating a new one.", els.indexName)
		_, err := els.esClient.CreateIndex(els.indexName).Body(mapping).Do(ctx)

		if err != nil {
			return err
		}
	}
	return nil
}

// Create Create ES Document.
func (els *LyricsSearchService) Create(ctx context.Context, title string, content string) (err error) {
	res, err := els.esClient.Index().
		Index(els.indexName).
		BodyJson(LyricsBody{title, content}).
		Do(ctx)
	if err != nil {
		return nil
	}
	log.Printf("Created item with ID %v", res.Id)
	_, err = els.esClient.Flush().Index(els.indexName).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (els *LyricsSearchService) deleteByIndex(ctx context.Context, itemID string) (err error) {
	res, err := els.esClient.Delete().Index(els.indexName).Id(itemID).Do(ctx)
	if err != nil {
		return err
	}
	log.Printf("Delete response %v", res.Result)
	return nil
}

// Search to search for song lyrics across ES Index.
func (els *LyricsSearchService) Search(ctx context.Context, key string, text string) (lyrics []LyricsBody, err error) {
	//termQuery := elastic.NewTermQuery("content", "All Around The World")
	matchQuery := elastic.NewMatchQuery(key, text)
	searchResult, err := els.esClient.Search().
		Index(els.indexName).
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
		log.Println(hit.Type, hit.Id, hitLyricsBody.Title)
	}
	log.Printf("Results %v", searchResult.Hits.TotalHits.Value)
	return lyricsResults, nil
}

package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"movie-service/internal/models"
)

type ElasticSearchRepository interface {
	IndexMovie(movieID string, movieData models.Movie) error
	SearchMovies(query string) ([]map[string]interface{}, error)
}

type elasticSearchRepository struct {
	client *elasticsearch.Client
}

func NewElasticSearchRepository() (ElasticSearchRepository, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://elasticsearch:9200", // Адрес вашего сервера Elasticsearch в Docker
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("Error creating es search: %w", err)
	}

	return &elasticSearchRepository{es}, nil
}

func (r *elasticSearchRepository) IndexMovie(movieID string, movieData models.Movie) error {
	data, err := json.Marshal(movieData)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      "movies",
		DocumentID: movieID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), r.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("Error indexing movie: %s", res.String())
	}

	return nil
}

func (r *elasticSearchRepository) SearchMovies(query string) ([]map[string]interface{}, error) {
	var buf bytes.Buffer
	queryStr := fmt.Sprintf(`{
		"query": {
			"match": {
				"title": %s
			}
		}
	}`, query)
	buf.WriteString(queryStr)

	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex("movies"),
		r.client.Search.WithBody(&buf),
		r.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error searching movies: %s", res.String())
	}

	var rmap map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&rmap); err != nil {
		return nil, err
	}

	hits := rmap["hits"].(map[string]interface{})["hits"].([]interface{})
	var results []map[string]interface{}
	for _, hit := range hits {
		results = append(results, hit.(map[string]interface{})["_source"].(map[string]interface{}))
	}

	return results, nil
}

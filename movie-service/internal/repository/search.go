package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"io"
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
		return fmt.Errorf("could not marshal movie data: %w", err)
	}

	fmt.Printf("Indexing movie data: 	%s\n", string(data))

	res, err := r.client.Index(
		"movies", // Не используем fmt.Sprintf, просто конкатенация строк
		bytes.NewReader(data),
		r.client.Index.WithDocumentID(movieID), // Здесь передаем ID документа
		r.client.Index.WithRefresh("true"),
	)

	// Убедитесь, что res не nil перед вызовом IsError
	if res == nil {
		return fmt.Errorf("response is nil")
	}

	if res.IsError() {
		body, _ := io.ReadAll(res.Body) // Читаем тело ответа для диагностики
		return fmt.Errorf("error indexing movie: %s, response body: %s", res.String(), string(body))
	}

	if err != nil {
		return fmt.Errorf("could not index movie: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing movie: %s", res.String())
	}

	return nil
}

func (r *elasticSearchRepository) SearchMovies(query string) ([]map[string]interface{}, error) {
	var buf bytes.Buffer

	queryStr := fmt.Sprintf(`{
		"query": {
			"match_phrase_prefix": {
				"title": {
					"query": "%s"
				}
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

	// Извлекаем результаты
	hits := rmap["hits"].(map[string]interface{})["hits"].([]interface{})
	var results []map[string]interface{}
	for _, hit := range hits {
		results = append(results, hit.(map[string]interface{})["_source"].(map[string]interface{}))
	}

	return results, nil
}

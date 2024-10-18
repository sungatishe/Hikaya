package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (m *MovieHandler) SearchMovies(rw http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(rw, "Invalid query", http.StatusBadRequest)
		return
	}

	results, err := m.movieService.SearchMovies(query)
	log.Println(results)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Invalid query: %s ", err.Error()), http.StatusBadRequest)
		log.Println(err)
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(results)
}

package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func Alias(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		alias, _ := parseAlias(r)
		log.Printf("Searcing alias: %s", alias)

	})
}

func parseAlias(r *http.Request) (string, error) {
	return r.URL.Path, nil
}

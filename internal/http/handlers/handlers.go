package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func Alias(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias, _ := parseAlias(r)
		log.Printf("Searcing alias: %s", alias)
		var url string
		err := db.QueryRowContext(r.Context(),
			"SELECT target_url FROM link WHERE alias=$1", alias,
		).Scan(&url)
		if err != nil {
			log.Println(err)
			http.Error(w, "not found", 404)
			return
		}
		log.Printf("Redirect to %s", url)
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func parseAlias(r *http.Request) (string, error) {
	return r.URL.Path, nil
}

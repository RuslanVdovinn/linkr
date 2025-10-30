package handlers

import (
	"encoding/json"
	"linkr/internal/domain"
	"log"
	"net/http"

	"gorm.io/gorm"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func Alias(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias, _ := parseAlias(r)
		log.Printf("Searcing alias: %s", alias)
		var link domain.Link
		if err := db.First(&link, "alias = ?", alias).Error; err != nil {
			http.Error(w, "not found", 404)
			return
		}
		// err := db.QueryRowContext(r.Context(),
		// 	"SELECT target_url FROM link WHERE alias=$1", alias,
		// ).Scan(&url)
		// if err != nil {
		// 	log.Println(err)
		// 	http.Error(w, "not found", 404)
		// 	return
		// }
		log.Printf("Redirect to %s", link.TargetURL)
		http.Redirect(w, r, link.TargetURL, http.StatusFound)
	}
}

func parseAlias(r *http.Request) (string, error) {
	return r.URL.Path, nil
}

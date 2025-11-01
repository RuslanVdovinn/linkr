package handlers

import (
	"encoding/json"
	"linkr/internal/domain"
	"log"
	"net/http"
	"strings"

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
		log.Printf("Redirect to %s", link.TargetURL)
		http.Redirect(w, r, link.TargetURL, http.StatusFound)
	}
}

func GetAlias(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias, _ := parseAlias(r)
		log.Printf("Searcing alias: %s", alias)
		var link domain.Link
		if err := db.First(&link, "alias = ?", alias).Error; err != nil {
			http.Error(w, "not found", 404)
			return
		}
		json.NewEncoder(w).Encode(link)
	}
}

func PatchAlias(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias, _ := parseAlias(r)
		log.Printf("Patching alias: %s", alias)
		var link domain.Link
		if err := db.First(&link, "alias = ?", alias).Error; err != nil {
			http.Error(w, "not found", 404)
			return
		}
		json.NewDecoder(r.Body).Decode(&link)
		if err := db.Save(&link).Error; err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		json.NewEncoder(w).Encode(link)
	}
}

func DeleteAlias(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias, _ := parseAlias(r)
		log.Printf("Deleting alias: %s", alias)
		var link domain.Link
		link.Alias = alias
		if err := db.Delete(&link, "alias = ?", alias).Error; err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		log.Printf("Successfully deleted alias: %s", alias)
		w.Write([]byte("{'status':'ok'}"))
	}
}

func CreateLink(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var link domain.Link
		json.NewDecoder(r.Body).Decode(&link)
		user := domain.AppUser{
			ID:    1,
			Email: "email",
			Name:  "test",
		}
		link.User = &user
		log.Printf("Req: %v", link)
		if err := db.Create(&link).Error; err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		w.Write([]byte("{'status':'ok'}"))
	}
}

func parseAlias(r *http.Request) (string, error) {
	splits := strings.Split(r.URL.Path, "/")
	return splits[len(splits)-1], nil
}

package domain

import (
	"time"

	"github.com/lib/pq"
)

type Link struct {
	ID        int64          `gorm:"column:id;primaryKey"               json:"id"`
	UserID    *int64         `gorm:"column:user_id"                     json:"user_id"`
	User      *AppUser       `gorm:"foreignKey:UserID;references:ID"    json:"user,omitempty"`
	Alias     string         `gorm:"column:alias"                       json:"alias"`
	TargetURL string         `gorm:"column:target_url"                  json:"target_url"`
	Title     *string        `gorm:"column:title"                       json:"title,omitempty"`
	Tags      pq.StringArray `gorm:"column:tags;type:text[]"            json:"tags"`
	IsActive  bool           `gorm:"column:is_active"                   json:"is_active"`
	ExpireAt  *time.Time     `gorm:"column:expire_at"                   json:"expire_at,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at"                  json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"                  json:"updated_at"`
}

func (Link) TableName() string { return "link" }

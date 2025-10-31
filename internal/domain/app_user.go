package domain

import "time"

type AppUser struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	Email     string    `gorm:"column:email"`
	Name      string    `gorm:"column:name"`
	CreatedAt time.Time `gorm:"column:created_at"                  json:"created_at"`
}

func (AppUser) TableName() string { return "app_user" }

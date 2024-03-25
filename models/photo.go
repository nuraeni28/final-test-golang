package models

import (
	"errors"
	"time"
)

type Photo struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoURL  string    `json:"photo_url"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Photo) Validate() error {

	if p.Title == "" {
		return errors.New("title is required")
	}

	// Validate caption field
	if p.PhotoURL == "" {
		return errors.New("Url is required")
	}

	return nil
}

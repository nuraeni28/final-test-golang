package models

import (
	"errors"
	"time"
)

type SocialMedia struct {
	ID             uint      `gorm:"primary_key" json:"id"`
	Name           string    `gorm:"not null" json:"name"`
	SocialMediaURL string    `gorm:"not null" json:"social_media_url"`
	UserID         uint      `json:"user_id"`
	User           User      `gorm:"foreignKey:UserID" json:"user"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (p *SocialMedia) Validate() error {

	if p.Name == "" {
		return errors.New("Name is required")
	}
	if p.SocialMediaURL == "" {
		return errors.New("Url is required")
	}

	return nil
}

package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	PhotoID   uint      `json:"photo_id"`
	Photo     Photo     `gorm:"foreignKey:PhotoID" json:"photo"`
	Message   string    `gorm:"not null" json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Comment) Validate(db *gorm.DB) error {
	// Validasi bahwa message tidak kosong
	if p.Message == "" {
		return errors.New("message is required")
	}

	var photo Photo
	if err := db.First(&photo, p.PhotoID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("photoID does not exist")
		}
		return err
	}

	return nil
}
func (p *Comment) ValidateEdit(db *gorm.DB) error {
	// Validasi bahwa message tidak kosong
	if p.Message == "" {
		return errors.New("message is required")
	}

	return nil
}

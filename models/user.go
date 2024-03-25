package models

import (
	"errors"
	"regexp"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"password"`
	Age       uint      `gorm:"not null" json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate function to validate user fields
func (u *User) Validate(db *gorm.DB) error {
	// Email validation
	if u.Email == "" {
		return errors.New("email is required")
	}
	if !isValidEmail(u.Email) {
		return errors.New("invalid email format")
	}

	// Username validation
	if u.Username == "" {
		return errors.New("username is required")
	}

	// Password validation
	if u.Password == "" {
		return errors.New("password is required")
	}
	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	// Age validation
	if u.Age == 0 {
		return errors.New("age is required")
	}
	if u.Age <= 8 {
		return errors.New("age must be greater than 8")
	}

	// Check uniqueness of email
	var existingEmailUser User
	if err := db.Where("email = ?", u.Email).First(&existingEmailUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No existing user found with the same email
		} else {
			// Other database error occurred
			return err
		}
	} else {
		// Email is not unique
		return errors.New("email must be unique")
	}

	// Check uniqueness of username
	var existingUsernameUser User
	if err := db.Where("username = ?", u.Username).First(&existingUsernameUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No existing user found with the same username
		} else {
			// Other database error occurred
			return err
		}
	} else {
		// Username is not unique
		return errors.New("username must be unique")
	}

	return nil
}

// Helper function to validate email format
func isValidEmail(email string) bool {
	// Email format validation using regular expression
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(emailRegex).MatchString(email)
}

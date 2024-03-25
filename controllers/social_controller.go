package controllers

import (
	"MyGram/config"
	"MyGram/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateSocialMediaProfile(c *gin.Context) {
	// Mendapatkan input data dari body request
	var social models.SocialMedia
	if err := c.ShouldBindJSON(&social); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi input data
	if err := social.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var userID uint
	switch v := userIDInterface.(type) {
	case float64:
		userID = uint(v)
	case int:
		userID = uint(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user ID"})
		return
	}

	social.UserID = userID
	social.CreatedAt = time.Now()
	social.UpdatedAt = time.Now()

	// Simpan profil ke database
	if err := config.GetDB().Create(&social).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create social media profile"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": social})
}
func GetSocialMediaWithUser(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var userID uint
	switch v := userIDInterface.(type) {
	case float64:
		userID = uint(v)
	case int:
		userID = uint(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user ID"})
		return
	}
	db := config.GetDB()

	var socialMediaWithUser []models.SocialMedia
	if err := db.Where("user_id = ?", userID).Preload("User").Find(&socialMediaWithUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch social media"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": socialMediaWithUser})
}

func UpdateSocialMedia(c *gin.Context) {

	socialMediaID := c.Param("socialmediaID")

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var userID uint
	switch v := userIDInterface.(type) {
	case float64:
		userID = uint(v)
	case int:
		userID = uint(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user ID"})
		return
	}

	var updatedSocialMedia models.SocialMedia

	if err := c.ShouldBindJSON(&updatedSocialMedia); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := updatedSocialMedia.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.GetDB()

	var socialMediaToUpdate models.SocialMedia
	if err := db.Where("id = ? AND user_id = ?", socialMediaID, userID).First(&socialMediaToUpdate).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Social media not found"})
		return
	}

	socialMediaToUpdate.Name = updatedSocialMedia.Name
	socialMediaToUpdate.SocialMediaURL = updatedSocialMedia.SocialMediaURL
	socialMediaToUpdate.UpdatedAt = time.Now()

	if err := db.Save(&socialMediaToUpdate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update social media"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": socialMediaToUpdate})
}
func DeleteSocialMedia(c *gin.Context) {

	socialMediaID := c.Param("socialmediaID")

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var userID uint
	switch v := userIDInterface.(type) {
	case float64:
		userID = uint(v)
	case int:
		userID = uint(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user ID"})
		return
	}

	db := config.GetDB()

	var socialMediaToDelete models.SocialMedia
	if err := db.Where("id = ? AND user_id = ?", socialMediaID, userID).First(&socialMediaToDelete).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Social media not found"})
		return
	}

	if err := db.Delete(&socialMediaToDelete).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete social media"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "your social media has been  successfully deleted"})
}

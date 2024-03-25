package controllers

import (
	"net/http"
	"time"

	"MyGram/config"
	"MyGram/models"

	"github.com/gin-gonic/gin"
)

func CreateComment(c *gin.Context) {
	var commentData models.Comment
	if err := c.ShouldBindJSON(&commentData); err != nil {
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
	db := config.GetDB()

	if err := commentData.Validate(db); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	commentData.UserID = userID
	commentData.CreatedAt = time.Now()
	commentData.UpdatedAt = time.Now()

	if err := db.Create(&commentData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create message"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": commentData})
}
func GetComment(c *gin.Context) {
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
	var comment []models.Comment
	if err := db.Preload("User").Preload("Photo").Where("user_id = ?", userID).Find(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	// Mengembalikan respons JSON dengan data comment
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": comment})
}
func UpdateComment(c *gin.Context) {
	db := config.GetDB()
	commentID := c.Param("commentId")
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

	var updatedComment models.Comment
	if err := updatedComment.ValidateEdit(db); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&updatedComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var commentToUpdate models.Comment
	if err := db.First(&commentToUpdate, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}
	if commentToUpdate.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	commentToUpdate.Message = updatedComment.Message
	commentToUpdate.UpdatedAt = time.Now()

	if err := db.Save(&commentToUpdate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	// Mengembalikan respons berhasil
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": commentToUpdate})
}
func DeleteComment(c *gin.Context) {
	db := config.GetDB()
	commentID := c.Param("commentId")

	// Mendapatkan user ID dari token JWT
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

	// Mencari komentar yang akan dihapus
	var commentToDelete models.Comment
	if err := db.First(&commentToDelete, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Memastikan bahwa komentar yang ingin dihapus dimiliki oleh user yang sesuai
	if commentToDelete.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Menghapus komentar dari database
	if err := db.Delete(&commentToDelete).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	// Mengembalikan respons berhasil
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "your comment has been  successfully deleted"})
}

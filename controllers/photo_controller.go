// controllers/post_controller.go
package controllers

import (
	"fmt"
	"net/http"
	"time"

	"MyGram/config"
	"MyGram/models"

	"github.com/gin-gonic/gin"
)

type PhotoResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoURL  string    `json:"photo_url"`
	UpdatedAt time.Time `json:"updated_at"`
}
type PhotoResponseCreate struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoURL  string    `json:"photo_url"`
	CreatedAt time.Time `json:"created_at"`
}

func CreatePost(c *gin.Context) {
	var postData models.Photo
	db := config.GetDB()

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

	if err := c.ShouldBindJSON(&postData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := postData.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	postData.UserID = userID

	postData.CreatedAt = time.Now()
	postData.UpdatedAt = time.Now()
	Data := PhotoResponseCreate{
		ID:        postData.ID,
		Title:     postData.Title,
		Caption:   postData.Caption,
		PhotoURL:  postData.PhotoURL,
		CreatedAt: postData.CreatedAt,
	}
	if err := db.Create(&postData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": Data})
}
func EditPhoto(c *gin.Context) {

	photoID := c.Param("photoID")
	fmt.Println("Type of userID:", photoID)
	db := config.GetDB()

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

	var updatedPhotoData models.Photo
	if err := c.ShouldBindJSON(&updatedPhotoData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi data photo yang akan diubah
	if err := updatedPhotoData.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mencari photo berdasarkan ID
	var photoToUpdate models.Photo
	if err := db.First(&photoToUpdate, photoID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	// Mengecek apakah pengguna yang sedang login adalah pemilik data yang akan diubah
	if photoToUpdate.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	photoToUpdate.Title = updatedPhotoData.Title
	photoToUpdate.Caption = updatedPhotoData.Caption
	photoToUpdate.PhotoURL = updatedPhotoData.PhotoURL
	photoToUpdate.UpdatedAt = time.Now()

	// Simpan perubahan ke dalam database
	if err := db.Save(&photoToUpdate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update photo"})
		return
	}
	Data := PhotoResponse{
		ID:        photoToUpdate.ID,
		Title:     photoToUpdate.Title,
		Caption:   photoToUpdate.Caption,
		PhotoURL:  photoToUpdate.PhotoURL,
		UpdatedAt: photoToUpdate.UpdatedAt,
	}

	// Menyiapkan data respons
	responseData := gin.H{
		"message": "Photo updated successfully",
		"data":    Data,
	}

	// Mengirim respons
	c.JSON(http.StatusOK, responseData)
}

func GetPosts(c *gin.Context) {
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
	var posts []models.Photo
	if err := db.Preload("User").Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	// Mengembalikan respons JSON dengan data posts
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": posts})
}

func DeletePhoto(c *gin.Context) {
	db := config.GetDB()
	photoID := c.Param("photoID")

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

	var photoToDelete models.Photo
	if err := db.First(&photoToDelete, photoID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	// Mengecek apakah pengguna yang sedang login adalah pemilik data yang akan dihapus
	if photoToDelete.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Menghapus photo dari database
	if err := db.Delete(&photoToDelete).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete photo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "your photo has been  successfully deleted"})
}

// controllers/user_controller.go
package controllers

import (
	"MyGram/config"
	"MyGram/models"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var user models.User
	db := config.GetDB()
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := user.Validate(db); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	// Save user to database
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Respond with only the input data
	responseData := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		ID       uint   `json:"id"`
		Age      uint   `json:"age"`
	}{
		Username: user.Username,
		Email:    user.Email,
		ID:       user.ID,
		Age:      user.Age,
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": responseData})
}

func Login(c *gin.Context) {
	var userInput models.User
	db := config.GetDB()

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFromDB models.User
	if err := db.Where("email = ?", userInput.Email).First(&userFromDB).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(userInput.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userFromDB.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expiration time (1 day)
	})

	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func UpdateUser(c *gin.Context) {
	// Mendapatkan user ID dari path parameter
	userID := c.Param("userId")

	userIDFloat, err := strconv.ParseFloat(userID, 64)
	if err != nil {
		userIDType := reflect.TypeOf(userID)
		fmt.Println("Type of userID:", userIDType)
	}
	userIDFromToken, _ := c.Get("user_id")

	if userIDFloat != userIDFromToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var updatedUserData struct {
		Email     string    `json:"email"`
		Username  string    `json:"username"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	if err := c.ShouldBindJSON(&updatedUserData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.GetDB()

	var userToUpdate models.User

	if err := db.First(&userToUpdate, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update email jika ada perubahan
	if updatedUserData.Email != "" {
		userToUpdate.Email = updatedUserData.Email
	}

	// Update username jika ada perubahan
	if updatedUserData.Username != "" {
		userToUpdate.Username = updatedUserData.Username
	}

	// Set updated_at ke waktu saat ini
	userToUpdate.UpdatedAt = time.Now()

	// Simpan perubahan ke dalam database
	if err := db.Save(&userToUpdate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	responseData := gin.H{
		"message":    "User updated successfully",
		"user":       userToUpdate,
		"updated_at": userToUpdate.UpdatedAt.Format(time.RFC3339),
	}
	c.JSON(http.StatusOK, responseData)
}

func DeleteUser(c *gin.Context) {
	// Mendapatkan user ID dari path parameter
	userID := c.Param("userId")

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Account ID"})
		return
	}

	userIDFloat, err := strconv.ParseFloat(userID, 64)
	if err != nil {
		userIDType := reflect.TypeOf(userID)
		fmt.Println("Type of userID:", userIDType)
	}
	userIDFromToken, _ := c.Get("user_id")

	if userIDFloat != userIDFromToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Mendapatkan data pengguna yang akan dihapus dari database
	db := config.GetDB()
	var userToDelete models.User
	if err := db.First(&userToDelete, userIDInt).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	// Hapus data pengguna dari database
	if err := db.Delete(&userToDelete).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Account"})
		return
	}

	// Berhasil menghapus data pengguna
	c.JSON(http.StatusOK, gin.H{"message": "your Account has been  successfully deleted "})
}

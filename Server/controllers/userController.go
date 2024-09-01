package controllers

import (
	"net/http"
	"project-its/initializers"
	"project-its/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var user models.User
	var foundUser models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	result := initializers.DB.Where("email = ?", user.Email).First(&foundUser)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kata sandi salah"})
		return
	}

	token, err := GenerateJWT(foundUser) // Fungsi untuk menghasilkan token JWT
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Simpan token ke database
	userToken := models.UserToken{
		UserID: foundUser.ID,
		Token:  token,
		Expiry: time.Now().Add(72 * time.Hour), // Misalnya token berlaku 24 jam
	}
	initializers.DB.Create(&userToken)

	// Tambahkan informasi pengguna dalam response
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"username": foundUser.Username,
			"email":    foundUser.Email,
		},
	})
}

func GenerateJWT(foundUser models.User) (string, error) {
	claims := jwt.MapClaims{
		"username": foundUser.Username,
		"email":    foundUser.Email,
		"role":     foundUser.Role,                        // Jika ada field role
		"sub":      foundUser.ID,                          // Menyimpan userID di klaim
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token berlaku selama 24 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("rahasia")) // Ganti "rahasia" dengan secret key Anda
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Register(c *gin.Context) {
	var newUser models.User

	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi kata sandi"})
		return
	}

	newUser.Password = string(hashedPassword)

	result := initializers.DB.Create(&newUser)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func GetUserToken(userID uint) (string, error) {
	var userToken models.UserToken
	if err := initializers.DB.Where("user_id = ?", userID).First(&userToken).Error; err != nil {
		return "", err
	}
	return userToken.Token, nil
}

func Logout(c *gin.Context) {
	// Ambil userID dari context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID tidak ditemukan"})
		return
	}

	// Hapus token dari database
	if err := initializers.DB.Where("user_id = ?", userID).Delete(&models.UserToken{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
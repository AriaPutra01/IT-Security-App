package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"project-gin/initializers"
	"project-gin/models"
	"time"
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": foundUser.Email,
		"role":     foundUser.Role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte("rahasia"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	// Tambahkan informasi pengguna dalam response
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"username": foundUser.Username,
			"email": foundUser.Email,
		},
	})
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
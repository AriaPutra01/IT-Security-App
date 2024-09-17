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

type requestUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Info     string `json:"info"`
}

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

	// Simpan token di database
	userToken := models.UserToken{
		UserID: foundUser.ID,
		Token:  token,
		Expiry: time.Now().Add(time.Hour * 1 * 24 * 30), // Token berlaku selama 30 hari
	}
	if err := initializers.DB.Create(&userToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token to database"})
		return
	}

	// Set cookie dengan HttpOnly
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: false,          // Menetapkan cookie sebagai HttpOnly
		MaxAge:   3600 * 24 * 30, // Masa berlaku cookie (30 hari)
		// secure: true, // Uncomment jika menggunakan HTTPS
	})

	// Tambahkan informasi pengguna dalam response
	c.JSON(http.StatusOK, gin.H{
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
		"role":     foundUser.Role,                                 // Jika ada field role
		"sub":      foundUser.ID,                                   // Menyimpan userID di klaim
		"exp":      time.Now().Add(time.Hour * 1 * 24 * 30).Unix(), // Token berlaku selama 30 hari
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("KopikapBasi123||Djarumsuper01||Akuganteng123||qwe234223")) // Ganti "rahasia" dengan secret key Anda
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

	// Hapus cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:   "token",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Menghapus cookie
	})

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func UserIndex(c *gin.Context) {

	// Get models from DB
	var users []models.User
	initializers.DB.Find(&users)

	//Respond with them
	c.JSON(200, gin.H{
		"users": users,
	})
}

func UserUpdate(c *gin.Context) {

	var requestBody requestUser

	if err := c.BindJSON(&requestBody); err != nil {
		c.Status(400)
		c.Error(err) // log the error
		return
	}

	id := c.Params.ByName("id")

	var users models.User
	initializers.DB.First(&users, id)

	if err := initializers.DB.First(&users, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "user tidak ditemukan"})
		return
	}

	if requestBody.Username != "" {
		users.Username = requestBody.Username
	} else {
		users.Username = users.Username // gunakan nilai yang ada dari database
	}

	if requestBody.Email != "" {
		users.Email = requestBody.Email
	} else {
		users.Email = users.Email // gunakan nilai yang ada dari database
	}
	if requestBody.Password != "" {
		users.Password = requestBody.Password
	} else {
		users.Password = users.Password // gunakan nilai yang ada dari database
	}

	initializers.DB.Model(&users).Updates(users)

	c.JSON(200, gin.H{
		"users": users,
	})

}

func UserDelete(c *gin.Context) {

	//get id
	id := c.Params.ByName("id")

	// find the user
	var users models.User

	if err := initializers.DB.First(&users, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "users not found"})
		return
	}

	/// delete it
	if err := initializers.DB.Delete(&users).Error; err != nil {
		c.JSON(404, gin.H{"error": "users Failed to Delete"})
		return
	}

	c.JSON(200, gin.H{
		"users": "users deleted",
	})
}

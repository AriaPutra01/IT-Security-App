package controllers

import (
	"html/template"
	"net/http"
	"project-gin/initializers"
	"project-gin/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AdminCreateUser membuat pengguna baru dengan peran tertentu
func AdminCreateUser(c *gin.Context) {
	// Periksa apakah pengguna adalah admin
	if !IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Ambil data dari formulir
	username := c.PostForm("username")
	password := c.PostForm("password")
	role := c.PostForm("role")

	// Validasi input
	if username == "" || password == "" || role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Buat pengguna baru
	user, err := CreateUser(username, password, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// AdminGetUsers mengambil daftar semua pengguna
func AdminGetUsers(c *gin.Context) {
	// Periksa apakah pengguna adalah admin
	if !IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Ambil semua pengguna
	users, err := GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// aku ingin menajadi mimpi indah dalam hidupku aku ingin kau tahu bahwa ku selalu memujamu tanpamu aku ingin kau tahu bahwa ku selalu memujamu tanpamu
// aku ingin menajadi mimpi indah dalam hidupku aku ingin kau tahu bahwa ku selalu memujamu tanpamu
// aku ingin menajadi mimpi indah dalam hidupku aku ingin kau tahu bahwa ku selalu memujamu tanpamu
// aku ingin menajadi mimpi indah dalam hidupku aku ingin kau tahu bahwa ku selalu memujamu tanpamu

// Logout menghapus sesi pengguna dan mengakhiri sesi
// hanya pengguna yang terautentikasi yang dapat melakukan logout
// Login melakukan autentikasi pengguna dan memulai sesi
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Autentikasi pengguna
	user, err := AuthenticateUser(username, password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Atur sesi
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
}

// GetLoginView menampilkan halaman login
func GetLoginView(c *gin.Context) {
	temp, err := template.ParseFiles("views/login.html")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	temp.Execute(c.Writer, gin.H{})
}

// GetRegisterView menampilkan halaman pendaftaran
func GetRegisterView(c *gin.Context) {
	temp, err := template.ParseFiles("views/register.html")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	temp.Execute(c.Writer, gin.H{})
}

// Logout menghapus sesi pengguna dan mengakhiri sesi
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user_id")
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// IsAdmin memeriksa apakah pengguna adalah admin
func IsAdmin(c *gin.Context) bool {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		return false
	}

	user, err := GetUserByID(userID)
	if err != nil {
		return false
	}

	return user.Role == "admin"
}

// CreateUser membuat pengguna baru di database
func CreateUser(username, password, role string) (*models.User, error) {
	user := &models.User{
		Username: username,
		Password: password, // Perlu di-hash untuk produksi
		Role:     role,
	}

	// Gunakan koneksi database yang ada dari initializers.DB
	db := initializers.DB
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// GetAllUsers mengambil semua pengguna dari database
func GetAllUsers() ([]*models.User, error) {
	var users []*models.User
	db := initializers.DB
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// AuthenticateUser melakukan autentikasi pengguna
func AuthenticateUser(username, password string) (*models.User, error) {
	var user models.User
	db := initializers.DB
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByID mengambil pengguna berdasarkan ID
func GetUserByID(userID interface{}) (*models.User, error) {
	var user models.User
	db := initializers.DB
	if err := db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

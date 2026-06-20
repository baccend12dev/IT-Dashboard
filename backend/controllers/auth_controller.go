package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary      Register a new user
// @Description  Register a new user in the system with a specific role
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body   models.RegisterRequest  true  "Register Request Payload"
// @Success      201  {object}  models.User
// @Failure      400  {object}  map[string]string "error: Invalid request body / Invalid role"
// @Failure      500  {object}  map[string]string "error: Username already exists / Database save error"
// @Router       /api/auth/register [post]
func Register(c *gin.Context) {
	var request models.RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	if !models.IsValidRole(request.Role) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role tidak valid. Role yang diperbolehkan: 'Administrator', 'Developer', 'Viewer'"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses kata sandi"})
		return
	}

	user := models.User{
		Username: request.Username,
		Password: hashedPassword,
		Role:     request.Role,
	}

	// Save to DB
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Username sudah terdaftar atau terjadi masalah database"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary      Authenticate user
// @Description  Login with username and password to retrieve a Bearer token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body   models.LoginRequest  true  "Login Request Payload"
// @Success      200  {object}  models.LoginResponse
// @Failure      400  {object}  map[string]string "error: Invalid request body"
// @Failure      401  {object}  map[string]string "error: Username atau password salah"
// @Router       /api/auth/login [post]
func Login(c *gin.Context) {
	var request models.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", request.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	// Verify password
	if !utils.CheckPasswordHash(request.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	// Generate Token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token otorisasi"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
		User:  user,
	})
}

// Me godoc
// @Summary      Get current user profile
// @Description  Get the profile details of the currently authenticated user
// @Tags         auth
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      401  {object}  map[string]string "error: Unauthorized / User tidak ditemukan"
// @Router       /api/auth/me [get]
func Me(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Logout godoc
// @Summary      Logout user
// @Description  Logout user session (returns success confirmation)
// @Tags         auth
// @Produce      json
// @Success      200  {object}  map[string]string "message: Logout successfully"
// @Router       /api/auth/logout [post]
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logout successfully"})
}

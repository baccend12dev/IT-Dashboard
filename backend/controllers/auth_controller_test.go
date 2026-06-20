package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/config"
	"backend/models"
	"backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupAuthTestDB(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	config.DB = db
}

func TestAuthRegister(t *testing.T) {
	setupAuthTestDB(t)

	r := gin.Default()
	routes.SetupAuthRoutes(r)

	// Case 1: Success registration
	payload := models.RegisterRequest{
		Username: "newuser",
		Password: "password123",
		Role:     "Developer",
	}
	bodyBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code 201, got %d", w.Code)
	}

	var returnedUser models.User
	if err := json.Unmarshal(w.Body.Bytes(), &returnedUser); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if returnedUser.Username != payload.Username || returnedUser.Role != payload.Role {
		t.Errorf("Response username or role mismatch")
	}

	// Case 2: Register existing username
	wDup := httptest.NewRecorder()
	reqDup, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(bodyBytes))
	reqDup.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wDup, reqDup)

	if wDup.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code 500 for duplicate user registration, got %d", wDup.Code)
	}
}

func TestAuthLogin(t *testing.T) {
	setupAuthTestDB(t)

	r := gin.Default()
	routes.SetupAuthRoutes(r)

	// Register a test user first
	regPayload := models.RegisterRequest{
		Username: "loginuser",
		Password: "password123",
		Role:     "Viewer",
	}
	regBytes, _ := json.Marshal(regPayload)
	reqReg, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(regBytes))
	reqReg.Header.Set("Content-Type", "application/json")
	wReg := httptest.NewRecorder()
	r.ServeHTTP(wReg, reqReg)

	// Case 1: Success login
	loginPayload := models.LoginRequest{
		Username: "loginuser",
		Password: "password123",
	}
	loginBytes, _ := json.Marshal(loginPayload)
	reqLogin, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(loginBytes))
	reqLogin.Header.Set("Content-Type", "application/json")
	wLogin := httptest.NewRecorder()
	r.ServeHTTP(wLogin, reqLogin)

	if wLogin.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", wLogin.Code)
	}

	var response models.LoginResponse
	if err := json.Unmarshal(wLogin.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Token == "" {
		t.Errorf("Expected token, but got empty string")
	}

	// Case 2: Login with incorrect password
	wrongPayload := models.LoginRequest{
		Username: "loginuser",
		Password: "wrongpassword",
	}
	wrongBytes, _ := json.Marshal(wrongPayload)
	reqWrong, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(wrongBytes))
	reqWrong.Header.Set("Content-Type", "application/json")
	wWrong := httptest.NewRecorder()
	r.ServeHTTP(wWrong, reqWrong)

	if wWrong.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code 401, got %d", wWrong.Code)
	}
}

func TestAuthMe(t *testing.T) {
	setupAuthTestDB(t)

	r := gin.Default()
	routes.SetupAuthRoutes(r)

	token := generateTestToken(t, "Developer")

	// Case 1: Success Fetch
	req, _ := http.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	var user models.User
	if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if user.Username != "test_Developer" || user.Role != "Developer" {
		t.Errorf("Expected user username 'test_Developer' and role 'Developer', got: %v", user)
	}

	// Case 2: Unauthorized
	reqUnauth, _ := http.NewRequest(http.MethodGet, "/api/auth/me", nil)
	wUnauth := httptest.NewRecorder()
	r.ServeHTTP(wUnauth, reqUnauth)

	if wUnauth.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code 401, got %d", wUnauth.Code)
	}
}

func TestAuthLogout(t *testing.T) {
	setupAuthTestDB(t)

	r := gin.Default()
	routes.SetupAuthRoutes(r)

	token := generateTestToken(t, "Viewer")

	req, _ := http.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
}

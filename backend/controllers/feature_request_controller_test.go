package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"backend/config"
	"backend/models"
	"backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupFeatureRequestTestDB(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Migrate related schemas
	err = db.AutoMigrate(&models.Server{}, &models.System{}, &models.FeatureRequest{}, &models.User{})
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	config.DB = db
}

func TestGetFeatureRequestsBySystemID(t *testing.T) {
	setupFeatureRequestTestDB(t)

	// Seed system
	system := models.System{Name: "Core Portal", Type: "Web", Links: "http://core", Status: "Active"}
	config.DB.Create(&system)

	// Seed feature requests
	requests := []models.FeatureRequest{
		{SystemId: system.ID, Title: "Dark Mode", Description: "Implement system wide dark mode", Status: "Pending"},
		{SystemId: system.ID, Title: "Export PDF", Description: "Export page data to PDF", Status: "In Progress"},
	}
	for _, req := range requests {
		config.DB.Create(&req)
	}

	r := gin.Default()
	routes.SetupFeatureRequestRoutes(r)

	token := generateTestToken(t, "Viewer")

	req, _ := http.NewRequest(http.MethodGet, "/api/systems/"+strconv.Itoa(int(system.ID))+"/feature-requests", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	var returnedRequests []models.FeatureRequest
	if err := json.Unmarshal(w.Body.Bytes(), &returnedRequests); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if len(returnedRequests) != 2 {
		t.Errorf("Expected 2 feature requests, but got %d", len(returnedRequests))
	}
}

func TestCreateFeatureRequest(t *testing.T) {
	setupFeatureRequestTestDB(t)

	// Seed system
	system := models.System{Name: "Core Portal", Type: "Web", Links: "http://core", Status: "Active"}
	config.DB.Create(&system)

	r := gin.Default()
	routes.SetupFeatureRequestRoutes(r)

	tokenDev := generateTestToken(t, "Developer")

	// Case 1: Success Creation
	reqBody := models.CreateFeatureRequest{
		Title:       "User Management",
		Description: "Add complete role based access control",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/systems/"+strconv.Itoa(int(system.ID))+"/feature-requests", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenDev)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, but got %d", http.StatusCreated, w.Code)
	}

	// Case 2: Unauthorized (Viewer role trying to POST)
	tokenViewer := generateTestToken(t, "Viewer")
	wUnauth := httptest.NewRecorder()
	reqUnauth, _ := http.NewRequest(http.MethodPost, "/api/systems/"+strconv.Itoa(int(system.ID))+"/feature-requests", bytes.NewBuffer(bodyBytes))
	reqUnauth.Header.Set("Content-Type", "application/json")
	reqUnauth.Header.Set("Authorization", "Bearer "+tokenViewer)
	r.ServeHTTP(wUnauth, reqUnauth)

	if wUnauth.Code != http.StatusForbidden {
		t.Errorf("Expected status code 403 Forbidden for Viewer role, but got %d", wUnauth.Code)
	}
}

func TestUpdateFeatureRequest(t *testing.T) {
	setupFeatureRequestTestDB(t)

	// Seed system & feature request
	system := models.System{Name: "Core Portal"}
	config.DB.Create(&system)

	feature := models.FeatureRequest{SystemId: system.ID, Title: "Old Title", Description: "Old Desc", Status: "Pending"}
	config.DB.Create(&feature)

	r := gin.Default()
	routes.SetupFeatureRequestRoutes(r)

	tokenDev := generateTestToken(t, "Developer")

	// Case 1: Success update status and details
	updatePayload := map[string]interface{}{
		"title":       "New Title",
		"description": "New Description",
		"status":      "Completed",
	}
	bodyBytes, _ := json.Marshal(updatePayload)
	req, _ := http.NewRequest(http.MethodPut, "/api/feature-requests/"+strconv.Itoa(int(feature.ID)), bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenDev)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}
}

func TestDeleteFeatureRequest(t *testing.T) {
	setupFeatureRequestTestDB(t)

	// Seed system & feature request
	system := models.System{Name: "Core Portal"}
	config.DB.Create(&system)

	feature := models.FeatureRequest{SystemId: system.ID, Title: "To Delete", Description: "Desc", Status: "Pending"}
	config.DB.Create(&feature)

	r := gin.Default()
	routes.SetupFeatureRequestRoutes(r)

	tokenDev := generateTestToken(t, "Developer")

	// Case 1: Success deletion
	req, _ := http.NewRequest(http.MethodDelete, "/api/feature-requests/"+strconv.Itoa(int(feature.ID)), nil)
	req.Header.Set("Authorization", "Bearer "+tokenDev)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}
}

func TestGetPendingFeatureRequests(t *testing.T) {
	setupFeatureRequestTestDB(t)

	// Seed system & feature requests
	system := models.System{Name: "Core Portal"}
	config.DB.Create(&system)

	requests := []models.FeatureRequest{
		{SystemId: system.ID, Title: "Pending Request 1", Description: "Desc 1", Status: "Pending"},
		{SystemId: system.ID, Title: "In Progress Request", Description: "Desc 2", Status: "In Progress"},
		{SystemId: system.ID, Title: "Pending Request 2", Description: "Desc 3", Status: "Pending"},
	}
	for _, req := range requests {
		config.DB.Create(&req)
	}

	r := gin.Default()
	routes.SetupFeatureRequestRoutes(r)

	token := generateTestToken(t, "Viewer")

	req, _ := http.NewRequest(http.MethodGet, "/api/feature-requests/pending", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	var response []models.FeatureRequest
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 pending feature requests, but got %d", len(response))
	}

	for _, item := range response {
		if item.Status != "Pending" {
			t.Errorf("Expected status to be 'Pending', but got '%s'", item.Status)
		}
		if item.System == nil || item.System.Name != "Core Portal" {
			t.Errorf("Expected associated system metadata to be preloaded, but got nil or mismatch")
		}
	}
}

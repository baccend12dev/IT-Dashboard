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

func setupDocumentationTestDB(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&models.Server{}, &models.System{}, &models.Documentation{}, &models.User{})
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	config.DB = db
}

func TestGetDocumentationsBySystemID(t *testing.T) {
	setupDocumentationTestDB(t)

	system := models.System{Name: "Portal HRD"}
	config.DB.Create(&system)

	docs := []models.Documentation{
		{SystemId: system.ID, Title: "Flow Bisnis Cuti", Category: "Business Flow", Content: "Diagram flow cuti karyawan"},
		{SystemId: system.ID, Title: "Panduan Deploy", Category: "Deployment Guide", Content: "Langkah build & deploy"},
	}
	for _, doc := range docs {
		config.DB.Create(&doc)
	}

	r := gin.Default()
	routes.SetupDocumentationRoutes(r)

	token := generateTestToken(t, "Viewer")

	// Test case 1: List all for system
	req, _ := http.NewRequest(http.MethodGet, "/api/systems/"+strconv.Itoa(int(system.ID))+"/documentations", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	var returnedDocs []models.Documentation
	json.Unmarshal(w.Body.Bytes(), &returnedDocs)
	if len(returnedDocs) != 2 {
		t.Errorf("Expected 2 documentations, got %d", len(returnedDocs))
	}
}

func TestGetDocumentationByID(t *testing.T) {
	setupDocumentationTestDB(t)

	system := models.System{Name: "Portal HRD"}
	config.DB.Create(&system)

	doc := models.Documentation{SystemId: system.ID, Title: "Database Schema", Category: "Database Documentation", Content: "Table list"}
	config.DB.Create(&doc)

	r := gin.Default()
	routes.SetupDocumentationRoutes(r)

	token := generateTestToken(t, "Viewer")

	// Case 1: Success fetch
	req, _ := http.NewRequest(http.MethodGet, "/api/documentations/"+strconv.Itoa(int(doc.ID)), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestCreateDocumentation(t *testing.T) {
	setupDocumentationTestDB(t)

	system := models.System{Name: "Portal HRD"}
	config.DB.Create(&system)

	r := gin.Default()
	routes.SetupDocumentationRoutes(r)

	tokenDev := generateTestToken(t, "Developer")

	// Case 1: Success creation
	payload := models.CreateDocumentationRequest{
		Title:    "API List",
		Category: "API Documentation",
		Content:  "List of REST APIs",
	}
	bodyBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/systems/"+strconv.Itoa(int(system.ID))+"/documentations", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenDev)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code 201, got %d", w.Code)
	}

	// Case 2: Unauthorized (Viewer role trying to create)
	tokenViewer := generateTestToken(t, "Viewer")
	wUnauth := httptest.NewRecorder()
	reqUnauth, _ := http.NewRequest(http.MethodPost, "/api/systems/"+strconv.Itoa(int(system.ID))+"/documentations", bytes.NewBuffer(bodyBytes))
	reqUnauth.Header.Set("Content-Type", "application/json")
	reqUnauth.Header.Set("Authorization", "Bearer "+tokenViewer)
	r.ServeHTTP(wUnauth, reqUnauth)

	if wUnauth.Code != http.StatusForbidden {
		t.Errorf("Expected status code 403 Forbidden for Viewer role, got %d", wUnauth.Code)
	}
}

func TestUpdateDocumentation(t *testing.T) {
	setupDocumentationTestDB(t)

	system := models.System{Name: "Portal HRD"}
	config.DB.Create(&system)

	doc := models.Documentation{SystemId: system.ID, Title: "Draft Panduan", Category: "User Manual", Content: "Draft content"}
	config.DB.Create(&doc)

	r := gin.Default()
	routes.SetupDocumentationRoutes(r)

	tokenDev := generateTestToken(t, "Developer")

	// Case 1: Success update
	updatePayload := map[string]interface{}{
		"title":    "Panduan Penggunaan Akhir",
		"category": "User Manual",
		"content":  "Final content",
	}
	bodyBytes, _ := json.Marshal(updatePayload)
	req, _ := http.NewRequest(http.MethodPut, "/api/documentations/"+strconv.Itoa(int(doc.ID)), bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenDev)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestDeleteDocumentation(t *testing.T) {
	setupDocumentationTestDB(t)

	system := models.System{Name: "Portal HRD"}
	config.DB.Create(&system)

	doc := models.Documentation{SystemId: system.ID, Title: "To Delete", Category: "Technical Flow", Content: "Temp documentation"}
	config.DB.Create(&doc)

	r := gin.Default()
	routes.SetupDocumentationRoutes(r)

	tokenDev := generateTestToken(t, "Developer")

	// Case 1: Success deletion
	req, _ := http.NewRequest(http.MethodDelete, "/api/documentations/"+strconv.Itoa(int(doc.ID)), nil)
	req.Header.Set("Authorization", "Bearer "+tokenDev)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
}

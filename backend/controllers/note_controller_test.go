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

func setupNoteTestDB(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&models.Server{}, &models.System{}, &models.Note{}, &models.User{})
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	config.DB = db
}

func TestGetNotesBySystemID(t *testing.T) {
	setupNoteTestDB(t)

	system := models.System{Name: "System A"}
	config.DB.Create(&system)

	notes := []models.Note{
		{SystemId: system.ID, Title: "Note 1", Content: "Content 1"},
		{SystemId: system.ID, Title: "Note 2", Content: "Content 2"},
	}
	for _, note := range notes {
		config.DB.Create(&note)
	}

	r := gin.Default()
	routes.SetupNoteRoutes(r)

	token := generateTestToken(t, "Viewer")

	req, _ := http.NewRequest(http.MethodGet, "/api/systems/"+strconv.Itoa(int(system.ID))+"/notes", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	var returnedNotes []models.Note
	json.Unmarshal(w.Body.Bytes(), &returnedNotes)

	if len(returnedNotes) != 2 {
		t.Errorf("Expected 2 notes, got %d", len(returnedNotes))
	}
}

func TestCreateNote(t *testing.T) {
	setupNoteTestDB(t)

	system := models.System{Name: "System A"}
	config.DB.Create(&system)

	r := gin.Default()
	routes.SetupNoteRoutes(r)

	tokenDev := generateTestToken(t, "Developer")

	// Case 1: Success creation
	payload := models.CreateNoteRequest{
		Title:   "New Note",
		Content: "New Content",
	}
	bodyBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/systems/"+strconv.Itoa(int(system.ID))+"/notes", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenDev)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code 201, got %d", w.Code)
	}

	// Case 2: Unauthorized (Viewer role)
	tokenViewer := generateTestToken(t, "Viewer")
	reqUnauth, _ := http.NewRequest(http.MethodPost, "/api/systems/"+strconv.Itoa(int(system.ID))+"/notes", bytes.NewBuffer(bodyBytes))
	reqUnauth.Header.Set("Content-Type", "application/json")
	reqUnauth.Header.Set("Authorization", "Bearer "+tokenViewer)
	wUnauth := httptest.NewRecorder()
	r.ServeHTTP(wUnauth, reqUnauth)

	if wUnauth.Code != http.StatusForbidden {
		t.Errorf("Expected status code 403 Forbidden for Viewer role, got %d", wUnauth.Code)
	}
}

func TestUpdateNote(t *testing.T) {
	setupNoteTestDB(t)

	system := models.System{Name: "System A"}
	config.DB.Create(&system)

	note := models.Note{SystemId: system.ID, Title: "Old Title", Content: "Old Content"}
	config.DB.Create(&note)

	r := gin.Default()
	routes.SetupNoteRoutes(r)

	tokenDev := generateTestToken(t, "Developer")

	// Case 1: Success update
	payload := models.CreateNoteRequest{
		Title:   "Updated Title",
		Content: "Updated Content",
	}
	bodyBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPut, "/api/notes/"+strconv.Itoa(int(note.ID)), bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenDev)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestDeleteNote(t *testing.T) {
	setupNoteTestDB(t)

	system := models.System{Name: "System A"}
	config.DB.Create(&system)

	note := models.Note{SystemId: system.ID, Title: "To Delete", Content: "Content"}
	config.DB.Create(&note)

	r := gin.Default()
	routes.SetupNoteRoutes(r)

	tokenDev := generateTestToken(t, "Developer")

	// Case 1: Success delete
	req, _ := http.NewRequest(http.MethodDelete, "/api/notes/"+strconv.Itoa(int(note.ID)), nil)
	req.Header.Set("Authorization", "Bearer "+tokenDev)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

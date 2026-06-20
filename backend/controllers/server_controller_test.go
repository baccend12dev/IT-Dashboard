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

func setupTestDB(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to in-memory test database: %v", err)
	}

	err = db.AutoMigrate(&models.Server{}, &models.User{})
	if err != nil {
		t.Fatalf("Failed to migrate database schemas: %v", err)
	}

	config.DB = db
}

func TestGetServers(t *testing.T) {
	setupTestDB(t)

	// Seed some test server records
	servers := []models.Server{
		{Name: "Server 1", IP: "192.168.1.10", OS: "Ubuntu 22.04", Location: "Data Center A"},
		{Name: "Server 2", IP: "192.168.1.11", OS: "Rocky Linux 9", Location: "Data Center B"},
	}
	for _, s := range servers {
		config.DB.Create(&s)
	}

	r := gin.Default()
	routes.SetupServerRoutes(r)

	token := generateTestToken(t, "Viewer") // Viewers should be allowed to GET

	req, _ := http.NewRequest(http.MethodGet, "/api/servers/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	var returnedServers []models.Server
	if err := json.Unmarshal(w.Body.Bytes(), &returnedServers); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(returnedServers) != 2 {
		t.Errorf("Expected 2 servers, but got %d", len(returnedServers))
	}
}

func TestGetServerByID(t *testing.T) {
	setupTestDB(t)

	server := models.Server{Name: "Server 1", IP: "192.168.1.10", OS: "Ubuntu 22.04", Location: "Data Center A"}
	config.DB.Create(&server)

	r := gin.Default()
	routes.SetupServerRoutes(r)

	token := generateTestToken(t, "Viewer")

	// Case 1: ID exists
	req, _ := http.NewRequest(http.MethodGet, "/api/servers/"+strconv.Itoa(int(server.ID)), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d for existing server, but got %d", http.StatusOK, w.Code)
	}

	var returnedServer models.Server
	if err := json.Unmarshal(w.Body.Bytes(), &returnedServer); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if returnedServer.ID != server.ID || returnedServer.Name != server.Name {
		t.Errorf("Server details mismatch")
	}

	// Case 2: ID does not exist
	reqNotFound, _ := http.NewRequest(http.MethodGet, "/api/servers/999", nil)
	reqNotFound.Header.Set("Authorization", "Bearer "+token)
	wNotFound := httptest.NewRecorder()
	r.ServeHTTP(wNotFound, reqNotFound)

	if wNotFound.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d for non-existent server, but got %d", http.StatusNotFound, wNotFound.Code)
	}
}

func TestCreateServer(t *testing.T) {
	setupTestDB(t)

	r := gin.Default()
	routes.SetupServerRoutes(r)

	// Case 1: Valid input and authorized (Developer role)
	tokenDev := generateTestToken(t, "Developer")
	reqBody := models.CreateServerRequest{
		Name:     "Test Server",
		IP:       "10.0.0.5",
		OS:       "Debian 12",
		Location: "Rack A1",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/servers/", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenDev)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, but got %d", http.StatusCreated, w.Code)
	}

	// Case 2: Unauthorized role (Viewer trying to create)
	tokenViewer := generateTestToken(t, "Viewer")
	wUnauth := httptest.NewRecorder()
	reqUnauth, _ := http.NewRequest(http.MethodPost, "/api/servers/", bytes.NewBuffer(bodyBytes))
	reqUnauth.Header.Set("Content-Type", "application/json")
	reqUnauth.Header.Set("Authorization", "Bearer "+tokenViewer)
	r.ServeHTTP(wUnauth, reqUnauth)

	if wUnauth.Code != http.StatusForbidden {
		t.Errorf("Expected status code 403 Forbidden for Viewer role, but got %d", wUnauth.Code)
	}
}

func TestUpdateServer(t *testing.T) {
	setupTestDB(t)

	server := models.Server{Name: "Old Name", IP: "192.168.1.10", OS: "Ubuntu 20.04", Location: "Data Center A"}
	config.DB.Create(&server)

	r := gin.Default()
	routes.SetupServerRoutes(r)

	tokenAdmin := generateTestToken(t, "Administrator")

	// Case 1: Update success
	updatedServerData := map[string]interface{}{
		"Name":     "New Name",
		"IP":       "192.168.1.20",
		"OS":       "Ubuntu 22.04",
		"Location": "Data Center B",
	}
	bodyBytes, _ := json.Marshal(updatedServerData)
	req, _ := http.NewRequest(http.MethodPut, "/api/servers/"+strconv.Itoa(int(server.ID)), bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenAdmin)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	var returnedServer models.Server
	if err := json.Unmarshal(w.Body.Bytes(), &returnedServer); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if returnedServer.Name != updatedServerData["Name"].(string) {
		t.Errorf("Returned server does not match updated values")
	}
}

func TestDeleteServer(t *testing.T) {
	setupTestDB(t)

	server := models.Server{Name: "Server To Delete", IP: "192.168.1.10", OS: "Ubuntu 22.04", Location: "Data Center A"}
	config.DB.Create(&server)

	r := gin.Default()
	routes.SetupServerRoutes(r)

	tokenAdmin := generateTestToken(t, "Administrator")

	// Case 1: Delete success
	req, _ := http.NewRequest(http.MethodDelete, "/api/servers/"+strconv.Itoa(int(server.ID)), nil)
	req.Header.Set("Authorization", "Bearer "+tokenAdmin)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}
}

package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetDocumentationsBySystemID retrieves all documentation entries for a specific system (with optional category filter)
// GetDocumentationsBySystemID godoc
// @Summary      Get system documentations
// @Description  Get a list of documentations associated with a system ID, optionally filtered by category
// @Tags         documentations
// @Produce      json
// @Param        id        path      int     true   "System ID"
// @Param        category  query     string  false  "Optional category filter"
// @Success      200  {array}   models.Documentation
// @Failure      500  {object}  map[string]string "error: Gagal mengambil data dokumentasi"
// @Router       /api/systems/{id}/documentations [get]
func GetDocumentationsBySystemID(c *gin.Context) {
	systemID := c.Param("id")
	categoryQuery := c.Query("category")
	var docs []models.Documentation

	query := config.DB.Where("system_id = ?", systemID)
	if categoryQuery != "" {
		query = query.Where("category = ?", categoryQuery)
	}

	query.Find(&docs)

	if config.DB.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data dokumentasi"})
		return
	}

	c.JSON(http.StatusOK, docs)
}

// GetDocumentationByID retrieves a single documentation record by ID
// GetDocumentationByID godoc
// @Summary      Get documentation by ID
// @Description  Get details of a single documentation entry by ID
// @Tags         documentations
// @Produce      json
// @Param        id   path      int  true  "Documentation ID"
// @Success      200  {object}  models.Documentation
// @Failure      404  {object}  map[string]string "error: Dokumentasi tidak ditemukan"
// @Router       /api/documentations/{id} [get]
func GetDocumentationByID(c *gin.Context) {
	var doc models.Documentation
	id := c.Param("id")

	if err := config.DB.First(&doc, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dokumentasi tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, doc)
}

// CreateDocumentation creates a new documentation entry for a system
// CreateDocumentation godoc
// @Summary      Create documentation
// @Description  Create a new documentation entry for a system
// @Tags         documentations
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "System ID"
// @Param        request body   models.CreateDocumentationRequest  true  "Documentation payload"
// @Success      201  {object}  models.Documentation
// @Failure      400  {object}  map[string]string "error: ID System tidak valid / Invalid payload / Kategori tidak valid"
// @Failure      404  {object}  map[string]string "error: System tidak ditemukan"
// @Failure      500  {object}  map[string]string "error: Gagal menyimpan ke database"
// @Router       /api/systems/{id}/documentations [post]
func CreateDocumentation(c *gin.Context) {
	systemIDStr := c.Param("id")
	val, err := strconv.ParseUint(systemIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID System tidak valid"})
		return
	}
	systemID := uint(val)

	// 1. Verify System exists
	var system models.System
	if err := config.DB.First(&system, systemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "System tidak ditemukan. Tidak bisa membuat dokumentasi untuk sistem yang tidak ada.",
		})
		return
	}

	// 2. Bind JSON
	var request models.CreateDocumentationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// 3. Validate Category
	if !models.IsValidDocumentationCategory(request.Category) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Kategori dokumentasi tidak valid. Kategori yang diperbolehkan adalah: 'Business Flow', 'Technical Flow', 'API Documentation', 'Database Documentation', 'Deployment Guide', 'User Manual'",
		})
		return
	}

	doc := models.Documentation{
		SystemId: systemID,
		Title:    request.Title,
		Category: request.Category,
		Content:  request.Content,
	}

	// 4. Save to database
	if err := config.DB.Create(&doc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan ke database"})
		return
	}

	c.JSON(http.StatusCreated, doc)
}

// UpdateDocumentation updates a documentation entry
// UpdateDocumentation godoc
// @Summary      Update documentation
// @Description  Update details of an existing documentation entry by ID
// @Tags         documentations
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Documentation ID"
// @Param        request body   models.UpdateDocumentationRequest  true  "Documentation update fields"
// @Success      200  {object}  models.Documentation
// @Failure      400  {object}  map[string]string "error: Invalid payload / Kategori tidak valid"
// @Failure      404  {object}  map[string]string "error: Dokumentasi tidak ditemukan"
// @Failure      500  {object}  map[string]string "error: Gagal memperbarui database"
// @Router       /api/documentations/{id} [put]
func UpdateDocumentation(c *gin.Context) {
	var doc models.Documentation
	id := c.Param("id")

	// 1. Verify existence
	if err := config.DB.First(&doc, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dokumentasi tidak ditemukan"})
		return
	}

	// 2. Bind payload
	var request models.UpdateDocumentationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Validate category if provided
	if request.Category != "" && !models.IsValidDocumentationCategory(request.Category) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Kategori dokumentasi tidak valid. Kategori yang diperbolehkan adalah: 'Business Flow', 'Technical Flow', 'API Documentation', 'Database Documentation', 'Deployment Guide', 'User Manual'",
		})
		return
	}

	// 4. Apply updates
	updates := make(map[string]interface{})
	if request.Title != "" {
		updates["title"] = request.Title
	}
	if request.Category != "" {
		updates["category"] = request.Category
	}
	if request.Content != "" {
		updates["content"] = request.Content
	}

	if len(updates) > 0 {
		if err := config.DB.Model(&doc).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui database"})
			return
		}
	}

	// Fetch updated record
	config.DB.First(&doc, id)
	c.JSON(http.StatusOK, doc)
}

// DeleteDocumentation deletes a documentation entry
// DeleteDocumentation godoc
// @Summary      Delete documentation
// @Description  Delete a documentation entry from database by ID
// @Tags         documentations
// @Produce      json
// @Param        id   path      int  true  "Documentation ID"
// @Success      200  {object}  map[string]string "message: Dokumentasi deleted successfully"
// @Failure      404  {object}  map[string]string "error: Dokumentasi tidak ditemukan"
// @Failure      500  {object}  map[string]string "error: Gagal menghapus dokumentasi"
// @Router       /api/documentations/{id} [delete]
func DeleteDocumentation(c *gin.Context) {
	var doc models.Documentation
	id := c.Param("id")

	// 1. Verify existence
	if err := config.DB.First(&doc, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dokumentasi tidak ditemukan"})
		return
	}

	// 2. Delete record
	if err := config.DB.Delete(&doc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus dokumentasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Dokumentasi deleted successfully",
		"id":      id,
	})
}

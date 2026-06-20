package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetFeatureRequestsBySystemID retrieves all feature requests for a specific system
// GetFeatureRequestsBySystemID godoc
// @Summary      Get feature requests by system ID
// @Description  Get a list of feature requests associated with a specific system
// @Tags         feature-requests
// @Produce      json
// @Param        id   path      int  true  "System ID"
// @Success      200  {array}   models.FeatureRequest
// @Failure      500  {object}  map[string]string "error: Gagal mengambil data feature request"
// @Router       /api/systems/{id}/feature-requests [get]
func GetFeatureRequestsBySystemID(c *gin.Context) {
	systemID := c.Param("id")
	var featureRequests []models.FeatureRequest

	config.DB.Where("system_id = ?", systemID).Find(&featureRequests)

	// Check if database error occurred
	if config.DB.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data feature request"})
		return
	}

	c.JSON(http.StatusOK, featureRequests)
}

// CreateFeatureRequest creates a new feature request for a system
// CreateFeatureRequest godoc
// @Summary      Create a new feature request
// @Description  Create a new feature request for a system
// @Tags         feature-requests
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "System ID"
// @Param        request body   models.CreateFeatureRequest  true  "Feature request payload"
// @Success      201  {object}  models.FeatureRequest
// @Failure      400  {object}  map[string]string "error: ID System tidak valid / Invalid request body"
// @Failure      404  {object}  map[string]string "error: System tidak ditemukan"
// @Failure      500  {object}  map[string]string "error: Gagal menyimpan ke database"
// @Router       /api/systems/{id}/feature-requests [post]
func CreateFeatureRequest(c *gin.Context) {
	systemIDStr := c.Param("id")
	val, err := strconv.ParseUint(systemIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID System tidak valid"})
		return
	}
	systemID := uint(val)

	// 1. Check if System exists
	var system models.System
	if err := config.DB.First(&system, systemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "System tidak ditemukan. Tidak bisa membuat request fitur untuk sistem yang tidak ada.",
		})
		return
	}

	// 2. Bind request body
	var request models.CreateFeatureRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	featureRequest := models.FeatureRequest{
		SystemId:    systemID,
		Title:       request.Title,
		Description: request.Description,
		Status:      "Pending", // Default status is Pending
	}

	// 3. Save to database
	if err := config.DB.Create(&featureRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan ke database"})
		return
	}

	c.JSON(http.StatusCreated, featureRequest)
}

// UpdateFeatureRequest updates a feature request details (title, description, status)
// UpdateFeatureRequest godoc
// @Summary      Update a feature request
// @Description  Update details/status of an existing feature request by feature request ID
// @Tags         feature-requests
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Feature Request ID"
// @Param        request body   models.UpdateFeatureRequest  true  "Feature request update payload"
// @Success      200  {object}  models.FeatureRequest
// @Failure      400  {object}  map[string]string "error: Bad request"
// @Failure      404  {object}  map[string]string "error: Feature Request tidak ditemukan"
// @Failure      500  {object}  map[string]string "error: Gagal memperbarui database"
// @Router       /api/feature-requests/{id} [put]
func UpdateFeatureRequest(c *gin.Context) {
	var featureRequest models.FeatureRequest
	id := c.Param("id")

	// 1. Verify existence
	if err := config.DB.First(&featureRequest, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Feature Request tidak ditemukan"})
		return
	}

	// 2. Bind payload
	var request models.UpdateFeatureRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Prepare partial updates
	updates := make(map[string]interface{})
	if request.Title != "" {
		updates["title"] = request.Title
	}
	if request.Description != "" {
		updates["description"] = request.Description
	}
	if request.Status != "" {
		updates["status"] = request.Status
	}

	if len(updates) > 0 {
		if err := config.DB.Model(&featureRequest).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui database"})
			return
		}
	}

	// Fetch updated record
	config.DB.First(&featureRequest, id)
	c.JSON(http.StatusOK, featureRequest)
}

// DeleteFeatureRequest deletes a feature request by ID
// DeleteFeatureRequest godoc
// @Summary      Delete a feature request
// @Description  Delete a feature request record from database by feature request ID
// @Tags         feature-requests
// @Produce      json
// @Param        id   path      int  true  "Feature Request ID"
// @Success      200  {object}  map[string]string "message: Feature Request deleted successfully"
// @Failure      404  {object}  map[string]string "error: Feature Request tidak ditemukan"
// @Failure      500  {object}  map[string]string "error: Gagal menghapus data"
// @Router       /api/feature-requests/{id} [delete]
func DeleteFeatureRequest(c *gin.Context) {
	var featureRequest models.FeatureRequest
	id := c.Param("id")

	// 1. Check if exists
	if err := config.DB.First(&featureRequest, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Feature Request tidak ditemukan"})
		return
	}

	// 2. Delete record
	if err := config.DB.Delete(&featureRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Feature Request deleted successfully",
		"id":      id,
	})
}

// GetPendingFeatureRequests retrieves all feature requests with 'Pending' status across all systems
// GetPendingFeatureRequests godoc
// @Summary      Get all pending feature requests
// @Description  Get a list of all pending feature requests preloaded with system info
// @Tags         feature-requests
// @Produce      json
// @Success      200  {array}   models.FeatureRequest
// @Failure      500  {object}  map[string]string "error: Gagal mengambil data pending requests"
// @Router       /api/feature-requests/pending [get]
func GetPendingFeatureRequests(c *gin.Context) {
	var featureRequests []models.FeatureRequest

	// Fetch all where status is 'Pending' and preload their associated System
	if err := config.DB.Preload("System").Where("status = ?", "Pending").Find(&featureRequests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pending requests"})
		return
	}

	c.JSON(http.StatusOK, featureRequests)
}

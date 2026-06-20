package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSystems godoc
// @Summary Get all systems
// @Description Get a list of all systems
// GetSystems retrieves all systems from the database and sends them as a JSON response.
// It queries the database for all system records and returns them to the client.
// If no systems are found, it returns an empty list.
// If a database error occurs, it returns an appropriate error response.
// GetSystems godoc
// @Summary      Get all systems
// @Description  Get a list of all systems
// @Tags         systems
// @Produce      json
// @Success      200  {array}   models.System
// @Router       /api/systems [get]
func GetSystems(c *gin.Context) {
	var systems []models.System // Assuming you have a System model defined in your models package
	config.DB.Find(&systems)

	c.JSON(http.StatusOK, systems)
}

// GetSystemByID godoc
// @Summary      Get system by ID
// @Description  Get detailed information of a system including server
// @Tags         systems
// @Produce      json
// @Param        id   path      int  true  "System ID"
// @Success      200  {object}  models.System
// @Failure      404  {object}  map[string]string "error: Id not found"
// @Router       /api/systems/{id} [get]
func GetSystemByID(c *gin.Context) {
	var system models.System
	id := c.Param("id")
	if err := config.DB.Preload("Server").First(&system, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Id not found"})
		return
	}
	c.JSON(http.StatusOK, system)
}

// // create system database entry simple input example
// func CreateSystem(c *gin.Context) {
// 	var system models.System

// 	//validasi request body
// 	if err := c.ShouldBindJSON(&system); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// simpan data ke database
// 	config.DB.Create(&system)
// 	c.JSON(http.StatusCreated, system)

// }

// recomended create system database entry with request body struct
// CreateSystem godoc
// @Summary      Create a new system
// @Description  Create a new system database entry by providing system fields
// @Tags         systems
// @Accept       json
// @Produce      json
// @Param        request body   models.CreateSystemRequest  true  "System creation payload"
// @Success      201  {object}  models.System
// @Failure      400  {object}  map[string]string "error: Invalid request body"
// @Failure      404  {object}  map[string]string "error: Server not found"
// @Failure      500  {object}  map[string]string "error: Failed to create system"
// @Router       /api/systems [post]
func CreateSystem(c *gin.Context) {
	var request models.CreateSystemRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	// cek apakah server dengan ID yang diberikan ada di database
	var server models.Server
	if err := config.DB.First(&server, request.ServerId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Server not found",
			"details": "Server dengan ID tersebut tidak tersedia di database",
		})
		return
	}

	//create system sesuai dengan request body struct
	system := models.System{
		Name:        request.Name,
		Type:        request.Type,
		Links:       request.Links,
		ServerId:    request.ServerId,
		Status:      request.Status,
		Description: request.Description,
	}

	// simpan data ke database
	if err := config.DB.Create(&system).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create system",
			"details": err.Error(),
		})
		return
	}
	// return response dengan data system yang baru dibuat
	c.JSON(http.StatusCreated, gin.H{
		"message": "System created successfully",
		"system":  system,
	})

}

// update function for system
// UpdateSystem godoc
// @Summary      Update a system
// @Description  Update details of an existing system by ID
// @Tags         systems
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "System ID"
// @Param        request body   models.System  true  "Fields to update"
// @Success      200  {object}  models.System
// @Failure      400  {object}  map[string]string "error: Bad request"
// @Failure      404  {object}  map[string]string "error: Id not found"
// @Failure      500  {object}  map[string]string "error: Failed to update database"
// @Router       /api/systems/{id} [put]
func UpdateSystem(c *gin.Context) {
	var system models.System
	id := c.Param("id")

	// 1. Cari data lama berdasarkan ID
	// Gunakan .First() untuk memastikan data memang ada di DB
	if err := config.DB.First(&system, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Id not found"})
		return
	}

	// 2. Bind JSON dari request body ke struct system
	// Ini akan menimpa field yang dikirimkan di JSON
	if err := c.ShouldBindJSON(&system); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Simpan perubahan (Gunakan Updates, bukan Save)
	// .Model(&system) memastikan GORM tahu record mana yang diupdate (berdasarkan ID)
	if err := config.DB.Model(&system).Updates(&system).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update database"})
		return
	}

	// 4. Kembalikan data yang sudah diperbarui
	c.JSON(http.StatusOK, system)
}

// DeleteSystem is a handler function that deletes a system from the database based on the provided ID.
// It retrieves the system record using the ID from the URL parameters, and if found, deletes it from the database.
// If the system is not found, it returns a 404 Not Found response. If the deletion is successful, it returns a 204 No Content response.

// DeleteSystem godoc
// @Summary      Delete a system
// @Description  Delete a system record and its related notes/features from database
// @Tags         systems
// @Produce      json
// @Param        id   path      int  true  "System ID"
// @Success      200  {object}  map[string]string "message: Delete successfully"
// @Failure      404  {object}  map[string]string "error: Id not found"
// @Failure      500  {object}  map[string]string "error: Failed to delete data"
// @Router       /api/systems/{id} [delete]
func DeleteSystem(c *gin.Context) {
	var system models.System
	id := c.Param("id")

	// 1. Cari data terlebih dahulu
	if err := config.DB.First(&system, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Id not found"})
		return
	}

	// 2. Hapus data
	if err := config.DB.Delete(&system).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete data"})
		return
	}

	// 3. Kembalikan response sukses dengan pesan
	c.JSON(http.StatusOK, gin.H{
		"message": "Delete successfully",
		"id":      id,
	})
}

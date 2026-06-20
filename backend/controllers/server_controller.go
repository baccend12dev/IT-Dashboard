package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetServers godoc
// @Summary Get all servers
// GetServers godoc
// @Summary      Get all servers
// @Description  Get a list of all server records
// @Tags         servers
// @Produce      json
// @Success      200  {array}   models.Server
// @Router       /api/servers [get]
func GetServers(c *gin.Context) {
	var servers []models.Server
	config.DB.Find(&servers)
	c.JSON(http.StatusOK, servers)
}

// GetServerByID godoc
// @Summary      Get server by ID
// @Description  Get details of a server record by ID
// @Tags         servers
// @Produce      json
// @Param        id   path      int  true  "Server ID"
// @Success      200  {object}  models.Server
// @Failure      404  {object}  map[string]string "error: Id not found"
// @Router       /api/servers/{id} [get]
func GetServerByID(c *gin.Context) {
	var server models.Server
	id := c.Param("id")
	if err := config.DB.First(&server, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Id not found"})
		return
	}
	c.JSON(http.StatusOK, server)
}

// CreateServer godoc
// @Summary      Create a new server
// @Description  Create a new server record
// @Tags         servers
// @Accept       json
// @Produce      json
// @Param        request body   models.CreateServerRequest  true  "Server payload"
// @Success      201  {object}  models.Server
// @Failure      400  {object}  map[string]string "error: Invalid request body"
// @Router       /api/servers [post]
func CreateServer(c *gin.Context) {
	var request models.CreateServerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	server := models.Server{
		Name:     request.Name,
		IP:       request.IP,
		Location: request.Location,
		OS:       request.OS,
	}
	config.DB.Create(&server)
	c.JSON(http.StatusCreated, server)
}

// update server
// UpdateServer godoc
// @Summary      Update a server
// @Description  Update details of an existing server record
// @Tags         servers
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Server ID"
// @Param        request body   models.Server  true  "Fields to update"
// @Success      200  {object}  models.Server
// @Failure      400  {object}  map[string]string "error: Bad request"
// @Failure      404  {object}  map[string]string "error: Id not found"
// @Router       /api/servers/{id} [put]
func UpdateServer(c *gin.Context) {
	var server models.Server
	id := c.Param("id")
	if err := config.DB.First(&server, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Id not found"})
		return
	}
	if err := c.ShouldBindJSON(&server); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&server)
	c.JSON(http.StatusOK, server)
}

// Deleted server function, because we don't want to delete server data, we can just update the status to "inactive" or something like that. If you want to implement delete server function, you can uncomment the code below.
// DeleteServer godoc
// @Summary      Delete a server
// @Description  Delete a server record from database
// @Tags         servers
// @Produce      json
// @Param        id   path      int  true  "Server ID"
// @Success      200  {object}  map[string]string "message: Server deleted successfully"
// @Failure      404  {object}  map[string]string "error: Id not found"
// @Router       /api/servers/{id} [delete]
func DeleteServer(c *gin.Context) {
	var server models.Server
	id := c.Param("id")
	if err := config.DB.First(&server, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Id not found"})
		return
	}
	config.DB.Delete(&server)
	c.JSON(http.StatusOK, gin.H{"message": "Server deleted successfully"})
}

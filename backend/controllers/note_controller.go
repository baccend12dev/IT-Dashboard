package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetNotesBySystemID godoc
// kemungkinan tidak diguankan karena sudah ada di system controller
// tapi tetap dibuat jika note banyak untuk bisa filter dan paginate
// GetNotesBySystemID godoc
// @Summary      Get notes by system ID
// @Description  Get a list of notes associated with a specific system
// @Tags         notes
// @Produce      json
// @Param        id   path      int  true  "System ID"
// @Success      200  {array}   models.Note
// @Failure      500  {object}  map[string]string "error: Gagal mengambil catatan"
// @Router       /api/systems/{id}/notes [get]
func GetNotesBySystemID(c *gin.Context) {
	systemID := c.Param("id")
	var notes []models.Note
	config.DB.Where("system_id = ?", systemID).Find(&notes)

	// Cek jika terjadi error pada database
	if config.DB.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil catatan"})
		return
	}
	// Return the notes as JSON response
	c.JSON(http.StatusOK, notes)
}

// CreateNote godoc
// @Summary Create a new note for a system
// @Description Create a new note for a system by providing the system ID and note content
// CreateNote godoc
// @Summary      Create a new note
// @Description  Create a new note for a system
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "System ID"
// @Param        request body   models.CreateNoteRequest  true  "Note content payload"
// @Success      201  {object}  models.Note
// @Failure      400  {object}  map[string]string "error: ID System tidak valid / Invalid request body"
// @Failure      404  {object}  map[string]string "error: System tidak ditemukan"
// @Failure      500  {object}  map[string]string "error: Gagal menyimpan ke database"
// @Router       /api/systems/{id}/notes [post]
func CreateNote(c *gin.Context) {
	systemIDStr := c.Param("id")
	// Konversi string ke uint64 dulu, lalu ke uint
	val, err := strconv.ParseUint(systemIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID System tidak valid"})
		return
	}
	systemID := uint(val)

	// 1. CEK DULU: Apakah System dengan ID ini ada di database?
	var system models.System
	if err := config.DB.First(&system, systemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "System tidak ditemukan. Tidak bisa membuat catatan untuk sistem yang tidak ada.",
		})
		return
	}

	var request models.CreateNoteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	note := models.Note{
		SystemId: systemID,
		Title:    request.Title,
		Content:  request.Content,
	}

	// Sebaiknya tangkap error jika database gagal simpan
	if err := config.DB.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan ke database"})
		return
	}

	c.JSON(http.StatusCreated, note)
}

// UpdateNote updates a note's title and content by ID
// UpdateNote godoc
// @Summary      Update a note
// @Description  Update details of an existing note record by note ID
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Note ID"
// @Param        request body   models.CreateNoteRequest  true  "Note content payload"
// @Success      200  {object}  models.Note
// @Failure      400  {object}  map[string]string "error: Invalid request body"
// @Failure      404  {object}  map[string]string "error: Catatan tidak ditemukan"
// @Failure      500  {object}  map[string]string "error: Gagal memperbarui catatan di database"
// @Router       /api/notes/{id} [put]
func UpdateNote(c *gin.Context) {
	var note models.Note
	id := c.Param("id")

	// 1. Verify existence
	if err := config.DB.First(&note, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Catatan tidak ditemukan"})
		return
	}

	// 2. Bind request body
	var request models.CreateNoteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	note.Title = request.Title
	note.Content = request.Content

	if err := config.DB.Save(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui catatan di database"})
		return
	}

	c.JSON(http.StatusOK, note)
}

// DeleteNote deletes a note by ID
// DeleteNote godoc
// @Summary      Delete a note
// @Description  Delete a note record from database by note ID
// @Tags         notes
// @Produce      json
// @Param        id   path      int  true  "Note ID"
// @Success      200  {object}  map[string]string "message: Catatan deleted successfully"
// @Failure      404  {object}  map[string]string "error: Catatan tidak ditemukan"
// @Failure      500  {object}  map[string]string "error: Gagal menghapus catatan"
// @Router       /api/notes/{id} [delete]
func DeleteNote(c *gin.Context) {
	var note models.Note
	id := c.Param("id")

	// 1. Check if exists
	if err := config.DB.First(&note, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Catatan tidak ditemukan"})
		return
	}

	// 2. Delete record
	if err := config.DB.Delete(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus catatan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Catatan deleted successfully",
		"id":      id,
	})
}


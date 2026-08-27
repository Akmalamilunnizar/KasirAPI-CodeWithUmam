package handlers

import (
	"net/http"

	"kasirApi/config"
	"kasirApi/models"

	"github.com/gin-gonic/gin"
)

type KontakHandler struct{}

func NewKontakHandler() *KontakHandler {
	return &KontakHandler{}
}

// GetAll retrieves contacts with optional ?tipe=SUPPLIER|CUSTOMER & ?search=
func (h *KontakHandler) GetAll(c *gin.Context) {
	var kontaks []models.Kontak
	query := config.DB.Model(&models.Kontak{})

	if tipe := c.Query("tipe"); tipe != "" {
		query = query.Where("tipe = ?", tipe)
	}

	if entitas := c.Query("jenis_entitas"); entitas != "" {
		query = query.Where("jenis_entitas = ?", entitas)
	}

	if search := c.Query("search"); search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("nama_kontak LIKE ? OR kode_kontak LIKE ?", searchTerm, searchTerm)
	}

	if err := query.Order("id ASC").Find(&kontaks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": kontaks, "total": len(kontaks)})
}

// Create creates a new contact
func (h *KontakHandler) Create(c *gin.Context) {
	var input models.Kontak
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.JenisEntitas == "" {
		input.JenisEntitas = models.EntitasIndividu
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": input, "message": "Kontak berhasil dibuat"})
}

// Update updates an existing contact
func (h *KontakHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var kontak models.Kontak
	if err := config.DB.First(&kontak, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kontak tidak ditemukan"})
		return
	}

	var input models.Kontak
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kontak.NamaKontak = input.NamaKontak
	kontak.Tipe = input.Tipe
	kontak.JenisEntitas = input.JenisEntitas
	kontak.NoTelp = input.NoTelp
	kontak.Alamat = input.Alamat
	if input.KodeKontak != "" {
		kontak.KodeKontak = input.KodeKontak
	}

	if err := config.DB.Save(&kontak).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": kontak, "message": "Kontak berhasil diupdate"})
}

// Delete soft deletes a contact
func (h *KontakHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	var kontak models.Kontak
	if err := config.DB.First(&kontak, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kontak tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&kontak).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kontak berhasil dihapus"})
}

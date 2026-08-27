package handlers

import (
	"net/http"
	"strconv"

	"kasirApi/config"
	"kasirApi/models"

	"github.com/gin-gonic/gin"
)

type BarangHandler struct{}

func NewBarangHandler() *BarangHandler {
	return &BarangHandler{}
}

// GetAll retrieves all products with search & category filters and active stock calculation
func (h *BarangHandler) GetAll(c *gin.Context) {
	var barangs []models.Barang
	query := config.DB.Model(&models.Barang{})

	if kategori := c.Query("kategori"); kategori != "" {
		query = query.Where("kategori = ?", kategori)
	}

	if search := c.Query("search"); search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("nama_barang LIKE ? OR kode_barang LIKE ?", searchTerm, searchTerm)
	}

	if err := query.Order("id ASC").Find(&barangs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch barangs: " + err.Error()})
		return
	}

	// Compute total active stock per barang
	for i := range barangs {
		var totalStok int64
		config.DB.Model(&models.StokBatch{}).
			Where("barang_id = ? AND total_qty > 0", barangs[i].ID).
			Select("COALESCE(SUM(total_qty), 0)").
			Scan(&totalStok)

		barangs[i].TotalStok = int(totalStok)
		if barangs[i].PcsPerBox > 0 {
			barangs[i].TotalBox = int(totalStok) / barangs[i].PcsPerBox
			barangs[i].SisaPcs = int(totalStok) % barangs[i].PcsPerBox
		}

		if totalStok == 0 {
			barangs[i].StokStatus = "OUT_OF_STOCK"
		} else if int(totalStok) < barangs[i].MinStok {
			barangs[i].StokStatus = "CRITICAL"
		} else if int(totalStok) == barangs[i].MinStok {
			barangs[i].StokStatus = "WARNING"
		} else {
			barangs[i].StokStatus = "NORMAL"
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  barangs,
		"total": len(barangs),
	})
}

// GetByID retrieves a single product by ID
func (h *BarangHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	var barang models.Barang
	if err := config.DB.First(&barang, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Barang not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": barang})
}

// Create handles creating a new barang
func (h *BarangHandler) Create(c *gin.Context) {
	var input models.Barang
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.PcsPerBox <= 0 {
		input.PcsPerBox = 12
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create barang: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": input, "message": "Barang created successfully"})
}

// Update handles updating an existing barang
func (h *BarangHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var barang models.Barang
	if err := config.DB.First(&barang, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Barang not found"})
		return
	}

	var input models.Barang
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	barang.NamaBarang = input.NamaBarang
	barang.Kategori = input.Kategori
	barang.MinStok = input.MinStok
	barang.PcsPerBox = input.PcsPerBox
	if input.KodeBarang != "" {
		barang.KodeBarang = input.KodeBarang
	}

	if err := config.DB.Save(&barang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update barang: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": barang, "message": "Barang updated successfully"})
}

// Delete handles soft deleting a barang
func (h *BarangHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var barang models.Barang
	if err := config.DB.First(&barang, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Barang not found"})
		return
	}

	if err := config.DB.Delete(&barang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete barang: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Barang deleted successfully (Soft Deleted)"})
}

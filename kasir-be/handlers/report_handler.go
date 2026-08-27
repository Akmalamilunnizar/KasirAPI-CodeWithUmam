package handlers

import (
	"math"
	"net/http"
	"time"

	"kasirApi/config"
	"kasirApi/models"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct{}

func NewReportHandler() *ReportHandler {
	return &ReportHandler{}
}

// GetTopBarang returns top 3 best-selling items
func (h *ReportHandler) GetTopBarang(c *gin.Context) {
	type Result struct {
		BarangID     uint   `json:"barang_id"`
		KodeBarang   string `json:"kode_barang"`
		NamaBarang   string `json:"nama_barang"`
		Kategori     string `json:"kategori"`
		TotalTerjual int    `json:"total_terjual"`
	}

	results := make([]Result, 0)

	err := config.DB.Table("transaksi_stoks").
		Select("barangs.id as barang_id, barangs.kode_barang, barangs.nama_barang, barangs.kategori, COALESCE(SUM(transaksi_stoks.qty), 0) as total_terjual").
		Joins("JOIN barangs ON barangs.id = transaksi_stoks.barang_id").
		Where("transaksi_stoks.jenis_transaksi = ?", "OUT").
		Where("barangs.deleted_at IS NULL").
		Group("barangs.id, barangs.kode_barang, barangs.nama_barang, barangs.kategori").
		Order("total_terjual DESC").
		Limit(3).
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "data": results})
}

// GetTrackingStok returns mutation audit trail logs with Kontak (Party) details
func (h *ReportHandler) GetTrackingStok(c *gin.Context) {
	logs := make([]models.TransaksiStok, 0)
	query := config.DB.Preload("Barang").Preload("Kontak").Order("tanggal_transaksi DESC")

	if barangID := c.Query("barang_id"); barangID != "" {
		query = query.Where("barang_id = ?", barangID)
	}
	if kontakID := c.Query("kontak_id"); kontakID != "" {
		query = query.Where("kontak_id = ?", kontakID)
	}
	if jenis := c.Query("jenis_transaksi"); jenis != "" {
		query = query.Where("jenis_transaksi = ?", jenis)
	}

	if err := query.Limit(100).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"data":   logs,
		"total":  len(logs),
	})
}

// GetStokPerBarang returns active batches with FEFO status and supplier party details
func (h *ReportHandler) GetStokPerBarang(c *gin.Context) {
	var batches []models.StokBatch
	err := config.DB.Preload("Barang").Preload("Kontak").
		Where("total_qty > 0").
		Order("expired_date ASC, barang_id ASC").
		Find(&batches).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	response := make([]models.StokPerBarangResponse, 0)

	for _, b := range batches {
		if b.Barang == nil {
			continue
		}

		pcsPerBox := b.Barang.PcsPerBox
		if pcsPerBox <= 0 {
			pcsPerBox = 12
		}

		box := b.TotalQty / pcsPerBox
		pcs := b.TotalQty % pcsPerBox

		diffHours := b.ExpiredDate.Sub(now).Hours()
		daysUntilED := int(diffHours / 24)
		isNearED := daysUntilED <= 90

		supplierName := "-"
		jenisEntitas := "-"
		if b.Kontak != nil {
			supplierName = b.Kontak.NamaKontak
			jenisEntitas = b.Kontak.JenisEntitas
		}

		response = append(response, models.StokPerBarangResponse{
			BatchID:      b.ID,
			BarangID:     b.BarangID,
			KodeBarang:   b.Barang.KodeBarang,
			NamaBarang:   b.Barang.NamaBarang,
			Kategori:     b.Barang.Kategori,
			PcsPerBox:    pcsPerBox,
			NamaSupplier: supplierName,
			JenisEntitas: jenisEntitas,
			ExpiredDate:  b.ExpiredDate,
			TotalQtyPcs:  b.TotalQty,
			Box:          box,
			Pcs:          pcs,
			IsNearED:     isNearED,
			DaysUntilED:  daysUntilED,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"data":   response,
		"total":  len(response),
	})
}

// GetLowStockAlerts returns low stock alerts
func (h *ReportHandler) GetLowStockAlerts(c *gin.Context) {
	var barangs []models.Barang
	if err := config.DB.Where("deleted_at IS NULL").Find(&barangs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	alerts := make([]models.LowStockAlertResponse, 0)

	for _, b := range barangs {
		var currentStok int64
		config.DB.Model(&models.StokBatch{}).
			Where("barang_id = ? AND total_qty > 0", b.ID).
			Select("COALESCE(SUM(total_qty), 0)").
			Scan(&currentStok)

		curStokInt := int(currentStok)
		if curStokInt <= b.MinStok {
			status := "WARNING"
			if curStokInt == 0 {
				status = "OUT_OF_STOCK"
			} else if curStokInt < int(math.Ceil(float64(b.MinStok)/2.0)) {
				status = "CRITICAL"
			}

			alerts = append(alerts, models.LowStockAlertResponse{
				BarangID:    b.ID,
				KodeBarang:  b.KodeBarang,
				NamaBarang:  b.NamaBarang,
				Kategori:    b.Kategori,
				MinStok:     b.MinStok,
				CurrentStok: curStokInt,
				Status:      status,
				Deficit:     b.MinStok - curStokInt,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "SUCCESS",
		"total_alert": len(alerts),
		"data":        alerts,
	})
}

package handlers

import (
	"fmt"
	"net/http"
	"time"

	"kasirApi/config"
	"kasirApi/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransaksiHandler struct{}

func NewTransaksiHandler() *TransaksiHandler {
	return &TransaksiHandler{}
}

// HandleMutasi executes stock mutation (IN with Supplier Kontak, OUT with Customer Kontak via FEFO)
func (h *TransaksiHandler) HandleMutasi(c *gin.Context) {
	var req models.MutasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var barang models.Barang
	if err := config.DB.First(&barang, req.BarangID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Barang tidak ditemukan"})
		return
	}

	var kontak models.Kontak
	if err := config.DB.First(&kontak, req.KontakID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kontak (Supplier/Customer) tidak ditemukan"})
		return
	}

	var resultTransaksi models.TransaksiStok
	var affectedBatches []gin.H

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		if req.JenisTransaksi == "IN" {
			if req.ExpiredDate == "" {
				return fmt.Errorf("expired_date wajib diisi untuk transaksi IN (barang masuk)")
			}

			ed, err := time.Parse("2006-01-02", req.ExpiredDate)
			if err != nil {
				ed, err = time.Parse(time.RFC3339, req.ExpiredDate)
				if err != nil {
					return fmt.Errorf("format expired_date tidak valid (gunakan YYYY-MM-DD): %w", err)
				}
			}

			batch := models.StokBatch{
				BarangID:    req.BarangID,
				KontakID:    &req.KontakID,
				ExpiredDate: ed,
				TotalQty:    req.Qty,
				CreatedAt:   now,
			}
			if err := tx.Create(&batch).Error; err != nil {
				return fmt.Errorf("gagal membuat batch stok: %w", err)
			}

			transaksi := models.TransaksiStok{
				BarangID:         req.BarangID,
				KontakID:         req.KontakID,
				TanggalTransaksi: now,
				JenisTransaksi:   "IN",
				Qty:              req.Qty,
				Keterangan:       req.Keterangan,
				CreatedAt:        now,
			}
			if err := tx.Create(&transaksi).Error; err != nil {
				return fmt.Errorf("gagal menyimpan transaksi log: %w", err)
			}

			resultTransaksi = transaksi
			affectedBatches = append(affectedBatches, gin.H{
				"batch_id":     batch.ID,
				"expired_date": batch.ExpiredDate.Format("2006-01-02"),
				"qty_added":    req.Qty,
			})

		} else if req.JenisTransaksi == "OUT" {
			// FEFO Greedy Deduction
			var batches []models.StokBatch
			err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("barang_id = ? AND total_qty > 0", req.BarangID).
				Order("expired_date ASC, id ASC").
				Find(&batches).Error

			if err != nil {
				return fmt.Errorf("gagal membaca batch stok: %w", err)
			}

			totalAvailable := 0
			for _, b := range batches {
				totalAvailable += b.TotalQty
			}

			if totalAvailable < req.Qty {
				return fmt.Errorf("stok tidak mencukupi! Tersedia: %d pcs, Diminta: %d pcs", totalAvailable, req.Qty)
			}

			sisaDibutuhkan := req.Qty
			for i := range batches {
				if sisaDibutuhkan <= 0 {
					break
				}

				batch := &batches[i]
				qtyPotong := 0

				if batch.TotalQty >= sisaDibutuhkan {
					qtyPotong = sisaDibutuhkan
					batch.TotalQty -= sisaDibutuhkan
					sisaDibutuhkan = 0
				} else {
					qtyPotong = batch.TotalQty
					sisaDibutuhkan -= batch.TotalQty
					batch.TotalQty = 0
				}

				if err := tx.Model(&models.StokBatch{}).
					Where("id = ?", batch.ID).
					Update("total_qty", batch.TotalQty).Error; err != nil {
					return fmt.Errorf("gagal mengupdate batch #%d: %w", batch.ID, err)
				}

				affectedBatches = append(affectedBatches, gin.H{
					"batch_id":     batch.ID,
					"expired_date": batch.ExpiredDate.Format("2006-01-02"),
					"qty_deducted": qtyPotong,
					"sisa_batch":   batch.TotalQty,
				})
			}

			transaksi := models.TransaksiStok{
				BarangID:         req.BarangID,
				KontakID:         req.KontakID, // Zero NULL: Customer ID recorded
				TanggalTransaksi: now,
				JenisTransaksi:   "OUT",
				Qty:              req.Qty,
				Keterangan:       req.Keterangan,
				CreatedAt:        now,
			}
			if err := tx.Create(&transaksi).Error; err != nil {
				return fmt.Errorf("gagal menyimpan transaksi log: %w", err)
			}

			resultTransaksi = transaksi
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "FAILED",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":           "SUCCESS",
		"message":          fmt.Sprintf("Mutasi stok %s berhasil diproses", req.JenisTransaksi),
		"transaksi":        resultTransaksi,
		"affected_batches": affectedBatches,
	})
}

package routes

import (
	"time"

	"kasirApi/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter registers all API endpoints
func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "AgriStock Inventory API (Party Model)",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	barangHandler := handlers.NewBarangHandler()
	kontakHandler := handlers.NewKontakHandler()
	transaksiHandler := handlers.NewTransaksiHandler()
	reportHandler := handlers.NewReportHandler()

	api := r.Group("/api")
	{
		// 1. Master Barang
		api.GET("/barang", barangHandler.GetAll)
		api.GET("/barang/:id", barangHandler.GetByID)
		api.POST("/barang", barangHandler.Create)
		api.PUT("/barang/:id", barangHandler.Update)
		api.DELETE("/barang/:id", barangHandler.Delete)

		// 2. Master Kontaks (Party Model: Supplier & Customer)
		api.GET("/kontaks", kontakHandler.GetAll)
		api.POST("/kontaks", kontakHandler.Create)
		api.PUT("/kontaks/:id", kontakHandler.Update)
		api.DELETE("/kontaks/:id", kontakHandler.Delete)

		// 3. Mutasi Stok (IN / OUT FEFO)
		api.POST("/transaksi-stok", transaksiHandler.HandleMutasi)

		// 4. Analytics & Reports
		api.GET("/top-barang", reportHandler.GetTopBarang)
		api.GET("/tracking-stok", reportHandler.GetTrackingStok)
		api.GET("/stok-per-barang", reportHandler.GetStokPerBarang)
		api.GET("/alerts/low-stock", reportHandler.GetLowStockAlerts)
	}

	return r
}

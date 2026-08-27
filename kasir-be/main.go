package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"kasirApi/config"
	"kasirApi/routes"

	"github.com/spf13/viper"
)

func main() {
	// Setup Viper configuration
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	port := viper.GetString("PORT")
	if port == "" {
		port = "8080"
	}

	dbDriver := viper.GetString("DB_DRIVER")
	if dbDriver == "" {
		dbDriver = "mysql"
	}

	dbConn := viper.GetString("DB_CONN")
	if dbConn == "" {
		dbConn = "root:@tcp(127.0.0.1:3306)/toko_pertanian?parseTime=true&multiStatements=true"
	}

	// Initialize Database (GORM + Connection Pool + AutoMigrate + Seed)
	log.Printf("[SERVER] Initializing database (driver: %s)...\n", dbDriver)
	_, err := config.InitDatabase(dbConn, dbDriver)
	if err != nil {
		log.Fatalf("[SERVER] Database initialization failed: %v", err)
	}

	// Setup Gin Engine and Routes
	r := routes.SetupRouter()

	addr := ":" + port
	fmt.Println("=======================================================")
	fmt.Printf("🌱 AgriStock Backend API running on http://localhost:%s\n", port)
	fmt.Println("=======================================================")

	if err := r.Run(addr); err != nil {
		log.Fatalf("[SERVER] Failed to run server: %v", err)
	}
}

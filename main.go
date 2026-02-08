package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"kasirApi/database"
	"kasirApi/handlers"
	"kasirApi/models"
	"kasirApi/repositories"
	"kasirApi/services"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	// "data"
)

func main() {
	// GET localhost:8080/api/Category/{id}
	// PUT localhost:8080/api/Category/{id}
	// DELETE localhost:8080/api/Category/{id}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// Setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialze database", err)
	}
	defer db.Close()

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	// Transaction
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	http.HandleFunc("/api/checkout", transactionHandler.HandleCheckout) // POST
	http.HandleFunc("/api/report/hari-ini", transactionHandler.GetReport)

	// Setup routes
	http.HandleFunc("/api/category", categoryHandler.HandleCategory)
	http.HandleFunc("/api/category/", categoryHandler.HandleCategoryByID)

	http.HandleFunc("/api/Category/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategoryByID(w, r)
		} else if r.Method == "PUT" {
			updateCategory(w, r)
		} else if r.Method == "DELETE" {
			deleteCategory(w, r)
		}

	})

	// GET localhost:8080/api/Category
	// POST localhost:8080/api/Category
	// http.HandleFunc("/api/Category", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method == "GET" {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		json.NewEncoder(w).Encode(Category)
	// 	} else if r.Method == "POST" {
	// 		// baca data from reqquest
	// 		var CategoryBaru Category
	// 		err := json.NewDecoder(r.Body).Decode(&CategoryBaru)
	// 		if err != nil {
	// 			http.Error(w, "Invalid request", http.StatusBadRequest)
	// 			return
	// 		}

	// 		// masukin data to variable Category
	// 		CategoryBaru.ID = len(Category) + 1
	// 		Category = append(Category, CategoryBaru)
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusCreated) //Success or 201
	// 		json.NewEncoder(w).Encode(CategoryBaru)
	// 	}

	// })

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})
	fmt.Println("Server running di localhost:8080")

	// err = http.ListenAndServe(":8080", nil)
	// if err != nil {
	// 	fmt.Println("gagal running server")
	// }

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running di", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("gagal running server", err)
	}
}

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

var Category = []models.Category{
	{ID: 1, Name: "Makanan", Price: 5000, Stock: 15},
	{ID: 2, Name: "Minuman", Price: 7000, Stock: 25},
	{ID: 3, Name: "Camilan", Price: 3000, Stock: 30},
}

func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/Category/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	for _, p := range Category {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "Category belum ada", http.StatusNotFound)

}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	// get id dari request
	idStr :=
		strings.TrimPrefix(r.URL.Path, "/api/Category/")

	// ganti int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	//get data dari request
	var updateCategory models.Category
	err = json.NewDecoder(r.Body).Decode(&updateCategory)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return

	}

	// loop Category, find id, change data type as requested
	for i := range Category {
		if Category[i].ID == id {
			updateCategory.ID = id
			Category[i] = updateCategory

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateCategory)
			return
		}
	}
	http.Error(w, "Category belum ada", http.StatusNotFound)
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	// get id
	idStr := strings.TrimPrefix(r.URL.Path, "/api/Category/")
	//  ganti id to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}
	// loop Category cari ID, dapat index yg mau dihapus
	for i, p := range Category {
		if p.ID == id {
			// bikin slicec baru dengan data sebelum dan sesudah indexing
			Category = append(Category[:i], Category[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "sukses delete",
			})

			return
		}
	}
	http.Error(w, "Category belum ada", http.StatusNotFound)
}

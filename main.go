package main

import (
	"kasirApi/handlers"
	"kasirApi/models"
	"kasirApi/repositories"
	"kasirApi/services"
	"kasirApi/database"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"strconv"
	"strings"
	"log"
	// "data"
)

func main() {
	// GET localhost:8080/api/categories/{id}
	// PUT localhost:8080/api/categories/{id}
	// DELETE localhost:8080/api/categories/{id}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port: viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running di", addr)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("gagal running server", err)
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
	
	// Setup routes
	http.HandleFunc("/api/category", categoryHandler.HandleCategory)
	http.HandleFunc("/api/category/", categoryHandler.HandleCategoryByID)

	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategoriesByID(w, r)
		} else if r.Method == "PUT" {
			updateCategories(w, r)
		} else if r.Method == "DELETE" {
			deleteCategories(w, r)
		}

	})

	

	// GET localhost:8080/api/categories
	// POST localhost:8080/api/categories
	// http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method == "GET" {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		json.NewEncoder(w).Encode(categories)
	// 	} else if r.Method == "POST" {
	// 		// baca data from reqquest
	// 		var categoriesBaru Category
	// 		err := json.NewDecoder(r.Body).Decode(&categoriesBaru)
	// 		if err != nil {
	// 			http.Error(w, "Invalid request", http.StatusBadRequest)
	// 			return
	// 		}

	// 		// masukin data to variable categories
	// 		categoriesBaru.ID = len(categories) + 1
	// 		categories = append(categories, categoriesBaru)
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusCreated) //Success or 201
	// 		json.NewEncoder(w).Encode(categoriesBaru)
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

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

var categories = []models.Category{
	{ID: 1, Name: "Makanan", Price: 5000, Stock: 15},
	{ID: 2, Name: "Minuman", Price: 7000, Stock: 25},
	{ID: 3, Name: "Camilan", Price: 3000, Stock: 30},
}

func getCategoriesByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Categories ID", http.StatusBadRequest)
		return
	}

	for _, p := range categories {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "Categories belum ada", http.StatusNotFound)

}

func updateCategories(w http.ResponseWriter, r *http.Request) {
	// get id dari request
	idStr :=
		strings.TrimPrefix(r.URL.Path, "/api/categories/")

	// ganti int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Categories ID", http.StatusBadRequest)
		return
	}

	//get data dari request
	var updateCategories models.Category
	err = json.NewDecoder(r.Body).Decode(&updateCategories)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return

	}

	// loop categories, find id, change data type as requested
	for i := range categories {
		if categories[i].ID == id {
			updateCategories.ID = id
			categories[i] = updateCategories

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateCategories)
			return
		}
	}
	http.Error(w, "Categories belum ada", http.StatusNotFound)
}

func deleteCategories(w http.ResponseWriter, r *http.Request) {
	// get id
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	//  ganti id to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Categories ID", http.StatusBadRequest)
		return
	}
	// loop categories cari ID, dapat index yg mau dihapus
	for i, p := range categories {
		if p.ID == id {
			// bikin slicec baru dengan data sebelum dan sesudah indexing
			categories = append(categories[:i], categories[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "sukses delete",
			})

			return
		}
	}
	http.Error(w, "Categories belum ada", http.StatusNotFound)
}

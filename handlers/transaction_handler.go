package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"kasirApi/models"
	"kasirApi/services"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// multiple item dengan quantity
func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.HandleCheckout(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TransactionHandler) Checkout (w http.ResponseWriter, r *http.Request) {
	var req models.CheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	transaction, err := h.service.Checkout(req.Items, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

func (h *TransactionHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	// Ambil parameter dari URL: ?start_date=...&end_date=...
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	// Logic Default: Jika user tidak kirim tanggal, anggap "HARI INI"
	if startDate == "" || endDate == "" {
		currentTime := time.Now()
		// Format tanggal ke string YYYY-MM-DD
		// 00:00:00 hari ini sampai 23:59:59 hari ini
		startDate = currentTime.Format("2006-01-02") + " 00:00:00"
		endDate = currentTime.Format("2006-01-02") + " 23:59:59"
	} else {
		// Jika user kirim tanggal, tambahkan jam biar akurat
		startDate = startDate + " 00:00:00"
		endDate = endDate + " 23:59:59"
	}

	// Panggil Service/Repo
	// (Anggap kamu langsung panggil repo di sini, idealnya lewat Service dulu)
	report, err := h.service.GetSalesReport(startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

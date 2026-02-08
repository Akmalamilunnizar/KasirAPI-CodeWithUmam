package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"kasirApi/models"
	// "kasirApi/repositories"
)

type transactionRepository struct {
	db *sql.DB
}

type TransactionRepository interface {
	// Update(transaction *models.Transaction) error
	Delete(id int) error
	// Checkout(items []models.CheckoutItem, useLock bool) (*models.Transaction, error)
	CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error)
	GetSalesReport(startDate, endDate string) (*models.SalesReport, error)
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (repo *transactionRepository) GetAll(nameFilter string) ([]models.Transaction, error) {
	query := "SELECT id, name, price, stock FROM transaction"

	args := []interface{}{}
	if nameFilter != "" {
		query += " WHERE name ILIKE $1"
		args = append(args, "%"+nameFilter+"%")
	}

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	transaction := make([]models.Transaction, 0)
	for rows.Next() {
		var p models.Transaction
		// err := rows.Scan(&p.ID, &p.ProductName, &p.Price, &p.Stock)
		if err != nil {
			return nil, err
		}
		transaction = append(transaction, p)
	}

	return transaction, nil
}

func (repo *transactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	if len(details) > 0 {
		query := "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES "
		var args []interface{}

		// loop
		for i, d := range details {
			d.TransactionID = transactionID
			n := i * 4
			// Rumus posisi parameter:
			// Baris 1: $1, $2, $3, $4
			// Baris 2: $5, $6, $7, $8

			// placeholder ke string query
			query += fmt.Sprintf("($%d, $%d, $%d, $%d)", n+1, n+2, n+3, n+4)

			// masukan ke slice artgs
			args = append(args, d.TransactionID, d.ProductID, d.Quantity, d.Subtotal)

		}

		query = query[:len(query)-1]

		_, err = tx.Exec(query, args...)
		if err != nil {
			return nil, err
		}
	}
	// old ways
	// for i := range details {
	// 	details[i].TransactionID = transactionID
	// 	_, err = tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
	// 		transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Update TransactionID di object return agar frontend dapat ID yang benar OPTIONAL NEXT
	// for i := range details {
	//     details[i].TransactionID = transactionID
	// }

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

func (r *transactionRepository) GetSalesReport(startDate, endDate string) (*models.SalesReport, error) {
	var report models.SalesReport

	// 1. Query Total Revenue & Total Transaksi
	// Gunakan COALESCE agar jika tidak ada data, hasilnya 0 (bukan NULL error)
	queryStats := `
		SELECT 
			COALESCE(SUM(total_amount), 0), 
			COUNT(id)
		FROM transactions 
		WHERE created_at >= $1 AND created_at <= $2`

	err := r.db.QueryRow(queryStats, startDate, endDate).Scan(&report.TotalRevenue, &report.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	// 2. Query Produk Terlaris (Top 1)
	// Kita join transaction_details ke products, lalu filter berdasarkan tanggal transaksi
	queryTopProduct := `
		SELECT 
			p.name, 
			COALESCE(SUM(td.quantity), 0) as total_qty
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE t.created_at >= $1 AND t.created_at <= $2
		GROUP BY p.name
		ORDER BY total_qty DESC
		LIMIT 1`

	err = r.db.QueryRow(queryTopProduct, startDate, endDate).Scan(&report.ProdukTerlaris.Nama, &report.ProdukTerlaris.QtyTerjual)

	// Handle kasus jika tidak ada penjualan sama sekali (sql.ErrNoRows)
	if err != nil {
		// Jika errornya bukan NoRows, return error beneran
		if err.Error() != "sql: no rows in result set" {
			// Kalau belum ada penjualan, set default strip/0 biar gak error
			report.ProdukTerlaris = models.ProductBestSeller{Nama: "-", QtyTerjual: 0}
		}
	}

	return &report, nil
}

func (repo *transactionRepository) Delete(id int) error {
	query := "DELETE FROM transaction WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("transaction tidak ditemukan")
	}

	return err
}

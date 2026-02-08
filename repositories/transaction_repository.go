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

	for i := range details {
		details[i].TransactionID = transactionID
		_, err = tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit();  err != nil{
		return nil, err
	}

	return &models.Transaction{
		ID: transactionID,
		TotalAmount: totalAmount,
		Details: details,
	}, nil
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

package services

import (
	"kasirApi/models"
	"kasirApi/repositories"
)

type TransactionService struct {
	repo repositories.TransactionRepository
}

func NewTransactionService(repo repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

// func (s *TransactionService) GetAll(name string) ([]models.Transaction, error) {
// 	return s.repo.GetAll(name)
// }

// func (s *TransactionService) Create(data *models.Transaction) error {
// 	return s.repo.Create(data)
// }

// func (s *TransactionService) GetByID(id int) (*models.Transaction, error) {
// 	return s.repo.GetByID(id)
// }

// func (s *TransactionService) Update(transaction *models.Transaction) error {
// 	return s.repo.Update(transaction)
// }

func (s *TransactionService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *TransactionService) Checkout(items []models.CheckoutItem, useLock bool) (*models.Transaction, error)  {
	return s.repo.CreateTransaction(items)
}

func (s *TransactionService) GetSalesReport(startDate, endDate string) (*models.SalesReport, error) {
    // Service acts as a bridge here
    return s.repo.GetSalesReport(startDate, endDate)
}
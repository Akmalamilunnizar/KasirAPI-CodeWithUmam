package services

import (
	"kasirApi/models"
	"kasirApi/repositories"
)

type ProdukService struct {
	repo repositories.ProdukRepository
}

func NewProdukService(repo repositories.ProdukRepository) *ProdukService {
	return &ProdukService{repo: repo}
}

func (s *ProdukService) GetAll(name string) ([]models.Produk, error) {
	return s.repo.GetAll(name)
}

func (s *ProdukService) Create(data *models.Produk) error {
	return s.repo.Create(data)
}

func (s *ProdukService) GetByID(id int) (*models.Produk, error) {
	return s.repo.GetByID(id)
}

func (s *ProdukService) Update(produk *models.Produk) error {
	return s.repo.Update(produk)
}

func (s *ProdukService) Delete(id int) error {
	return s.repo.Delete(id)
}

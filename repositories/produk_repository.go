package repositories

import (
	"database/sql"
	"errors"
	"kasirApi/models"

)

type produkRepository struct {
	db *sql.DB
}

type ProdukRepository interface {
	GetAll(name string) ([]models.Produk, error)
	Create(produk *models.Produk) error
	GetByID(id int) (*models.Produk, error)
	Update(produk *models.Produk) error
	Delete(id int) error
}

func NewProdukRepository(db *sql.DB) ProdukRepository {
	return &produkRepository{db: db}
}

func (repo *produkRepository) GetAll(nameFilter string) ([]models.Produk, error) {
	query := "SELECT id, name, price, stock FROM produk"
	
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

	produk := make([]models.Produk, 0)
	for rows.Next() {
		var p models.Produk
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
		if err != nil {
			return nil, err
		}
		produk = append(produk, p)
	}

	return produk, nil
}


func (repo *produkRepository) Create(produk *models.Produk) error {
	query := "INSERT INTO produk (name, price, stock) VALUES ($1, $2, $3) RETURNING id"
	err := repo.db.QueryRow(query, produk.Name, produk.Price, produk.Stock).Scan(&produk.ID)
	return err
}

// GetByID - ambil produk by ID
func (repo *produkRepository) GetByID(id int) (*models.Produk, error) {
	query := "SELECT id, name, price, stock FROM produk WHERE id = $1"

	var p models.Produk
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	if err == sql.ErrNoRows {
		return nil, errors.New("produk tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (repo *produkRepository) Update(produk *models.Produk) error {
	query := "UPDATE produk SET name = $1, price = $2, stock = $3 WHERE id = $4"
	result, err := repo.db.Exec(query, produk.Name, produk.Price, produk.Stock, produk.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return nil
}

func (repo *produkRepository) Delete(id int) error {
	query := "DELETE FROM produk WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return err
}

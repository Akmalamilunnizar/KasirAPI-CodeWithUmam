package models

import (
	"time"

	"gorm.io/gorm"
)

// Constants
const (
	TipeSupplier = "SUPPLIER"
	TipeCustomer = "CUSTOMER"

	EntitasBadanUsaha = "BADAN_USAHA"
	EntitasIndividu   = "INDIVIDU"

	KatBenihBibit    = "BENIH_BIBIT"
	KatPupukPadat    = "PUPUK_PADAT"
	KatPupukCair     = "PUPUK_CAIR"
	KatPestisida     = "PESTISIDA"
	KatAlatPertanian = "ALAT_PERTANIAN"
)

// Kontak represents Party model (Supplier or Customer, PT/CV or Individu)
type Kontak struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Tipe         string         `gorm:"size:20;not null;index" json:"tipe"` // SUPPLIER / CUSTOMER
	JenisEntitas string         `gorm:"size:20;not null;default:'INDIVIDU'" json:"jenis_entitas"` // BADAN_USAHA / INDIVIDU
	KodeKontak   string         `gorm:"size:50;not null;uniqueIndex" json:"kode_kontak"`
	NamaKontak   string         `gorm:"size:255;not null" json:"nama_kontak"`
	NoTelp       string         `gorm:"size:30" json:"no_telp"`
	Alamat       string         `gorm:"type:text" json:"alamat"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// Barang represents inventory item
type Barang struct {
	ID         uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	KodeBarang string         `gorm:"size:50;not null;uniqueIndex" json:"kode_barang"`
	NamaBarang string         `gorm:"size:255;not null" json:"nama_barang"`
	Kategori   string         `gorm:"size:50;not null;index" json:"kategori"`
	MinStok    int            `gorm:"not null;default:0" json:"min_stok"`
	PcsPerBox  int            `gorm:"not null;default:12" json:"pcs_per_box"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Virtual fields
	TotalStok  int    `gorm:"-" json:"total_stok"`
	TotalBox   int    `gorm:"-" json:"total_box"`
	SisaPcs    int    `gorm:"-" json:"sisa_pcs"`
	StokStatus string `gorm:"-" json:"stok_status"`
}

// StokBatch represents FEFO batch
type StokBatch struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	BarangID    uint      `gorm:"not null;index:idx_stok_batches_fefo,priority:1" json:"barang_id"`
	KontakID    *uint     `gorm:"index" json:"kontak_id"` // Supplier asal
	ExpiredDate time.Time `gorm:"type:date;not null;index:idx_stok_batches_fefo,priority:2" json:"expired_date"`
	TotalQty    int       `gorm:"not null;default:0;index:idx_stok_batches_fefo,priority:3" json:"total_qty"`
	CreatedAt   time.Time `json:"created_at"`

	Barang *Barang `gorm:"foreignKey:BarangID;constraint:OnDelete:RESTRICT" json:"barang,omitempty"`
	Kontak *Kontak `gorm:"foreignKey:KontakID;constraint:OnDelete:SET NULL" json:"kontak,omitempty"`
}

// TransaksiStok represents single-ledger stock mutation (Zero NULL on KontakID)
type TransaksiStok struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	BarangID         uint      `gorm:"not null;index" json:"barang_id"`
	KontakID         uint      `gorm:"not null;index" json:"kontak_id"` // Supplier (IN) or Customer (OUT)
	TanggalTransaksi time.Time `gorm:"not null;index" json:"tanggal_transaksi"`
	JenisTransaksi   string    `gorm:"size:3;not null" json:"jenis_transaksi"` // IN / OUT
	Qty              int       `gorm:"not null" json:"qty"`
	Keterangan       string    `gorm:"type:text" json:"keterangan"`
	CreatedAt        time.Time `json:"created_at"`

	Barang *Barang `gorm:"foreignKey:BarangID;constraint:OnDelete:RESTRICT" json:"barang,omitempty"`
	Kontak *Kontak `gorm:"foreignKey:KontakID;constraint:OnDelete:RESTRICT" json:"kontak,omitempty"`
}

// Request & Response DTOs
type MutasiRequest struct {
	BarangID       uint   `json:"barang_id" binding:"required"`
	KontakID       uint   `json:"kontak_id" binding:"required"`
	JenisTransaksi string `json:"jenis_transaksi" binding:"required,oneof=IN OUT"`
	Qty            int    `json:"qty" binding:"required,gt=0"`
	ExpiredDate    string `json:"expired_date"` // YYYY-MM-DD (Wajib jika IN)
	Keterangan     string `json:"keterangan"`
}

type TopBarangResponse struct {
	BarangID     uint   `json:"barang_id"`
	KodeBarang   string `json:"kode_barang"`
	NamaBarang   string `json:"nama_barang"`
	Kategori     string `json:"kategori"`
	TotalTerjual int    `json:"total_terjual"`
}

type StokPerBarangResponse struct {
	BatchID      uint      `json:"batch_id"`
	BarangID     uint      `json:"barang_id"`
	KodeBarang   string    `json:"kode_barang"`
	NamaBarang   string    `json:"nama_barang"`
	Kategori     string    `json:"kategori"`
	PcsPerBox    int       `json:"pcs_per_box"`
	NamaSupplier string    `json:"nama_supplier"`
	JenisEntitas string    `json:"jenis_entitas"`
	ExpiredDate  time.Time `json:"expired_date"`
	TotalQtyPcs  int       `json:"total_qty_pcs"`
	Box          int       `json:"box"`
	Pcs          int       `json:"pcs"`
	IsNearED     bool      `json:"is_near_ed"`
	DaysUntilED  int       `json:"days_until_ed"`
}

type LowStockAlertResponse struct {
	BarangID    uint   `json:"barang_id"`
	KodeBarang  string `json:"kode_barang"`
	NamaBarang  string `json:"nama_barang"`
	Kategori    string `json:"kategori"`
	MinStok     int    `json:"min_stok"`
	CurrentStok int    `json:"current_stok"`
	Status      string `json:"status"`
	Deficit     int    `json:"deficit"`
}

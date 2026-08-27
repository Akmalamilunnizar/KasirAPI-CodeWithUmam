package config

import (
	"log"
	"strings"
	"time"

	"kasirApi/models"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase initializes GORM with Connection Pool and Seed Data
func InitDatabase(dsn string, driver string) (*gorm.DB, error) {
	var dialector gorm.Dialector

	if driver == "sqlite" || strings.HasSuffix(dsn, ".db") {
		dialector = sqlite.Open(dsn)
	} else {
		dialector = mysql.Open(dsn)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	// Auto Migration
	err = db.AutoMigrate(
		&models.Kontak{},
		&models.Barang{},
		&models.StokBatch{},
		&models.TransaksiStok{},
	)
	if err != nil {
		log.Printf("[DATABASE] AutoMigrate warning: %v\n", err)
	}

	// Seed Initial Data
	SeedInitialData(db)

	DB = db
	return db, nil
}

// SeedInitialData populates demo agricultural contacts, items, batches, and transactions
func SeedInitialData(db *gorm.DB) {
	var count int64
	db.Model(&models.Kontak{}).Count(&count)
	if count > 0 {
		return
	}

	log.Println("[DATABASE] Seeding initial Party Model (Kontaks) demo data...")

	kontaks := []models.Kontak{
		// Suppliers (Badan Usaha)
		{ID: 1, Tipe: models.TipeSupplier, JenisEntitas: models.EntitasBadanUsaha, KodeKontak: "SPL-SYNG-01", NamaKontak: "PT Syngenta Indonesia", NoTelp: "021-78831234", Alamat: "Cilandak Commercial Estate, Jakarta Selatan"},
		{ID: 2, Tipe: models.TipeSupplier, JenisEntitas: models.EntitasBadanUsaha, KodeKontak: "SPL-PI-02", NamaKontak: "PT Pupuk Indonesia (Persero)", NoTelp: "021-5481200", Alamat: "Jl. Kemanggisan Utama No. 1, Jakarta Barat"},
		{ID: 3, Tipe: models.TipeSupplier, JenisEntitas: models.EntitasBadanUsaha, KodeKontak: "SPL-BISI-03", NamaKontak: "PT BISI International Tbk", NoTelp: "0354-682123", Alamat: "Jl. Raya Surabaya-Mojokerto KM 19, Kediri"},
		{ID: 4, Tipe: models.TipeSupplier, JenisEntitas: models.EntitasBadanUsaha, KodeKontak: "SPL-CBA-04", NamaKontak: "PT CBA Chemical Industry", NoTelp: "021-5961234", Alamat: "Kawasan Industri Modern Cikande, Serang"},

		// Customers (Badan Usaha & Individu)
		{ID: 5, Tipe: models.TipeCustomer, JenisEntitas: models.EntitasBadanUsaha, KodeKontak: "CST-POK-01", NamaKontak: "Gapoktan Makmur Sejahtera (Kelompok Tani)", NoTelp: "0812-3456-7890", Alamat: "Desa Sidomulyo RT 02/03, Kediri"},
		{ID: 6, Tipe: models.TipeCustomer, JenisEntitas: models.EntitasIndividu, KodeKontak: "CST-IND-02", NamaKontak: "Pak Slamet Riyadi (Petani Jagung)", NoTelp: "0857-1122-3344", Alamat: "Dusun Krajan, Pare"},
		{ID: 7, Tipe: models.TipeCustomer, JenisEntitas: models.EntitasIndividu, KodeKontak: "CST-IND-03", NamaKontak: "Haji Sukirman (Petani Bawang Merah)", NoTelp: "0813-9988-7766", Alamat: "Kecamatan Plemahan"},
		{ID: 8, Tipe: models.TipeCustomer, JenisEntitas: models.EntitasBadanUsaha, KodeKontak: "CST-CV-04", NamaKontak: "CV Tani Berkah Sub-Agen", NoTelp: "0354-998811", Alamat: "Jl. Pahlawan No. 45, Nganjuk"},
	}
	for _, k := range kontaks {
		db.FirstOrCreate(&k, models.Kontak{KodeKontak: k.KodeKontak})
	}

	barangs := []models.Barang{
		{ID: 1, KodeBarang: "BNH-JG-001", NamaBarang: "Benih Jagung Hibrida Pioneer P27", Kategori: models.KatBenihBibit, MinStok: 30, PcsPerBox: 20},
		{ID: 2, KodeBarang: "PPK-NPK-001", NamaBarang: "Pupuk NPK Phonska Plus 1kg", Kategori: models.KatPupukPadat, MinStok: 50, PcsPerBox: 25},
		{ID: 3, KodeBarang: "PST-CRC-500", NamaBarang: "Insektisida Curacron 500EC 100ml", Kategori: models.KatPestisida, MinStok: 25, PcsPerBox: 10},
		{ID: 4, KodeBarang: "PST-RND-486", NamaBarang: "Herbisida Roundup 486SL 200ml", Kategori: models.KatPestisida, MinStok: 40, PcsPerBox: 20},
		{ID: 5, KodeBarang: "ALT-SPR-016", NamaBarang: "Sprayer Elektrik CBA 16 Liter", Kategori: models.KatAlatPertanian, MinStok: 5, PcsPerBox: 1},
		{ID: 6, KodeBarang: "PPK-GND-100", NamaBarang: "Pupuk Daun Gandasil D 100gr", Kategori: models.KatPupukPadat, MinStok: 15, PcsPerBox: 50},
	}
	for _, b := range barangs {
		db.FirstOrCreate(&b, models.Barang{KodeBarang: b.KodeBarang})
	}

	now := time.Now()
	supp1 := uint(1)
	supp2 := uint(2)
	supp3 := uint(3)
	supp4 := uint(4)

	batches := []models.StokBatch{
		{ID: 1, BarangID: 1, KontakID: &supp3, ExpiredDate: now.AddDate(0, 2, 0), TotalQty: 45},
		{ID: 2, BarangID: 1, KontakID: &supp3, ExpiredDate: now.AddDate(0, 10, 0), TotalQty: 100},
		{ID: 3, BarangID: 2, KontakID: &supp2, ExpiredDate: now.AddDate(0, 24, 0), TotalQty: 120},
		{ID: 4, BarangID: 3, KontakID: &supp1, ExpiredDate: now.AddDate(0, 1, 15), TotalQty: 12},
		{ID: 5, BarangID: 3, KontakID: &supp1, ExpiredDate: now.AddDate(0, 12, 0), TotalQty: 30},
		{ID: 6, BarangID: 4, KontakID: &supp1, ExpiredDate: now.AddDate(0, 15, 0), TotalQty: 35},
		{ID: 7, BarangID: 5, KontakID: &supp4, ExpiredDate: now.AddDate(5, 0, 0), TotalQty: 8},
		{ID: 8, BarangID: 6, KontakID: &supp4, ExpiredDate: now.AddDate(0, 8, 0), TotalQty: 10},
	}
	for _, b := range batches {
		db.FirstOrCreate(&b, models.StokBatch{ID: b.ID})
	}

	transaksis := []models.TransaksiStok{
		{BarangID: 1, KontakID: 3, TanggalTransaksi: now.AddDate(0, 0, -5), JenisTransaksi: "IN", Qty: 150, Keterangan: "Penerimaan PO Distributor BISI"},
		{BarangID: 2, KontakID: 2, TanggalTransaksi: now.AddDate(0, 0, -4), JenisTransaksi: "IN", Qty: 120, Keterangan: "Penerimaan PO Pupuk Indonesia"},
		{BarangID: 3, KontakID: 1, TanggalTransaksi: now.AddDate(0, 0, -4), JenisTransaksi: "IN", Qty: 50, Keterangan: "Penerimaan PO PT Syngenta"},
		{BarangID: 1, KontakID: 5, TanggalTransaksi: now.AddDate(0, 0, -3), JenisTransaksi: "OUT", Qty: 35, Keterangan: "Penjualan Nota #TRX-001 Poktan Makmur"},
		{BarangID: 1, KontakID: 6, TanggalTransaksi: now.AddDate(0, 0, -2), JenisTransaksi: "OUT", Qty: 20, Keterangan: "Penjualan Nota #TRX-002 Pak Slamet"},
		{BarangID: 4, KontakID: 8, TanggalTransaksi: now.AddDate(0, 0, -1), JenisTransaksi: "OUT", Qty: 25, Keterangan: "Penjualan Nota #TRX-003 CV Tani Berkah"},
		{BarangID: 3, KontakID: 7, TanggalTransaksi: now.AddDate(0, 0, -1), JenisTransaksi: "OUT", Qty: 8, Keterangan: "Penjualan Nota #TRX-004 Haji Sukirman"},
	}
	for _, t := range transaksis {
		db.Create(&t)
	}

	log.Println("[DATABASE] Seed initial data completed")
}

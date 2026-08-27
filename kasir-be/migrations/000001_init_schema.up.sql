-- 1. Table: suppliers
CREATE TABLE IF NOT EXISTS suppliers (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    kode_supplier VARCHAR(50) NOT NULL UNIQUE,
    nama_supplier VARCHAR(255) NOT NULL,
    no_telp VARCHAR(30) NULL,
    alamat TEXT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    INDEX idx_suppliers_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 2. Table: barangs
CREATE TABLE IF NOT EXISTS barangs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    kode_barang VARCHAR(50) NOT NULL UNIQUE,
    nama_barang VARCHAR(255) NOT NULL,
    kategori VARCHAR(50) NOT NULL,
    min_stok INT NOT NULL DEFAULT 0,
    pcs_per_box INT NOT NULL DEFAULT 12,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    INDEX idx_barangs_kategori (kategori),
    INDEX idx_barangs_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3. Table: stok_batches (FEFO Support + Supplier Link)
CREATE TABLE IF NOT EXISTS stok_batches (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    barang_id BIGINT NOT NULL,
    supplier_id BIGINT NULL,
    expired_date DATE NOT NULL,
    total_qty INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_stok_batches_fefo (barang_id, expired_date ASC, total_qty),
    INDEX idx_stok_batches_supplier (supplier_id),
    CONSTRAINT fk_stok_batches_barang 
        FOREIGN KEY (barang_id) 
        REFERENCES barangs(id) 
        ON DELETE RESTRICT,
    CONSTRAINT fk_stok_batches_supplier 
        FOREIGN KEY (supplier_id) 
        REFERENCES suppliers(id) 
        ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 4. Table: transaksi_stoks (Audit Log Mutasi + Supplier Link)
CREATE TABLE IF NOT EXISTS transaksi_stoks (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    barang_id BIGINT NOT NULL,
    supplier_id BIGINT NULL,
    tanggal_transaksi DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    jenis_transaksi ENUM('IN', 'OUT') NOT NULL,
    qty INT NOT NULL,
    keterangan TEXT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_transaksi_stoks_barang_tgl (barang_id, tanggal_transaksi DESC),
    INDEX idx_transaksi_stoks_supplier (supplier_id),
    CONSTRAINT fk_transaksi_stoks_barang 
        FOREIGN KEY (barang_id) 
        REFERENCES barangs(id) 
        ON DELETE RESTRICT,
    CONSTRAINT fk_transaksi_stoks_supplier 
        FOREIGN KEY (supplier_id) 
        REFERENCES suppliers(id) 
        ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

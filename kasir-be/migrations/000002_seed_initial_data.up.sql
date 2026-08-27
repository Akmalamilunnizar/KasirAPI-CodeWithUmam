-- 1. Seed Suppliers
INSERT INTO suppliers (id, kode_supplier, nama_supplier, no_telp, alamat, created_at, updated_at)
VALUES 
    (1, 'SPL-SYNG-01', 'PT Syngenta Indonesia', '021-78831234', 'Cilandak Commercial Estate, Jakarta Selatan', NOW(), NOW()),
    (2, 'SPL-PI-02', 'PT Pupuk Indonesia (Persero)', '021-5481200', 'Jl. Kemanggisan Utama No. 1, Jakarta Barat', NOW(), NOW()),
    (3, 'SPL-BISI-03', 'PT BISI International Tbk (Pioneer)', '0354-682123', 'Jl. Raya Surabaya-Mojokerto KM 19, Kediri', NOW(), NOW()),
    (4, 'SPL-CBA-04', 'PT CBA Chemical Industry', '021-5961234', 'Kawasan Industri Modern Cikande, Serang', NOW(), NOW())
ON DUPLICATE KEY UPDATE nama_supplier=VALUES(nama_supplier);

-- 2. Seed Barangs
INSERT INTO barangs (id, kode_barang, nama_barang, kategori, min_stok, pcs_per_box, created_at, updated_at)
VALUES 
    (1, 'BNH-JG-001', 'Benih Jagung Hibrida Pioneer P27', 'BENIH_BIBIT', 20, 20, NOW(), NOW()),
    (2, 'PPK-NPK-001', 'Pupuk NPK Phonska Plus 1kg', 'PUPUK_PADAT', 50, 25, NOW(), NOW()),
    (3, 'PST-CRC-500', 'Insektisida Curacron 500EC 100ml', 'PESTISIDA', 15, 10, NOW(), NOW()),
    (4, 'PST-RND-486', 'Herbisida Roundup 486SL 200ml', 'PESTISIDA', 25, 20, NOW(), NOW()),
    (5, 'ALT-SPR-016', 'Sprayer Elektrik CBA 16 Liter', 'ALAT_PERTANIAN', 5, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE nama_barang=VALUES(nama_barang);

-- 3. Seed Stok Batches (FEFO + Supplier)
INSERT INTO stok_batches (id, barang_id, supplier_id, expired_date, total_qty, created_at)
VALUES 
    -- Jagung dari PT BISI (Pioneer)
    (1, 1, 3, DATE_ADD(CURRENT_DATE, INTERVAL 6 MONTH), 100, NOW()),
    (2, 1, 3, DATE_ADD(CURRENT_DATE, INTERVAL 12 MONTH), 150, NOW()),
    -- Pupuk NPK dari Pupuk Indonesia
    (3, 2, 2, DATE_ADD(CURRENT_DATE, INTERVAL 24 MONTH), 200, NOW()),
    -- Curacron dari Syngenta
    (4, 3, 1, DATE_ADD(CURRENT_DATE, INTERVAL 3 MONTH), 30, NOW()),
    (5, 3, 1, DATE_ADD(CURRENT_DATE, INTERVAL 18 MONTH), 50, NOW()),
    -- Roundup
    (6, 4, 1, DATE_ADD(CURRENT_DATE, INTERVAL 15 MONTH), 80, NOW()),
    -- Sprayer dari CBA
    (7, 5, 4, DATE_ADD(CURRENT_DATE, INTERVAL 5 YEAR), 10, NOW())
ON DUPLICATE KEY UPDATE total_qty=VALUES(total_qty);

-- 4. Seed Transaksi Stoks (Audit In/Out)
INSERT INTO transaksi_stoks (barang_id, supplier_id, tanggal_transaksi, jenis_transaksi, qty, keterangan, created_at)
VALUES 
    (1, 3, DATE_SUB(NOW(), INTERVAL 2 DAY), 'IN', 250, 'Penerimaan PO Supplier BISI International', NOW()),
    (2, 2, DATE_SUB(NOW(), INTERVAL 2 DAY), 'IN', 200, 'Penerimaan PO Distributor Pupuk Indonesia', NOW()),
    (3, 1, DATE_SUB(NOW(), INTERVAL 1 DAY), 'IN', 80, 'Penerimaan PO PT Syngenta Indonesia', NOW()),
    (1, NULL, DATE_SUB(NOW(), INTERVAL 1 DAY), 'OUT', 5, 'Penjualan Nota #TRX-001 Petani Pak Joko', NOW()),
    (3, NULL, DATE_SUB(NOW(), INTERVAL 5 HOUR), 'OUT', 2, 'Penjualan Nota #TRX-002', NOW());

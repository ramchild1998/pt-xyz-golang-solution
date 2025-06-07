CREATE DATABASE IF NOT EXISTS pt_xyz_db;

USE pt_xyz_db;

-- Tabel untuk menyimpan data personal konsumen
CREATE TABLE IF NOT EXISTS konsumen (
    nik VARCHAR(16) PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    legal_name VARCHAR(255) NOT NULL,
    tempat_lahir VARCHAR(100) NOT NULL,
    tanggal_lahir DATE NOT NULL,
    gaji DECIMAL(15, 2) NOT NULL,
    foto_ktp VARCHAR(255),
    foto_selfie VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Tabel untuk menyimpan limit tenor konsumen
CREATE TABLE IF NOT EXISTS limit_tenor (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    konsumen_nik VARCHAR(16) NOT NULL,
    tenor_bulan INT NOT NULL,
    sisa_limit DECIMAL(15, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (konsumen_nik) REFERENCES konsumen(nik) ON DELETE CASCADE,
    UNIQUE KEY `idx_nik_tenor` (`konsumen_nik`, `tenor_bulan`)
);

-- Tabel untuk mencatat transaksi
CREATE TABLE IF NOT EXISTS transaksi (
    nomor_kontrak VARCHAR(255) PRIMARY KEY,
    konsumen_nik VARCHAR(16) NOT NULL,
    otr DECIMAL(15, 2) NOT NULL,
    admin_fee DECIMAL(15, 2) NOT NULL,
    jumlah_cicilan INT NOT NULL,
    jumlah_bunga DECIMAL(15, 2) NOT NULL DEFAULT 0,
    nama_asset VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (konsumen_nik) REFERENCES konsumen(nik) ON DELETE CASCADE
);

-- Seed data contoh
INSERT INTO konsumen (nik, full_name, legal_name, tempat_lahir, tanggal_lahir, gaji) VALUES
('1111222233334444', 'Budi Santoso', 'Budi Santoso', 'Jakarta', '1990-01-15', 10000000),
('5555666677778888', 'Annisa Fitriani', 'Annisa Fitriani', 'Bandung', '1992-05-20', 15000000);

INSERT INTO limit_tenor (konsumen_nik, tenor_bulan, sisa_limit) VALUES
('1111222233334444', 1, 100000),
('1111222233334444', 2, 200000),
('1111222233334444', 3, 500000),
('1111222233334444', 6, 700000),
('5555666677778888', 1, 1000000),
('5555666677778888', 2, 1200000),
('5555666677778888', 3, 1500000),
('5555666677778888', 6, 2000000);
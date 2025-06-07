package entity

import "time"

// Transaksi merepresentasikan data satu transaksi
type Transaksi struct {
	NomorKontrak  string    `json:"nomor_kontrak" db:"nomor_kontrak"`
	KonsumenNIK   string    `json:"konsumen_nik" db:"konsumen_nik"`
	OTR           float64   `json:"otr" db:"otr"`
	AdminFee      float64   `json:"admin_fee" db:"admin_fee"`
	JumlahCicilan int       `json:"jumlah_cicilan" db:"jumlah_cicilan"`
	JumlahBunga   float64   `json:"jumlah_bunga" db:"jumlah_bunga"`
	NamaAsset     string    `json:"nama_asset" db:"nama_asset"`
	Status        string    `json:"status" db:"status"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// CreateTransaksiRequest adalah DTO untuk request pembuatan transaksi
type CreateTransaksiRequest struct {
	KonsumenNIK string  `json:"nik" binding:"required"`
	TenorBulan  int     `json:"tenor" binding:"required"`
	OTR         float64 `json:"otr" binding:"required,gt=0"`
	AdminFee    float64 `json:"admin_fee" binding:"required,gte=0"`
	NamaAsset   string  `json:"nama_asset" binding:"required"`
}

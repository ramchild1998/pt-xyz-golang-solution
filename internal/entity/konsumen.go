package entity

import "time"

// Konsumen merepresentasikan data personal konsumen
type Konsumen struct {
	NIK          string    `json:"nik" db:"nik"`
	FullName     string    `json:"full_name" db:"full_name"`
	LegalName    string    `json:"legal_name" db:"legal_name"`
	TempatLahir  string    `json:"tempat_lahir" db:"tempat_lahir"`
	TanggalLahir time.Time `json:"tanggal_lahir" db:"tanggal_lahir"`
	Gaji         float64   `json:"gaji" db:"gaji"`
	FotoKTP      string    `json:"foto_ktp" db:"foto_ktp"`
	FotoSelfie   string    `json:"foto_selfie" db:"foto_selfie"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// LimitTenor merepresentasikan limit kredit konsumen per tenor
type LimitTenor struct {
	ID          int64     `json:"id" db:"id"`
	KonsumenNIK string    `json:"konsumen_nik" db:"konsumen_nik"`
	TenorBulan  int       `json:"tenor_bulan" db:"tenor_bulan"`
	SisaLimit   float64   `json:"sisa_limit" db:"sisa_limit"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

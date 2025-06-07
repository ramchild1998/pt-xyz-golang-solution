package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ramchild1998/pt-xyz-case-study/internal/entity"
	"github.com/ramchild1998/pt-xyz-case-study/internal/repository"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrLimitNotFound     = errors.New("limit untuk tenor yang dipilih tidak ditemukan")
	ErrInsufficientLimit = errors.New("limit kredit tidak mencukupi untuk transaksi ini")
	ErrInternalServer    = errors.New("terjadi kesalahan internal pada server")
)

type TransaksiUsecase interface {
	CreateTransaksi(ctx context.Context, req *entity.CreateTransaksiRequest) (*entity.Transaksi, error)
}

type transaksiUsecase struct {
	transaksiRepo repository.TransaksiRepository
	db            *sqlx.DB
}

func NewTransaksiUsecase(repo repository.TransaksiRepository, db *sqlx.DB) TransaksiUsecase {
	return &transaksiUsecase{
		transaksiRepo: repo,
		db:            db,
	}
}

func (u *transaksiUsecase) CreateTransaksi(ctx context.Context, req *entity.CreateTransaksiRequest) (*entity.Transaksi, error) {
	// Memulai transaksi database (DB Transaction) untuk menjaga ACID
	tx, err := u.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, ErrInternalServer
	}
	// Defer rollback, jika terjadi error, transaksi akan dibatalkan
	defer tx.Rollback()

	// 1. Ambil limit konsumen dengan locking (FOR UPDATE)
	limit, err := u.transaksiRepo.GetLimitByNIKAndTenor(ctx, req.KonsumenNIK, req.TenorBulan)
	if err != nil {
		return nil, ErrInternalServer
	}
	if limit == nil {
		return nil, ErrLimitNotFound
	}

	// 2. Validasi limit
	totalBiayaTransaksi := req.OTR // Asumsi OTR adalah total yang mengurangi limit
	if limit.SisaLimit < totalBiayaTransaksi {
		return nil, ErrInsufficientLimit
	}

	// 3. Siapkan data transaksi baru
	newTransaksi := &entity.Transaksi{
		NomorKontrak:  fmt.Sprintf("KTR-%s", uuid.New().String()),
		KonsumenNIK:   req.KonsumenNIK,
		OTR:           req.OTR,
		AdminFee:      req.AdminFee,
		JumlahCicilan: req.TenorBulan, // Jumlah cicilan sama dengan tenor
		JumlahBunga:   0,              // Perhitungan bunga bisa ditambahkan di sini
		NamaAsset:     req.NamaAsset,
		Status:        "PENDING",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// 4. Update sisa limit
	limit.SisaLimit -= totalBiayaTransaksi
	limit.UpdatedAt = time.Now()

	// 5. Simpan transaksi dan update limit dalam satu DB transaction
	err = u.transaksiRepo.CreateTransaksiAndUpdateLimit(ctx, tx, newTransaksi, limit)
	if err != nil {
		return nil, ErrInternalServer
	}

	// Jika semua berhasil, commit transaksi
	if err := tx.Commit(); err != nil {
		return nil, ErrInternalServer
	}

	return newTransaksi, nil
}

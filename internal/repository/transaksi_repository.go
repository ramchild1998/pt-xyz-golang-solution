package repository

import (
	"context"
	"database/sql"

	"github.com/ramchild1998/pt-xyz-case-study/internal/entity"

	"github.com/jmoiron/sqlx"
)

type TransaksiRepository interface {
	GetLimitByNIKAndTenor(ctx context.Context, nik string, tenor int) (*entity.LimitTenor, error)
	CreateTransaksiAndUpdateLimit(ctx context.Context, tx *sqlx.Tx, transaksi *entity.Transaksi, limit *entity.LimitTenor) error
}

type transaksiRepository struct {
	db *sqlx.DB
}

func NewTransaksiRepository(db *sqlx.DB) TransaksiRepository {
	return &transaksiRepository{db: db}
}

// GetLimitByNIKAndTenor mengambil data limit berdasarkan NIK dan tenor
func (r *transaksiRepository) GetLimitByNIKAndTenor(ctx context.Context, nik string, tenor int) (*entity.LimitTenor, error) {
	var limit entity.LimitTenor
	query := "SELECT id, konsumen_nik, tenor_bulan, sisa_limit FROM limit_tenor WHERE konsumen_nik = ? AND tenor_bulan = ? FOR UPDATE"
	err := r.db.GetContext(ctx, &limit, query, nik, tenor)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Data tidak ditemukan dianggap bukan error fatal
		}
		return nil, err
	}
	return &limit, nil
}

// CreateTransaksiAndUpdateLimit menyimpan transaksi baru dan memperbarui sisa limit dalam satu transaksi DB.
// **PENANGANAN TRANSAKSI KONKUREN (ACID & ISOLATION)**
// Penggunaan `tx *sqlx.Tx` memastikan bahwa semua operasi di dalamnya (insert transaksi dan update limit)
// bersifat atomik. Jika salah satu gagal, seluruhnya akan di-rollback.
// Klausa `FOR UPDATE` pada query `GetLimitByNIKAndTenor` mengunci baris data limit selama transaksi,
// mencegah 'race condition' di mana dua transaksi untuk konsumen yang sama diproses secara bersamaan.
func (r *transaksiRepository) CreateTransaksiAndUpdateLimit(ctx context.Context, tx *sqlx.Tx, transaksi *entity.Transaksi, limit *entity.LimitTenor) error {
	// 1. Insert ke tabel transaksi
	queryTransaksi := `INSERT INTO transaksi (nomor_kontrak, konsumen_nik, otr, admin_fee, jumlah_cicilan, nama_asset, status) 
                       VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := tx.ExecContext(ctx, queryTransaksi, transaksi.NomorKontrak, transaksi.KonsumenNIK, transaksi.OTR, transaksi.AdminFee, transaksi.JumlahCicilan, transaksi.NamaAsset, transaksi.Status)
	if err != nil {
		return err
	}

	// 2. Update sisa limit di tabel limit_tenor
	queryLimit := `UPDATE limit_tenor SET sisa_limit = ? WHERE id = ?`
	_, err = tx.ExecContext(ctx, queryLimit, limit.SisaLimit, limit.ID)
	if err != nil {
		return err
	}

	return nil
}

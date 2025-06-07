package repository

import (
	"context"

	"github.com/ramchild1998/pt-xyz-case-study/internal/entity"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
)

type MockTransaksiRepository struct {
	mock.Mock
}

func (m *MockTransaksiRepository) GetLimitByNIKAndTenor(ctx context.Context, nik string, tenor int) (*entity.LimitTenor, error) {
	args := m.Called(ctx, nik, tenor)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.LimitTenor), args.Error(1)
}

func (m *MockTransaksiRepository) CreateTransaksiAndUpdateLimit(ctx context.Context, tx *sqlx.Tx, transaksi *entity.Transaksi, limit *entity.LimitTenor) error {
	args := m.Called(ctx, tx, transaksi, limit)
	return args.Error(0)
}

package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/ramchild1998/pt-xyz-case-study/internal/entity"
	"github.com/ramchild1998/pt-xyz-case-study/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransaksiUsecase_CreateTransaksi(t *testing.T) {
	// Setup mock DB
	mockDB, mockSQL, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Setup mock repository
	mockRepo := new(repository.MockTransaksiRepository)

	// Inisialisasi usecase dengan mock
	uc := NewTransaksiUsecase(mockRepo, sqlxDB)

	ctx := context.Background()

	t.Run("Success Case", func(t *testing.T) {
		req := &entity.CreateTransaksiRequest{
			KonsumenNIK: "12345",
			TenorBulan:  3,
			OTR:         400000,
		}

		mockLimit := &entity.LimitTenor{
			ID:          1,
			KonsumenNIK: "12345",
			TenorBulan:  3,
			SisaLimit:   500000,
		}

		mockRepo.On("GetLimitByNIKAndTenor", ctx, req.KonsumenNIK, req.TenorBulan).Return(mockLimit, nil).Once()
		mockRepo.On("CreateTransaksiAndUpdateLimit", ctx, mock.Anything, mock.AnythingOfType("*entity.Transaksi"), mock.AnythingOfType("*entity.LimitTenor")).Return(nil).Once()

		mockSQL.ExpectBegin()
		mockSQL.ExpectCommit()

		transaksi, err := uc.CreateTransaksi(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, transaksi)
		assert.Equal(t, req.KonsumenNIK, transaksi.KonsumenNIK)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Insufficient Limit", func(t *testing.T) {
		req := &entity.CreateTransaksiRequest{
			KonsumenNIK: "12345",
			TenorBulan:  3,
			OTR:         600000, // Melebihi limit
		}

		mockLimit := &entity.LimitTenor{
			ID:          1,
			KonsumenNIK: "12345",
			TenorBulan:  3,
			SisaLimit:   500000,
		}

		mockRepo.On("GetLimitByNIKAndTenor", ctx, req.KonsumenNIK, req.TenorBulan).Return(mockLimit, nil).Once()

		mockSQL.ExpectBegin()
		mockSQL.ExpectRollback()

		transaksi, err := uc.CreateTransaksi(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, transaksi)
		assert.True(t, errors.Is(err, ErrInsufficientLimit))
		mockRepo.AssertExpectations(t)
	})
}

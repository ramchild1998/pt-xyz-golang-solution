package handler

import (
	"errors"
	"net/http"

	"github.com/ramchild1998/pt-xyz-case-study/internal/entity"
	"github.com/ramchild1998/pt-xyz-case-study/internal/usecase"

	"github.com/gin-gonic/gin"
)

type TransaksiHandler struct {
	transaksiUsecase usecase.TransaksiUsecase
}

func NewTransaksiHandler(r *gin.Engine, uc usecase.TransaksiUsecase) {
	handler := &TransaksiHandler{transaksiUsecase: uc}
	r.POST("/transaksi", handler.CreateTransaksi)
}

func (h *TransaksiHandler) CreateTransaksi(c *gin.Context) {
	var req entity.CreateTransaksiRequest

	// **PENCEGAHAN KEAMANAN OWASP (Input Validation)**
	// Gin `ShouldBindJSON` secara default memvalidasi tipe data dan
	// aturan `binding` yang didefinisikan di struct. Ini mencegah serangan
	// seperti Injeksi (misal, tipe data tidak valid) pada level awal.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request: " + err.Error()})
		return
	}

	transaksi, err := h.transaksiUsecase.CreateTransaksi(c.Request.Context(), &req)
	if err != nil {
		// Mapping error dari usecase ke HTTP status code
		switch {
		case errors.Is(err, usecase.ErrLimitNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrInsufficientLimit):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
		return
	}

	c.JSON(http.StatusCreated, transaksi)
}

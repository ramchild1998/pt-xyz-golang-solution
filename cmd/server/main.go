package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ramchild1998/pt-xyz-case-study/internal/handler"
	"github.com/ramchild1998/pt-xyz-case-study/internal/repository"
	"github.com/ramchild1998/pt-xyz-case-study/internal/usecase"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {
	// **PENCEGAHAN KEAMANAN OWASP A05:2021 - Security Misconfiguration**
	// Menggunakan environment variables untuk konfigurasi sensitif (koneksi DB),
	// bukan hardcoding.
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// **PENCEGAHAN KEAMANAN OWASP A03:2021 - Injection**
	// Penggunaan prepared statements ditangani oleh driver `database/sql` dan `sqlx`
	// saat parameter diteruskan dengan `?`, ini adalah praktik standar yang aman.
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(127.0.0.1:3306)/pt_xyz_db?parseTime=true"
		log.Println("DATABASE_DSN not set, using default")
	}

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// **PENCEGAHAN KEAMANAN OWASP - Security Misconfiguration**
	// Nonaktifkan mode debug Gin di lingkungan produksi.
	gin.SetMode(os.Getenv("GIN_MODE")) // "release" for production

	router := gin.Default()

	// Dependency Injection
	transaksiRepo := repository.NewTransaksiRepository(db)
	transaksiUsecase := usecase.NewTransaksiUsecase(transaksiRepo, db)
	handler.NewTransaksiHandler(router, transaksiUsecase)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s\n", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

# Author : Ramadiansyah
# Position : Backend Developer


# PT XYZ Multifinance - Studi Kasus Solusi Microservice

Proyek ini adalah implementasi solusi backend untuk studi kasus PT XYZ Multifinance, dibangun menggunakan Go (Golang) dengan arsitektur *microservices*, *Clean Architecture*, dan dikemas menggunakan Docker.

## Arsitektur Aplikasi

Solusi ini mengadopsi arsitektur *microservices* untuk meningkatkan skalabilitas, *maintainability*, dan kecepatan *deployment*.

![Arsitektur Aplikasi](https://drive.google.com/file/d/14GFCw4GV8sVe_ajDHS5RVVZi_2NoW2cJ/view?usp=sharing)

## Desain Database (ERD)

Database dirancang untuk mendukung layanan yang ada.

**Database Transaksi:**
![ERD Database Transaksi](https://drive.google.com/file/d/1sES2UR882h4OAZyqP8vBKfPPlzaMFt1b/view?usp=sharing)

## Persyaratan Minimum yang Dipenuhi

-   **Git Flow**: Pengembangan menggunakan alur kerja Git Flow (`master`, `develop`, `feature/*`).
-   **Clean Architecture**: Struktur proyek memisahkan `entity`, `repository`, `usecase`, dan `handler`.
-   **Handling Concurrent Transaction**: Penanganan transaksi bersamaan diimplementasikan pada endpoint pembuatan transaksi (`POST /transaksi`) menggunakan *database transaction locking* (`SELECT ... FOR UPDATE`) untuk memastikan integritas data (ACID).
-   **Pencegahan Keamanan (OWASP TOP 10)**:
    1.  **A01:2021 - Broken Access Control**: Diterapkan melalui validasi input yang ketat.
    2.  **A03:2021 - Injection**: Dicegah dengan penggunaan *prepared statements* secara default oleh `database/sql` dan `sqlx`.
    3.  **A05:2021 - Security Misconfiguration**: Konfigurasi sensitif (koneksi DB) dikelola melalui *environment variables* (.env) dan tidak di-*hardcode*. Mode debug dinonaktifkan di produksi.
-   **Unit Test**: Unit test disediakan untuk lapisan *usecase* dengan *mocking* pada lapisan *repository*.

## Nilai Tambah

-   **Dockerize**: Seluruh aplikasi dan database-nya telah di-dockerisasi dan dapat dijalankan dengan satu perintah menggunakan `docker-compose`.

---

## Instalasi dan Menjalankan Proyek

### Prasyarat

-   [Docker](https://www.docker.com/products/docker-desktop/)
-   [Docker Compose](https://docs.docker.com/compose/install/)
-   [Git](https://git-scm.com/downloads)

### Langkah-langkah

1.  **Clone Repository**

    ```bash
    git clone https://github.com/ramchild1998/pt-xyz-golang-solution.git
    cd pt-xyz-case-study
    ```

2.  **Konfigurasi Environment**

    Proyek ini menggunakan file `.env` untuk konfigurasi. Salin file contoh dan sesuaikan jika perlu (konfigurasi default sudah diatur untuk `docker-compose`).

    ```bash
    # Tidak perlu disalin jika file .env sudah ada di repo
    cp .env.example .env
    ```

3.  **Jalankan dengan Docker Compose**

    Perintah ini akan membangun image aplikasi, membuat kontainer untuk aplikasi dan database MySQL, serta menjalankan skema database dari `schema.sql`.

    ```bash
    docker-compose up --build
    ```

    Aplikasi akan berjalan di `http://localhost:8080`.

## Menguji Aplikasi

Gunakan `curl` atau Postman untuk menguji endpoint.

### Endpoint: `POST /transaksi`

Endpoint ini untuk membuat transaksi baru.

**Contoh Request (Sukses):**

Melakukan transaksi untuk Budi (NIK `1111222233334444`) dengan tenor 3 bulan (limit 500.000) sebesar 450.000.

```bash
curl -X POST http://localhost:8080/transaksi \
-H "Content-Type: application/json" \
-d '{
    "nik": "1111222233334444",
    "tenor": 3,
    "otr": 450000,
    "admin_fee": 5000,
    "nama_asset": "Smartphone Canggih"
}'
```

**Contoh Request (Gagal - Limit Tidak Cukup):**

Melakukan transaksi lagi yang melebihi sisa limit.

```bash
curl -X POST http://localhost:8080/transaksi \
-H "Content-Type: application/json" \
-d '{
    "nik": "1111222233334444",
    "tenor": 3,
    "otr": 100000,
    "nama_asset": "Headphone Keren"
}'
```

Response yang diharapkan:
```json
{
    "error": "limit kredit tidak mencukupi untuk transaksi ini"
}
```

## Menjalankan Unit Test

Untuk menjalankan unit test, Anda perlu Go terinstal di mesin lokal Anda.

```bash
go test ./... -v
```
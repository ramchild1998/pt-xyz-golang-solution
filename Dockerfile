# --- Build Stage ---
    FROM golang:1.24-alpine AS builder

    # Set working directory
    WORKDIR /app
    
    # Copy go.mod and go.sum files
    COPY go.mod go.sum ./
    
    # Download dependencies
    RUN go mod download
    
    # Copy the source code
    COPY . .
    
    # Build the application
    # -ldflags="-w -s" -> untuk mengurangi ukuran binary
    # CGO_ENABLED=0 -> untuk static build
    RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /main ./cmd/server
    
    # --- Final Stage ---
    FROM alpine:latest
    
    # Set working directory
    WORKDIR /root/
    
    # Copy the pre-built binary from the builder stage
    COPY --from=builder /main .
    
    # Copy .env file
    # COPY .env . (Jika Anda ingin memasukkan .env ke dalam image, tidak direkomendasikan untuk produksi)
    
    # Expose port
    EXPOSE 8080
    
    # Command to run the executable
    CMD ["./main"]
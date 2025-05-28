# Sử dụng image chính thức của Golang
FROM golang:1.21 as builder

WORKDIR /app

# Copy go.mod và go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build ứng dụng
RUN go build -o main .

# Stage chạy app
FROM debian:bullseye-slim

# Cài đặt CA-certificates (bắt buộc cho HTTPS)
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Tạo thư mục ứng dụng
WORKDIR /app

# Copy binary từ stage builder
COPY --from=builder /app/main .

# Copy file .env (tùy chọn, có thể dùng Docker secrets hoặc ENV flags sau)
COPY .env .env

# Expose port của ứng dụng
EXPOSE 8003

# Chạy ứng dụng
CMD ["./main"]

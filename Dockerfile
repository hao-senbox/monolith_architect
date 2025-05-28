# Sử dụng image chính thức của Golang
FROM golang:1.24 as builder

WORKDIR /app

# Copy go.mod và go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build ứng dụng từ đúng path
RUN go build -o main ./cmd/server

# Stage chạy app
FROM debian:bullseye-slim

# Cài đặt CA-certificates (bắt buộc cho HTTPS)
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary từ stage builder
COPY --from=builder /app/main .

# Expose port của ứng dụng
EXPOSE 8003

# Chạy ứng dụng
CMD ["./main"]
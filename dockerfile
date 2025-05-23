FROM golang:1.24.1-alpine

# Cài git nếu dùng go get hoặc module cần clone
RUN apk add --no-cache git

# Tạo thư mục làm việc
WORKDIR /app

# Copy toàn bộ mã nguồn
COPY . .


# (Tùy chọn) tải module trước để tăng tốc build lần sau
RUN go mod download

# Chạy ứng dụng bằng `go run`
CMD ["go", "run", "cmd/main.go"]

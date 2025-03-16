FROM golang:1.21-alpine AS builder

WORKDIR /app

# Bağımlılıkları kopyala ve indir
COPY go.mod go.sum ./
RUN go mod download

# Kaynak kodları kopyala
COPY . .

# Uygulamayı derle
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./server/cmd/main.go

# Final imajı
FROM alpine:3.18

WORKDIR /app

# CA sertifikalarını ekle
RUN apk --no-cache add ca-certificates

# Builder aşamasından derlenmiş uygulamayı kopyala
COPY --from=builder /app/main .
COPY --from=builder /app/config /app/config
COPY --from=builder /app/migration /app/migration

# Uygulamayı çalıştır
CMD ["/app/main"] 
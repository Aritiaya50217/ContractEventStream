# ---------- Build stage ----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# install git (needed for go mod)
RUN apk add --no-cache git

# COPY เฉพาะไฟล์ go.mod และ go.sum จาก service
COPY go.mod go.sum ./

# ดาวน์โหลด dependencies
RUN go mod download

# COPY source code ของ service
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o workflow-service cmd/main.go

# ---------- Runtime stage ----------
FROM alpine:3.20

WORKDIR /app

# certs for postgres / kafka tls
RUN apk add --no-cache ca-certificates

# copy binary จาก build stage
COPY --from=builder /app/workflow-service .

EXPOSE 8080

CMD ["./workflow-service"]

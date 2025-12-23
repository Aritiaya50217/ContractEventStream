# ---------- Build stage ----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# install git (needed for go mod sometimes)
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o workflow-service cmd/main.go


# ---------- Runtime stage ----------
FROM alpine:3.20

WORKDIR /app

# certs for postgres / kafka tls (future proof)
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/workflow-service .

EXPOSE 8080

CMD ["./workflow-service"]

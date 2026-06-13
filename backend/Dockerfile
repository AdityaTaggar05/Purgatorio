# ---------- Build stage ----------
FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# IMPORTANT: correct build path
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -buildid=" -trimpath \
    -o purg-backend ./cmd/server

# ---------- Runtime stage ----------
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /app/purg-backend .
COPY ./keys ./keys

EXPOSE 8080

CMD ["./purg-backend"]


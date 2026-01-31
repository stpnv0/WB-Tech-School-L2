FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o calendar ./cmd/app

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/calendar /app/calendar

COPY config ./config

ENV CONFIG_PATH=/app/config/config.yaml

EXPOSE 8080

CMD ["/app/calendar"]

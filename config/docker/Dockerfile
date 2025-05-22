FROM golang:1.24-alpine AS builder

WORKDIR /opt/app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o crypto-knight-history ./cmd/cmd.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /opt/app

COPY --from=builder /opt/app/crypto-knight-history .

EXPOSE 8080

CMD ["./crypto-knight-history"]

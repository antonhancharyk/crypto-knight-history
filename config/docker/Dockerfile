FROM golang:alpine AS builder

WORKDIR /opt/app

COPY . .

RUN go build -o crypto-knight-history ./cmd/cmd.go

FROM alpine

WORKDIR /opt/app

COPY --from=builder /opt/app/crypto-knight-history /opt/app/crypto-knight-history

EXPOSE 8080

CMD ["./crypto-knight-history"]
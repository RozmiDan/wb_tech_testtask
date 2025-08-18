FROM golang:1.25 AS builder

WORKDIR /wb_service/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o app ./cmd/main.go

FROM debian:latest

COPY --from=builder /wb_service/app/app .

EXPOSE 8080

CMD ["/app"]

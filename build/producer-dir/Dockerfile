FROM golang:1.25 AS builder
ENV GOPROXY=https://goproxy.io,direct \
    GOSUMDB=sum.golang.org \
    GO111MODULE=on
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /producer-dir ./cmd/producer-dir

FROM gcr.io/distroless/base-debian12
COPY --from=builder /producer-dir /producer-dir
ENTRYPOINT ["/producer-dir"]

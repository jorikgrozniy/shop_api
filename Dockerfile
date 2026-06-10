FROM golang:1.26.1 AS builder
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/api

FROM scratch
COPY --from=builder /build/app .
CMD ["/app"]
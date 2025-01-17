FROM golang:1.23 AS builder
WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags="-w -s" -o ms cmd/microservice/main.go

FROM scratch
COPY --from=builder /app/ms /app/ms
CMD ["/app/ms"]
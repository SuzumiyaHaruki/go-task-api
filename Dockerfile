FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY . .
RUN go build -o /app/bin/server ./cmd/server

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/bin/server /app/server
EXPOSE 8080
CMD ["/app/server"]

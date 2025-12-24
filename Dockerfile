FROM golang:1.23.6-alpine AS builder
WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o /app/main ./cmd/urlshortener

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main /app/main
EXPOSE 8080
CMD ["/app/main"]


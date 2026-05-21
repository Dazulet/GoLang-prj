FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o manga-service ./cmd/server

FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
RUN mkdir -p uploads/covers uploads/avatars uploads/chapters
COPY --from=builder /app/manga-service .
EXPOSE 8082
CMD ["./manga-service"]

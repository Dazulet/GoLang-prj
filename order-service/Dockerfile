FROM golang:1.24-alpine

ENV GOTOOLCHAIN=local

WORKDIR /app

COPY . .

RUN go mod edit -go=1.24 && \
    go mod edit -require=github.com/gin-gonic/gin@v1.10.0 && \
    go mod edit -replace=golang.org/x/net=golang.org/x/net@v0.33.0 && \
    go mod tidy

RUN go build -o main .

CMD ["./main"]
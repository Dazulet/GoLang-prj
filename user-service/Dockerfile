FROM golang:1.24-alpine

# Запрещаем искать версию 1.25 в сети
ENV GOTOOLCHAIN=local

WORKDIR /app

# ШАГ 1: Сначала копируем ВЕСЬ код (включая плохой go.mod)
COPY . .

# ШАГ 2: Теперь ПЕРЕЗАПИСЫВАЕМ настройки прямо внутри Docker
# Мы заставляем его использовать 1.24 и старые версии библиотек
RUN go mod edit -go=1.24 && \
    go mod edit -require=github.com/gin-gonic/gin@v1.10.0 && \
    go mod edit -replace=golang.org/x/net=golang.org/x/net@v0.33.0 && \
    go mod tidy

# ШАГ 3: Теперь собираем. Теперь go.mod внутри Docker исправлен
RUN go build -o main .

CMD ["./main"]
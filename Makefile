.PHONY: up up-detached down logs health test-all tidy-all migrate-up migrate-down

## ── Docker ──────────────────────────────────────────────────────────────────

up:
	docker-compose up --build

up-detached:
	docker-compose up --build -d

down:
	docker-compose down

down-volumes:
	docker-compose down -v

logs:
	docker-compose logs -f

logs-auth:
	docker-compose logs -f auth-service

logs-manga:
	docker-compose logs -f manga-service

logs-comment:
	docker-compose logs -f comment-service

## ── Health ──────────────────────────────────────────────────────────────────

health:
	@echo "→ Auth Service"
	@curl -sf http://localhost:8081/health | python3 -m json.tool || echo "UNREACHABLE"
	@echo "→ Manga Service"
	@curl -sf http://localhost:8082/health | python3 -m json.tool || echo "UNREACHABLE"
	@echo "→ Comment Service"
	@curl -sf http://localhost:8083/health | python3 -m json.tool || echo "UNREACHABLE"

## ── Tests ───────────────────────────────────────────────────────────────────

test-auth:
	cd auth-service && go test ./tests/... -v

test-manga:
	cd manga-service && go test ./tests/... -v

test-comment:
	cd comment-service && go test ./tests/... -v

test-all: test-auth test-manga test-comment

cover-auth:
	cd auth-service && go test ./tests/... -cover

cover-all:
	cd auth-service    && go test ./tests/... -cover
	cd manga-service   && go test ./tests/... -cover
	cd comment-service && go test ./tests/... -cover

## ── Go tooling ──────────────────────────────────────────────────────────────

tidy-all:
	cd auth-service    && go mod tidy
	cd manga-service   && go mod tidy
	cd comment-service && go mod tidy

## ── Migrations (requires migrate CLI) ───────────────────────────────────────
##    Install: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-up:
	migrate -path migrations/auth    -database "postgres://auth_user:auth_secret@localhost:5433/auth_db?sslmode=disable" up
	migrate -path migrations/manga   -database "postgres://manga_user:manga_secret@localhost:5434/manga_db?sslmode=disable" up
	migrate -path migrations/comment -database "postgres://comment_user:comment_secret@localhost:5435/comment_db?sslmode=disable" up

migrate-down:
	migrate -path migrations/auth    -database "postgres://auth_user:auth_secret@localhost:5433/auth_db?sslmode=disable" down 1
	migrate -path migrations/manga   -database "postgres://manga_user:manga_secret@localhost:5434/manga_db?sslmode=disable" down 1
	migrate -path migrations/comment -database "postgres://comment_user:comment_secret@localhost:5435/comment_db?sslmode=disable" down 1

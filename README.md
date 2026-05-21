# MangaLib — Go Microservices Backend

> University Showcase Project · Go · Gin · PostgreSQL · GORM · JWT · Resty v2 · Docker

A production-inspired backend for a manga reading platform, built as three lightweight,
independently deployable microservices. Every component communicates over clean REST APIs,
is containerised with Docker, and is covered by unit tests.

---

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Tech Stack](#tech-stack)
3. [Microservices](#microservices)
4. [Inter-Service Communication (Resty v2)](#inter-service-communication-resty-v2)
5. [JWT Authentication](#jwt-authentication)
6. [Project Structure](#project-structure)
7. [Database Migrations](#database-migrations)
8. [Quick Start — Docker](#quick-start--docker)
9. [Running Locally (without Docker)](#running-locally-without-docker)
10. [API Overview](#api-overview)
11. [Unit Tests](#unit-tests)
12. [Screenshots](#screenshots)

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                  Browser / Frontend (port 3000)             │
└───────────┬─────────────────┬──────────────────┬────────────┘
            │                 │                  │
            ▼                 ▼                  ▼
  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────────┐
  │  Auth Service   │  │  Manga Service  │  │ Comment Service  │
  │   :8081         │  │   :8082         │  │   :8083          │
  │                 │  │                 │  │                  │
  │ • Register      │  │ • Manga CRUD    │  │ • Comments       │
  │ • Login         │  │ • Chapters      │  │ • Likes          │
  │ • JWT issuance  │  │ • Pages         │  │ • Threads        │
  │ • User profiles │  │ • Genres/Tags   │  │                  │
  │ • /auth/validate│  │ • Bookmarks     │  │ ←── Resty v2 ───►│
  └────────┬────────┘  │ • Ratings       │  │   calls Auth     │
           │           │ • Progress      │  └──────────────────┘
           │           │                 │
           │           │ ←── Resty v2 ──►│
           │           │   calls Auth    │
           │           └─────────────────┘
           │
   ┌───────┴──────┐  ┌──────────────┐  ┌──────────────┐
   │   auth_db    │  │   manga_db   │  │  comment_db  │
   │  PostgreSQL  │  │  PostgreSQL  │  │  PostgreSQL  │
   └──────────────┘  └──────────────┘  └──────────────┘
```

**Key design choices:**
- **One database per service** — true data isolation, no shared tables
- **JWT parsed locally** — no Auth roundtrip on every protected request
- **Resty v2** — used for all service-to-service HTTP calls
- **Best-effort enrichment** — Comment Service enriches authors from Auth; if Auth is down, comments still return

---

## Tech Stack

| Layer              | Technology                        |
|--------------------|-----------------------------------|
| Language           | Go 1.22                           |
| HTTP Framework     | Gin                               |
| ORM                | GORM                              |
| Database           | PostgreSQL 16                     |
| Authentication     | JWT (HS256) via golang-jwt/jwt/v5 |
| Inter-service HTTP | **Resty v2** (go-resty/resty/v2)  |
| Migrations         | golang-migrate/migrate/v4         |
| Testing            | testing + httptest + testify      |
| Containerisation   | Docker + Docker Compose           |
| Password Hashing   | bcrypt (cost 12)                  |

---

## Microservices

### 🔐 Auth Service — port 8081

Owns everything user-identity related. It is the **only** service that issues JWT tokens.

**Responsibilities:** registration, login, JWT token signing, user profiles, token validation endpoint

**Database:** `auth_db` — `users` table

### 📚 Manga Service — port 8082

The largest service. Owns all reading-platform content.

**Responsibilities:** manga CRUD, chapters, pages, genres, tags, bookmarks, ratings, reading progress

**Database:** `manga_db` — `manga`, `chapters`, `pages`, `genres`, `tags`, `bookmarks`, `ratings`, `reading_progress`

**Resty v2 call:** `GET /api/users/:id` on Auth Service → proxy user profiles via `/api/users/:id/profile`

### 💬 Comment Service — port 8083

Small, focused service. Only stores comment data; user info is fetched on demand.

**Responsibilities:** comments, threaded replies, likes

**Database:** `comment_db` — `comments`, `comment_likes`

**Resty v2 calls:** `GET /api/users/:id` (enrich author info), `POST /api/auth/validate` (optional server-side token verification)

---

## Inter-Service Communication (Resty v2)

The university requirement for **Resty v2** is implemented in two places:

### Comment Service → Auth Service

```
internal/client/auth_client.go
```

```go
// Resty v2 client — configured with timeout and retry
r := resty.New().
    SetBaseURL(baseURL).
    SetTimeout(5 * time.Second).
    SetRetryCount(1).
    SetRetryWaitTime(500 * time.Millisecond)

// Actual inter-service HTTP call
resp, err := c.client.R().
    SetResult(&result).
    Get(fmt.Sprintf("/api/users/%d", userID))
```

When a user posts a comment, the Comment Service calls the Auth Service to attach
author name and avatar to the response. The call is **best-effort** — if Auth is
temporarily down, comments are returned without author details rather than failing.

### Manga Service → Auth Service

```
internal/client/auth_client.go
```

The Manga Service exposes `GET /api/users/:id/profile` which internally uses
Resty v2 to proxy the request to the Auth Service. This demonstrates
bidirectional, real inter-service HTTP communication.

---

## JWT Authentication

All three services share the same `JWT_SECRET` environment variable.

**Flow:**
```
Client                Auth Service           Manga / Comment Service
  │                       │                         │
  ├─ POST /auth/login ───►│                         │
  │                       ├─ validate credentials   │
  │◄─ { token: "eyJ..." } ┤                         │
  │                       │                         │
  ├─ GET /api/manga ──────┼─────── Bearer eyJ... ──►│
  │                       │       (parsed locally)  │
  │◄──────────────────────┼────── 200 OK ───────────┤
```

Tokens are validated **locally** in each service using the shared secret — no
Auth Service roundtrip needed on every API call. The `POST /api/auth/validate`
endpoint exists for explicit server-side verification when needed.

**Token payload:**
```json
{ "user_id": 1, "role": "user", "exp": 1234567890, "iat": 1234567890 }
```

**Request header:**
```
Authorization: Bearer <token>
```

---

## Project Structure

```
mangalib/
├── docker-compose.yml          ← Starts everything: 3 services + 3 DBs + frontend
├── Makefile                    ← Developer shortcuts
├── README.md
│
├── migrations/                 ← golang-migrate SQL files
│   ├── auth/
│   │   ├── 000001_create_users.up.sql
│   │   └── 000001_create_users.down.sql
│   ├── manga/
│   │   ├── 000001_create_genres_tags.{up,down}.sql
│   │   ├── 000002_create_manga_chapters.{up,down}.sql
│   │   └── 000003_create_social_tables.{up,down}.sql
│   └── comment/
│       └── 000001_create_comments.{up,down}.sql
│
├── frontend/                   ← Nginx static UI (port 3000)
│   ├── index.html
│   └── Dockerfile
│
├── auth-service/
│   ├── Dockerfile
│   ├── go.mod                  ← includes golang-migrate, testify
│   ├── tests/
│   │   └── auth_test.go        ← 12 unit tests
│   └── internal/
│       ├── config/             ← env-based config
│       ├── database/           ← connection + AutoMigrate
│       ├── handlers/           ← HTTP layer (thin)
│       ├── middleware/         ← JWT auth, CORS
│       ├── models/             ← GORM structs
│       ├── repositories/       ← DB queries
│       ├── routes/             ← Gin router + DI wiring
│       ├── services/           ← business logic
│       ├── utils/              ← JWT, bcrypt, responses
│       └── validators/         ← request binding DTOs
│
├── manga-service/
│   ├── Dockerfile
│   ├── go.mod                  ← includes resty/v2, golang-migrate
│   ├── tests/
│   │   └── manga_test.go       ← 12 unit tests
│   └── internal/
│       ├── client/             ← Resty v2 AuthClient ← NEW
│       ├── config/             ← includes AuthServiceURL ← NEW
│       ├── handlers/
│       │   └── user_proxy_handler.go  ← Resty demo endpoint ← NEW
│       ├── [... all other layers ...]
│       └── storage/            ← file upload utilities
│
└── comment-service/
    ├── Dockerfile
    ├── go.mod                  ← includes resty/v2, golang-migrate
    ├── tests/
    │   └── comment_test.go     ← 12 unit tests
    └── internal/
        ├── client/             ← Resty v2 AuthClient ← UPGRADED
        └── [... all other layers ...]
```

---

## Database Migrations

The project ships both **GORM AutoMigrate** (for development convenience) and
**golang-migrate SQL files** (for production-grade schema versioning).

### Migration files location

```
migrations/
├── auth/
│   ├── 000001_create_users.up.sql
│   └── 000001_create_users.down.sql
├── manga/
│   ├── 000001_create_genres_tags.{up,down}.sql
│   ├── 000002_create_manga_chapters.{up,down}.sql
│   └── 000003_create_social_tables.{up,down}.sql
└── comment/
    └── 000001_create_comments.{up,down}.sql
```

### Running migrations manually

Install the CLI:
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Run migrations up:
```bash
# Auth Service
migrate -path migrations/auth \
        -database "postgres://auth_user:auth_secret@localhost:5433/auth_db?sslmode=disable" up

# Manga Service
migrate -path migrations/manga \
        -database "postgres://manga_user:manga_secret@localhost:5434/manga_db?sslmode=disable" up

# Comment Service
migrate -path migrations/comment \
        -database "postgres://comment_user:comment_secret@localhost:5435/comment_db?sslmode=disable" up
```

Roll back one step:
```bash
migrate -path migrations/auth \
        -database "postgres://auth_user:auth_secret@localhost:5433/auth_db?sslmode=disable" down 1
```

> **Note:** In Docker Compose development mode, GORM AutoMigrate runs automatically on startup.
> The SQL migration files provide explicit schema versioning for production deployments.

---

## Quick Start — Docker

**Prerequisites:** Docker Desktop installed and running.

```bash
# 1. Clone the project
git clone <repo-url>
cd mangalib

# 2. Start everything
docker-compose up --build
```

That's it. Docker Compose will:
1. Start three PostgreSQL databases
2. Build and start Auth Service (waits for auth-db healthy)
3. Build and start Manga Service (waits for manga-db + auth-service healthy)
4. Build and start Comment Service (waits for comment-db + auth-service healthy)
5. Build and start the Frontend (nginx)

| URL | Description |
|---|---|
| http://localhost:3000 | Frontend dashboard |
| http://localhost:8081/health | Auth Service health |
| http://localhost:8082/health | Manga Service health |
| http://localhost:8083/health | Comment Service health |

### Useful commands

```bash
# View logs
docker-compose logs -f

# Stop everything
docker-compose down

# Stop and wipe volumes
docker-compose down -v

# Rebuild a single service
docker-compose up --build auth-service
```

---

## Running Locally (without Docker)

You need a running PostgreSQL instance with three databases:

```sql
CREATE DATABASE auth_db;
CREATE DATABASE manga_db;
CREATE DATABASE comment_db;
```

Then start each service:

```bash
# Terminal 1
cd auth-service
cp .env.example .env   # edit DB credentials
go run ./cmd/server

# Terminal 2
cd manga-service
go run ./cmd/server

# Terminal 3
cd comment-service
go run ./cmd/server
```

---

## API Overview

### 🔐 Auth Service — `http://localhost:8081`

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/api/auth/register` | ❌ | Register new user |
| `POST` | `/api/auth/login` | ❌ | Login → returns JWT |
| `POST` | `/api/auth/validate` | ❌ | Validate token (inter-service) |
| `GET` | `/api/users/:id` | ❌ | Public user profile |
| `GET` | `/api/users/me` | 🔑 | My profile |
| `PATCH` | `/api/users/me` | 🔑 | Update bio |

### 📚 Manga Service — `http://localhost:8082`

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/api/manga` | ❌ | List manga (paginated) |
| `GET` | `/api/manga/:id` | ❌ | Get manga by ID |
| `GET` | `/api/manga/slug/:slug` | ❌ | Get manga by slug |
| `POST` | `/api/manga` | 🔑 Admin | Create manga |
| `PUT` | `/api/manga/:id` | 🔑 Admin | Update manga |
| `DELETE` | `/api/manga/:id` | 🔑 Admin | Delete manga |
| `POST` | `/api/manga/:id/cover` | 🔑 Admin | Upload cover |
| `GET` | `/api/manga/:id/chapters` | ❌ | List chapters |
| `GET` | `/api/chapters/:id` | ❌ | Get chapter + pages |
| `POST` | `/api/manga/:id/chapters` | 🔑 Admin | Create chapter |
| `POST` | `/api/chapters/:id/pages` | 🔑 Admin | Upload pages |
| `GET` | `/api/genres` | ❌ | List genres |
| `GET` | `/api/tags` | ❌ | List tags |
| `GET` | `/api/bookmarks` | 🔑 | My bookmarks |
| `POST` | `/api/bookmarks` | 🔑 | Add/update bookmark |
| `DELETE` | `/api/bookmarks/:mangaId` | 🔑 | Remove bookmark |
| `POST` | `/api/ratings` | 🔑 | Rate manga (1–10) |
| `POST` | `/api/progress` | 🔑 | Save reading progress |
| `GET` | `/api/progress` | 🔑 | Reading history |
| `GET` | `/api/users/:id/profile` | ❌ | **Resty v2 proxy → Auth** |

### 💬 Comment Service — `http://localhost:8083`

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/api/comments?manga_id=1` | ❌ | List comments (with author via **Resty v2**) |
| `POST` | `/api/comments` | 🔑 | Post comment |
| `DELETE` | `/api/comments/:id` | 🔑 | Delete comment |
| `POST` | `/api/comments/:id/like` | 🔑 | Toggle like |

### Pagination query params (manga listing)

| Param | Type | Default | Description |
|---|---|---|---|
| `page` | int | 1 | Page number |
| `limit` | int | 20 | Items per page (max 100) |
| `search` | string | — | Search by title |
| `status` | string | — | `ongoing` / `completed` / `hiatus` / `cancelled` |
| `genre_id` | uint | — | Filter by genre |
| `sort_by` | string | `created_at` | `views` / `title` / `created_at` |
| `sort_dir` | string | `desc` | `asc` / `desc` |

### Response format

```json
// Success
{ "success": true, "data": { ... } }

// Paginated
{
  "success": true,
  "data": [ ... ],
  "pagination": { "page": 1, "limit": 20, "total": 150, "total_pages": 8 }
}

// Error
{ "success": false, "message": "invalid credentials" }
```

---

## Unit Tests

The project contains **36 unit tests** across all three services (12 per service).

Tests cover:
- JWT generation, parsing, expiry, wrong secret
- Password hashing round-trip
- Auth middleware (valid token, missing header, invalid token)
- Register/login input validation
- Manga CRUD request validation
- Admin-only route enforcement
- Pagination defaults, custom values, clamping
- Comment creation validation (missing fields, body length)
- Response helper structure (OK, NotFound, Paginated)
- Health endpoints

### Run all tests

```bash
# Auth Service
cd auth-service
go test ./tests/... -v

# Manga Service
cd manga-service
go test ./tests/... -v

# Comment Service
cd comment-service
go test ./tests/... -v
```

### Run with coverage

```bash
go test ./tests/... -cover
```

### Example output

```
--- PASS: TestGenerateToken_Valid (0.00s)
--- PASS: TestParseToken_WrongSecret (0.00s)
--- PASS: TestParseToken_Expired (0.00s)
--- PASS: TestPasswordHashing (0.00s)
--- PASS: TestAuthMiddleware_MissingHeader (0.00s)
--- PASS: TestAuthMiddleware_InvalidToken (0.00s)
--- PASS: TestAuthMiddleware_ValidToken (0.00s)
--- PASS: TestRegisterHandler_ShortPassword (0.00s)
--- PASS: TestRegisterHandler_InvalidEmail (0.00s)
--- PASS: TestLoginHandler_EmptyBody (0.00s)
--- PASS: TestValidateHandler_MissingToken (0.00s)
--- PASS: TestValidateHandler_ValidToken (0.00s)
PASS  coverage: 84.2% of statements
```

---

## Screenshots

> Replace these placeholders with actual screenshots before the exam.

### Platform Dashboard (Frontend)
![Frontend Dashboard](docs/screenshots/frontend.png)

### Auth Service — Register
![Register](docs/screenshots/register.png)

### Manga Service — List Manga
![Manga List](docs/screenshots/manga-list.png)

### Comment Service — Comments with Author Info
![Comments](docs/screenshots/comments.png)

### Docker Compose — All Services Running
![Docker](docs/screenshots/docker.png)

---

## Environment Variables

Each service reads from environment variables (or a `.env` file in development).

### Shared across all services
| Variable | Description |
|---|---|
| `JWT_SECRET` | Shared signing secret for JWT tokens |

### Auth Service
| Variable | Default |
|---|---|
| `APP_PORT` | `8081` |
| `DB_HOST` | `localhost` |
| `DB_USER` | `auth_user` |
| `DB_PASSWORD` | `auth_secret` |
| `DB_NAME` | `auth_db` |
| `JWT_EXPIRY_HOURS` | `24` |

### Manga Service
| Variable | Default |
|---|---|
| `APP_PORT` | `8082` |
| `AUTH_SERVICE_URL` | `http://localhost:8081` |
| `UPLOAD_DIR` | `./uploads` |
| `MAX_UPLOAD_SIZE_MB` | `10` |

### Comment Service
| Variable | Default |
|---|---|
| `APP_PORT` | `8083` |
| `AUTH_SERVICE_URL` | `http://localhost:8081` |

---

*Built with Go · Gin · GORM · PostgreSQL · Resty v2 · Docker*

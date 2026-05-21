-- Manga Service Migration: 000003_create_social_tables

CREATE TABLE IF NOT EXISTS bookmarks (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    BIGINT      NOT NULL,
    manga_id   BIGINT      NOT NULL REFERENCES manga(id) ON DELETE CASCADE,
    status     VARCHAR(30) NOT NULL DEFAULT 'plan_to_read',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, manga_id)
);

CREATE INDEX IF NOT EXISTS idx_bookmarks_user_id ON bookmarks(user_id);

CREATE TABLE IF NOT EXISTS ratings (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    BIGINT      NOT NULL,
    manga_id   BIGINT      NOT NULL REFERENCES manga(id) ON DELETE CASCADE,
    score      SMALLINT    NOT NULL CHECK (score >= 1 AND score <= 10),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, manga_id)
);

CREATE TABLE IF NOT EXISTS reading_progress (
    id          BIGSERIAL   PRIMARY KEY,
    user_id     BIGINT      NOT NULL,
    chapter_id  BIGINT      NOT NULL REFERENCES chapters(id) ON DELETE CASCADE,
    manga_id    BIGINT      NOT NULL REFERENCES manga(id)    ON DELETE CASCADE,
    page_number INTEGER     NOT NULL DEFAULT 1,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, chapter_id)
);

CREATE INDEX IF NOT EXISTS idx_progress_user_id ON reading_progress(user_id);

-- Comment Service Migration: 000001_create_comments

CREATE TABLE IF NOT EXISTS comments (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    BIGINT      NOT NULL,
    manga_id   BIGINT      NOT NULL,
    parent_id  BIGINT      REFERENCES comments(id) ON DELETE CASCADE,
    body       TEXT        NOT NULL,
    likes      INTEGER     NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_comments_manga_id   ON comments(manga_id)              WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_comments_user_id    ON comments(user_id)               WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_comments_parent_id  ON comments(parent_id)             WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_comments_deleted_at ON comments(deleted_at);

CREATE TABLE IF NOT EXISTS comment_likes (
    user_id    BIGINT      NOT NULL,
    comment_id BIGINT      NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, comment_id)
);

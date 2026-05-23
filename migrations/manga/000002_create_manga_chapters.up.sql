
CREATE TABLE IF NOT EXISTS manga (
    id          BIGSERIAL    PRIMARY KEY,
    title       VARCHAR(500) NOT NULL,
    alt_title   VARCHAR(500),
    slug        VARCHAR(500) NOT NULL UNIQUE,
    description TEXT,
    cover       VARCHAR(500),
    author      VARCHAR(255),
    artist      VARCHAR(255),
    status      VARCHAR(20)  NOT NULL DEFAULT 'ongoing',
    year        INTEGER,
    age_rating  VARCHAR(10)  NOT NULL DEFAULT '13+',
    views       BIGINT       NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_manga_slug       ON manga(slug)       WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_manga_title      ON manga(title)      WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_manga_views      ON manga(views DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_manga_deleted_at ON manga(deleted_at);

-- Many-to-many join tables
CREATE TABLE IF NOT EXISTS manga_genres (
    manga_id BIGINT NOT NULL REFERENCES manga(id) ON DELETE CASCADE,
    genre_id BIGINT NOT NULL REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (manga_id, genre_id)
);

CREATE TABLE IF NOT EXISTS manga_tags (
    manga_id BIGINT NOT NULL REFERENCES manga(id) ON DELETE CASCADE,
    tag_id   BIGINT NOT NULL REFERENCES tags(id)  ON DELETE CASCADE,
    PRIMARY KEY (manga_id, tag_id)
);

CREATE TABLE IF NOT EXISTS chapters (
    id         BIGSERIAL    PRIMARY KEY,
    manga_id   BIGINT       NOT NULL REFERENCES manga(id) ON DELETE CASCADE,
    number     NUMERIC(8,1) NOT NULL,
    title      VARCHAR(500),
    volume     INTEGER      NOT NULL DEFAULT 0,
    views      BIGINT       NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_chapters_manga_id ON chapters(manga_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS pages (
    id         BIGSERIAL    PRIMARY KEY,
    chapter_id BIGINT       NOT NULL REFERENCES chapters(id) ON DELETE CASCADE,
    number     INTEGER      NOT NULL,
    image_url  VARCHAR(500) NOT NULL,
    width      INTEGER,
    height     INTEGER
);

CREATE INDEX IF NOT EXISTS idx_pages_chapter_id ON pages(chapter_id);

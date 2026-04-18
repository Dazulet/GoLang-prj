CREATE TABLE manga_genres (
    manga_id INTEGER REFERENCES manga(id) ON DELETE CASCADE,
    genre_id INTEGER REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (manga_id, genre_id)
);
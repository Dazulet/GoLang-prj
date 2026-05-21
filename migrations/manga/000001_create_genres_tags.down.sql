-- Manga Service Migration: 000001_create_genres_tags (rollback)
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS genres;

-- Manga Service Migration: 000002_create_manga_chapters (rollback)
DROP TABLE IF EXISTS pages;
DROP TABLE IF EXISTS chapters;
DROP TABLE IF EXISTS manga_tags;
DROP TABLE IF EXISTS manga_genres;
DROP TABLE IF EXISTS manga;

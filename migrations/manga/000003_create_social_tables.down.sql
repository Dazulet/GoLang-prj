-- Manga Service Migration: 000003_create_social_tables (rollback)
DROP TABLE IF EXISTS reading_progress;
DROP TABLE IF EXISTS ratings;
DROP TABLE IF EXISTS bookmarks;

-- Comment Service Migration: 000001_create_comments (rollback)
DROP TABLE IF EXISTS comment_likes;
DROP TABLE IF EXISTS comments;

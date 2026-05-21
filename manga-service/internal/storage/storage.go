package storage

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mangalib/manga-service/internal/config"
)

func EnsureDirectories(cfg *config.Config) error {
	for _, sub := range []string{"covers", "avatars", "chapters"} {
		dir := filepath.Join(cfg.UploadDir, sub)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("mkdir %s: %w", dir, err)
		}
	}
	return nil
}

func SaveFile(c *gin.Context, file *multipart.FileHeader, subDir, uploadDir string) (string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExt(ext) {
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dir := filepath.Join(uploadDir, subDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	dst := filepath.Join(dir, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		return "", err
	}
	return filepath.Join("uploads", subDir, filename), nil
}

func SaveChapterFiles(c *gin.Context, files []*multipart.FileHeader, chapterID uint, uploadDir string) ([]string, error) {
	subDir := filepath.Join("chapters", fmt.Sprintf("%d", chapterID))
	dir := filepath.Join(uploadDir, subDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	paths := make([]string, 0, len(files))
	for i, f := range files {
		ext := strings.ToLower(filepath.Ext(f.Filename))
		if !allowedExt(ext) {
			continue
		}
		dst := filepath.Join(dir, fmt.Sprintf("%04d%s", i+1, ext))
		if err := c.SaveUploadedFile(f, dst); err != nil {
			return nil, err
		}
		paths = append(paths, dst)
	}
	return paths, nil
}

func allowedExt(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif":
		return true
	}
	return false
}

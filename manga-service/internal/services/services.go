package services

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/mangalib/manga-service/internal/models"
	"github.com/mangalib/manga-service/internal/repositories"
	"github.com/mangalib/manga-service/internal/utils"
	"github.com/mangalib/manga-service/internal/validators"
)

// ── MangaService ──────────────────────────────────────────────────────────────

type MangaService struct {
	manga  *repositories.MangaRepository
	genres *repositories.GenreRepository
	tags   *repositories.TagRepository
}

func NewMangaService(manga *repositories.MangaRepository, genres *repositories.GenreRepository, tags *repositories.TagRepository) *MangaService {
	return &MangaService{manga: manga, genres: genres, tags: tags}
}

func (s *MangaService) Create(req validators.CreateMangaRequest) (*models.Manga, error) {
	m := models.Manga{
		Title:       req.Title,
		AltTitle:    req.AltTitle,
		Slug:        utils.Slugify(req.Title),
		Description: req.Description,
		Author:      req.Author,
		Artist:      req.Artist,
		Status:      models.MangaStatus(req.Status),
		Year:        req.Year,
		AgeRating:   req.AgeRating,
	}
	if m.Status == "" {
		m.Status = models.StatusOngoing
	}
	if m.AgeRating == "" {
		m.AgeRating = "13+"
	}
	if err := s.manga.Create(&m); err != nil {
		return nil, err
	}
	if len(req.GenreIDs) > 0 {
		genres, _ := s.genres.FindByIDs(req.GenreIDs)
		_ = s.manga.ReplaceGenres(&m, genres)
	}
	if len(req.TagIDs) > 0 {
		tags, _ := s.tags.FindByIDs(req.TagIDs)
		_ = s.manga.ReplaceTags(&m, tags)
	}
	return s.manga.FindByID(m.ID)
}

func (s *MangaService) GetByID(id uint) (*models.Manga, error) {
	m, err := s.manga.FindByID(id)
	if err != nil {
		return nil, errors.New("manga not found")
	}
	s.enrichRating(m)
	_ = s.manga.IncrementViews(id)
	return m, nil
}

func (s *MangaService) GetBySlug(slug string) (*models.Manga, error) {
	m, err := s.manga.FindBySlug(slug)
	if err != nil {
		return nil, errors.New("manga not found")
	}
	s.enrichRating(m)
	_ = s.manga.IncrementViews(m.ID)
	return m, nil
}

func (s *MangaService) List(f repositories.MangaFilter, page, limit int) ([]models.Manga, int64, error) {
	offset := (page - 1) * limit
	list, total, err := s.manga.List(f, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	for i := range list {
		s.enrichRating(&list[i])
	}
	return list, total, nil
}

func (s *MangaService) Update(id uint, req validators.UpdateMangaRequest) (*models.Manga, error) {
	m, err := s.manga.FindByID(id)
	if err != nil {
		return nil, errors.New("manga not found")
	}
	if req.Title != nil {
		m.Title = *req.Title
		m.Slug = utils.Slugify(*req.Title)
	}
	if req.AltTitle != nil {
		m.AltTitle = *req.AltTitle
	}
	if req.Description != nil {
		m.Description = *req.Description
	}
	if req.Author != nil {
		m.Author = *req.Author
	}
	if req.Artist != nil {
		m.Artist = *req.Artist
	}
	if req.Status != nil {
		m.Status = models.MangaStatus(*req.Status)
	}
	if req.Year != nil {
		m.Year = *req.Year
	}
	if req.AgeRating != nil {
		m.AgeRating = *req.AgeRating
	}
	if err := s.manga.Update(m); err != nil {
		return nil, err
	}
	if req.GenreIDs != nil {
		genres, _ := s.genres.FindByIDs(req.GenreIDs)
		_ = s.manga.ReplaceGenres(m, genres)
	}
	if req.TagIDs != nil {
		tags, _ := s.tags.FindByIDs(req.TagIDs)
		_ = s.manga.ReplaceTags(m, tags)
	}
	return s.manga.FindByID(id)
}

func (s *MangaService) Delete(id uint) error {
	if _, err := s.manga.FindByID(id); err != nil {
		return errors.New("manga not found")
	}
	return s.manga.Delete(id)
}

func (s *MangaService) UpdateCover(id uint, coverURL string) error {
	m, err := s.manga.FindByID(id)
	if err != nil {
		return errors.New("manga not found")
	}
	m.Cover = coverURL
	return s.manga.Update(m)
}

func (s *MangaService) enrichRating(m *models.Manga) {
	avg, count, err := s.manga.AggregateRating(m.ID)
	if err == nil {
		m.AvgRating = avg
		m.TotalRatings = count
	}
}

// ── ChapterService ────────────────────────────────────────────────────────────

type ChapterService struct {
	chapters *repositories.ChapterRepository
	pages    *repositories.PageRepository
	manga    *repositories.MangaRepository
}

func NewChapterService(chapters *repositories.ChapterRepository, pages *repositories.PageRepository, manga *repositories.MangaRepository) *ChapterService {
	return &ChapterService{chapters: chapters, pages: pages, manga: manga}
}

func (s *ChapterService) Create(mangaID uint, req validators.CreateChapterRequest) (*models.Chapter, error) {
	if _, err := s.manga.FindByID(mangaID); err != nil {
		return nil, errors.New("manga not found")
	}
	ch := models.Chapter{MangaID: mangaID, Number: req.Number, Title: req.Title, Volume: req.Volume}
	return &ch, s.chapters.Create(&ch)
}

func (s *ChapterService) GetByID(id uint) (*models.Chapter, error) {
	ch, err := s.chapters.FindByID(id)
	if err != nil {
		return nil, errors.New("chapter not found")
	}
	_ = s.chapters.IncrementViews(id)
	return ch, nil
}

func (s *ChapterService) ListByManga(mangaID uint, page, limit int) ([]models.Chapter, int64, error) {
	return s.chapters.ListByManga(mangaID, limit, (page-1)*limit)
}

func (s *ChapterService) Update(id uint, req validators.UpdateChapterRequest) (*models.Chapter, error) {
	ch, err := s.chapters.FindByID(id)
	if err != nil {
		return nil, errors.New("chapter not found")
	}
	if req.Number != nil {
		ch.Number = *req.Number
	}
	if req.Title != nil {
		ch.Title = *req.Title
	}
	if req.Volume != nil {
		ch.Volume = *req.Volume
	}
	return ch, s.chapters.Update(ch)
}

func (s *ChapterService) Delete(id uint) error {
	if _, err := s.chapters.FindByID(id); err != nil {
		return errors.New("chapter not found")
	}
	return s.chapters.Delete(id)
}

func (s *ChapterService) AddPages(chapterID, mangaID uint, filePaths []string, baseURL string) ([]models.Page, error) {
	if _, err := s.chapters.FindByID(chapterID); err != nil {
		return nil, errors.New("chapter not found")
	}
	_ = s.pages.DeleteByChapter(chapterID)
	pages := make([]models.Page, 0, len(filePaths))
	for i, fp := range filePaths {
		imageURL := fmt.Sprintf("%s/uploads/chapters/%d/%s", baseURL, chapterID, filepath.Base(fp))
		pages = append(pages, models.Page{ChapterID: chapterID, Number: i + 1, ImageURL: imageURL})
	}
	return pages, s.pages.CreateBatch(pages)
}

// ── GenreService ──────────────────────────────────────────────────────────────

type GenreService struct{ repo *repositories.GenreRepository }

func NewGenreService(repo *repositories.GenreRepository) *GenreService {
	return &GenreService{repo: repo}
}

func (s *GenreService) Create(req validators.CreateGenreRequest) (*models.Genre, error) {
	g := models.Genre{Name: req.Name, Slug: utils.Slugify(req.Name)}
	return &g, s.repo.Create(&g)
}

func (s *GenreService) List() ([]models.Genre, error) { return s.repo.FindAll() }

// ── TagService ────────────────────────────────────────────────────────────────

type TagService struct{ repo *repositories.TagRepository }

func NewTagService(repo *repositories.TagRepository) *TagService { return &TagService{repo: repo} }

func (s *TagService) Create(req validators.CreateTagRequest) (*models.Tag, error) {
	t := models.Tag{Name: req.Name, Slug: utils.Slugify(req.Name)}
	return &t, s.repo.Create(&t)
}

func (s *TagService) List() ([]models.Tag, error) { return s.repo.FindAll() }

// ── BookmarkService ───────────────────────────────────────────────────────────

type BookmarkService struct {
	repo *repositories.BookmarkRepository
}

func NewBookmarkService(repo *repositories.BookmarkRepository) *BookmarkService {
	return &BookmarkService{repo: repo}
}

func (s *BookmarkService) Upsert(userID uint, req validators.UpsertBookmarkRequest) (*models.Bookmark, error) {
	b := models.Bookmark{UserID: userID, MangaID: req.MangaID, Status: models.BookmarkStatus(req.Status)}
	return &b, s.repo.Upsert(&b)
}

func (s *BookmarkService) Remove(userID, mangaID uint) error {
	return s.repo.Delete(userID, mangaID)
}

func (s *BookmarkService) ListByUser(userID uint, page, limit int) ([]models.Bookmark, int64, error) {
	return s.repo.ListByUser(userID, limit, (page-1)*limit)
}

// ── RatingService ─────────────────────────────────────────────────────────────

type RatingService struct {
	repo *repositories.RatingRepository
}

func NewRatingService(repo *repositories.RatingRepository) *RatingService {
	return &RatingService{repo: repo}
}

func (s *RatingService) Upsert(userID uint, req validators.UpsertRatingRequest) (*models.Rating, error) {
	r := models.Rating{UserID: userID, MangaID: req.MangaID, Score: req.Score}
	return &r, s.repo.Upsert(&r)
}

func (s *RatingService) GetByManga(userID, mangaID uint) (*models.Rating, error) {
	return s.repo.FindByUserAndManga(userID, mangaID)
}

// ── ProgressService ───────────────────────────────────────────────────────────

type ProgressService struct {
	repo *repositories.ProgressRepository
}

func NewProgressService(repo *repositories.ProgressRepository) *ProgressService {
	return &ProgressService{repo: repo}
}

func (s *ProgressService) Update(userID uint, req validators.UpdateProgressRequest) error {
	p := models.ReadingProgress{UserID: userID, ChapterID: req.ChapterID, MangaID: req.MangaID, PageNumber: req.PageNumber}
	return s.repo.Upsert(&p)
}

func (s *ProgressService) History(userID uint) ([]models.ReadingProgress, error) {
	return s.repo.ListByUser(userID)
}

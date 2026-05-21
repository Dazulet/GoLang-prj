package repositories

import (
	"github.com/mangalib/manga-service/internal/models"
	"gorm.io/gorm"
)

// ── MangaFilter ───────────────────────────────────────────────────────────────

type MangaFilter struct {
	Search  string
	Status  string
	GenreID uint
	SortBy  string
	SortDir string
}

// ── MangaRepository ───────────────────────────────────────────────────────────

type MangaRepository struct{ db *gorm.DB }

func NewMangaRepository(db *gorm.DB) *MangaRepository { return &MangaRepository{db: db} }

func (r *MangaRepository) Create(m *models.Manga) error { return r.db.Create(m).Error }

func (r *MangaRepository) FindByID(id uint) (*models.Manga, error) {
	var m models.Manga
	err := r.db.Preload("Genres").Preload("Tags").First(&m, id).Error
	return &m, err
}

func (r *MangaRepository) FindBySlug(slug string) (*models.Manga, error) {
	var m models.Manga
	err := r.db.Preload("Genres").Preload("Tags").Where("slug = ?", slug).First(&m).Error
	return &m, err
}

func (r *MangaRepository) List(f MangaFilter, limit, offset int) ([]models.Manga, int64, error) {
	q := r.db.Model(&models.Manga{}).Preload("Genres").Preload("Tags")

	if f.Search != "" {
		q = q.Where("title ILIKE ? OR alt_title ILIKE ?", "%"+f.Search+"%", "%"+f.Search+"%")
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.GenreID != 0 {
		q = q.Joins("JOIN manga_genres ON manga_genres.manga_id = manga.id").
			Where("manga_genres.genre_id = ?", f.GenreID)
	}

	col := "created_at"
	if f.SortBy == "views" {
		col = "views"
	} else if f.SortBy == "title" {
		col = "title"
	}
	dir := "desc"
	if f.SortDir == "asc" {
		dir = "asc"
	}
	q = q.Order(col + " " + dir)

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []models.Manga
	err := q.Limit(limit).Offset(offset).Find(&list).Error
	return list, total, err
}

func (r *MangaRepository) Update(m *models.Manga) error { return r.db.Save(m).Error }

func (r *MangaRepository) Delete(id uint) error { return r.db.Delete(&models.Manga{}, id).Error }

func (r *MangaRepository) IncrementViews(id uint) error {
	return r.db.Model(&models.Manga{}).Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + 1")).Error
}

func (r *MangaRepository) ReplaceGenres(m *models.Manga, genres []models.Genre) error {
	return r.db.Model(m).Association("Genres").Replace(genres)
}

func (r *MangaRepository) ReplaceTags(m *models.Manga, tags []models.Tag) error {
	return r.db.Model(m).Association("Tags").Replace(tags)
}

func (r *MangaRepository) AggregateRating(mangaID uint) (avg float64, count int, err error) {
	row := r.db.Model(&models.Rating{}).
		Select("COALESCE(AVG(score),0), COUNT(*)").
		Where("manga_id = ?", mangaID).Row()
	err = row.Scan(&avg, &count)
	return
}

// ── ChapterRepository ─────────────────────────────────────────────────────────

type ChapterRepository struct{ db *gorm.DB }

func NewChapterRepository(db *gorm.DB) *ChapterRepository { return &ChapterRepository{db: db} }

func (r *ChapterRepository) Create(ch *models.Chapter) error { return r.db.Create(ch).Error }

func (r *ChapterRepository) FindByID(id uint) (*models.Chapter, error) {
	var ch models.Chapter
	err := r.db.Preload("Pages", func(db *gorm.DB) *gorm.DB {
		return db.Order("pages.number ASC")
	}).First(&ch, id).Error
	return &ch, err
}

func (r *ChapterRepository) ListByManga(mangaID uint, limit, offset int) ([]models.Chapter, int64, error) {
	var total int64
	r.db.Model(&models.Chapter{}).Where("manga_id = ?", mangaID).Count(&total)
	var chapters []models.Chapter
	err := r.db.Where("manga_id = ?", mangaID).Order("number ASC").
		Limit(limit).Offset(offset).Find(&chapters).Error
	return chapters, total, err
}

func (r *ChapterRepository) Update(ch *models.Chapter) error { return r.db.Save(ch).Error }

func (r *ChapterRepository) Delete(id uint) error { return r.db.Delete(&models.Chapter{}, id).Error }

func (r *ChapterRepository) IncrementViews(id uint) error {
	return r.db.Model(&models.Chapter{}).Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + 1")).Error
}

func (r *ChapterRepository) CountByManga(mangaID uint) (int64, error) {
	var n int64
	err := r.db.Model(&models.Chapter{}).Where("manga_id = ?", mangaID).Count(&n).Error
	return n, err
}

// ── PageRepository ────────────────────────────────────────────────────────────

type PageRepository struct{ db *gorm.DB }

func NewPageRepository(db *gorm.DB) *PageRepository { return &PageRepository{db: db} }

func (r *PageRepository) CreateBatch(pages []models.Page) error { return r.db.Create(&pages).Error }

func (r *PageRepository) DeleteByChapter(chapterID uint) error {
	return r.db.Where("chapter_id = ?", chapterID).Delete(&models.Page{}).Error
}

// ── GenreRepository ───────────────────────────────────────────────────────────

type GenreRepository struct{ db *gorm.DB }

func NewGenreRepository(db *gorm.DB) *GenreRepository { return &GenreRepository{db: db} }

func (r *GenreRepository) Create(g *models.Genre) error { return r.db.Create(g).Error }

func (r *GenreRepository) FindAll() ([]models.Genre, error) {
	var genres []models.Genre
	return genres, r.db.Order("name ASC").Find(&genres).Error
}

func (r *GenreRepository) FindByIDs(ids []uint) ([]models.Genre, error) {
	var genres []models.Genre
	return genres, r.db.Where("id IN ?", ids).Find(&genres).Error
}

// ── TagRepository ─────────────────────────────────────────────────────────────

type TagRepository struct{ db *gorm.DB }

func NewTagRepository(db *gorm.DB) *TagRepository { return &TagRepository{db: db} }

func (r *TagRepository) Create(t *models.Tag) error { return r.db.Create(t).Error }

func (r *TagRepository) FindAll() ([]models.Tag, error) {
	var tags []models.Tag
	return tags, r.db.Order("name ASC").Find(&tags).Error
}

func (r *TagRepository) FindByIDs(ids []uint) ([]models.Tag, error) {
	var tags []models.Tag
	return tags, r.db.Where("id IN ?", ids).Find(&tags).Error
}

// ── BookmarkRepository ────────────────────────────────────────────────────────

type BookmarkRepository struct{ db *gorm.DB }

func NewBookmarkRepository(db *gorm.DB) *BookmarkRepository { return &BookmarkRepository{db: db} }

func (r *BookmarkRepository) Upsert(b *models.Bookmark) error {
	return r.db.Where(models.Bookmark{UserID: b.UserID, MangaID: b.MangaID}).
		Assign(models.Bookmark{Status: b.Status}).FirstOrCreate(b).Error
}

func (r *BookmarkRepository) Delete(userID, mangaID uint) error {
	return r.db.Where("user_id = ? AND manga_id = ?", userID, mangaID).Delete(&models.Bookmark{}).Error
}

func (r *BookmarkRepository) ListByUser(userID uint, limit, offset int) ([]models.Bookmark, int64, error) {
	var total int64
	r.db.Model(&models.Bookmark{}).Where("user_id = ?", userID).Count(&total)
	var list []models.Bookmark
	err := r.db.Preload("Manga.Genres").Where("user_id = ?", userID).
		Order("updated_at DESC").Limit(limit).Offset(offset).Find(&list).Error
	return list, total, err
}

func (r *BookmarkRepository) FindByUserAndManga(userID, mangaID uint) (*models.Bookmark, error) {
	var b models.Bookmark
	return &b, r.db.Where("user_id = ? AND manga_id = ?", userID, mangaID).First(&b).Error
}

// ── RatingRepository ──────────────────────────────────────────────────────────

type RatingRepository struct{ db *gorm.DB }

func NewRatingRepository(db *gorm.DB) *RatingRepository { return &RatingRepository{db: db} }

func (r *RatingRepository) Upsert(rating *models.Rating) error {
	return r.db.Where(models.Rating{UserID: rating.UserID, MangaID: rating.MangaID}).
		Assign(models.Rating{Score: rating.Score}).FirstOrCreate(rating).Error
}

func (r *RatingRepository) FindByUserAndManga(userID, mangaID uint) (*models.Rating, error) {
	var rating models.Rating
	return &rating, r.db.Where("user_id = ? AND manga_id = ?", userID, mangaID).First(&rating).Error
}

// ── ProgressRepository ────────────────────────────────────────────────────────

type ProgressRepository struct{ db *gorm.DB }

func NewProgressRepository(db *gorm.DB) *ProgressRepository { return &ProgressRepository{db: db} }

func (r *ProgressRepository) Upsert(p *models.ReadingProgress) error {
	return r.db.Where(models.ReadingProgress{UserID: p.UserID, ChapterID: p.ChapterID}).
		Assign(models.ReadingProgress{PageNumber: p.PageNumber, MangaID: p.MangaID}).
		FirstOrCreate(p).Error
}

func (r *ProgressRepository) ListByUser(userID uint) ([]models.ReadingProgress, error) {
	var list []models.ReadingProgress
	err := r.db.Preload("Chapter").Where("user_id = ?", userID).
		Order("updated_at DESC").Find(&list).Error
	return list, err
}

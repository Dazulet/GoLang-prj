package models

import (
	"time"

	"gorm.io/gorm"
)

// ── Genre & Tag ───────────────────────────────────────────────────────────────

type Genre struct {
	ID    uint    `gorm:"primaryKey"                    json:"id"`
	Name  string  `gorm:"uniqueIndex;size:100;not null"  json:"name"`
	Slug  string  `gorm:"uniqueIndex;size:100;not null"  json:"slug"`
	Manga []Manga `gorm:"many2many:manga_genres;"        json:"-"`
}

type Tag struct {
	ID    uint    `gorm:"primaryKey"                    json:"id"`
	Name  string  `gorm:"uniqueIndex;size:100;not null"  json:"name"`
	Slug  string  `gorm:"uniqueIndex;size:100;not null"  json:"slug"`
	Manga []Manga `gorm:"many2many:manga_tags;"          json:"-"`
}

// ── Manga ─────────────────────────────────────────────────────────────────────

type MangaStatus string

const (
	StatusOngoing   MangaStatus = "ongoing"
	StatusCompleted MangaStatus = "completed"
	StatusHiatus    MangaStatus = "hiatus"
	StatusCancelled MangaStatus = "cancelled"
)

type Manga struct {
	ID          uint           `gorm:"primaryKey"                        json:"id"`
	Title       string         `gorm:"size:500;not null;index"           json:"title"`
	AltTitle    string         `gorm:"size:500"                          json:"alt_title,omitempty"`
	Slug        string         `gorm:"uniqueIndex;size:500;not null"     json:"slug"`
	Description string         `gorm:"type:text"                         json:"description,omitempty"`
	Cover       string         `gorm:"size:500"                          json:"cover,omitempty"`
	Author      string         `gorm:"size:255"                          json:"author,omitempty"`
	Artist      string         `gorm:"size:255"                          json:"artist,omitempty"`
	Status      MangaStatus    `gorm:"type:varchar(20);default:'ongoing'" json:"status"`
	Year        int            `                                         json:"year,omitempty"`
	AgeRating   string         `gorm:"size:10;default:'13+'"             json:"age_rating"`
	Views       uint64         `gorm:"default:0;index"                   json:"views"`
	CreatedAt   time.Time      `                                         json:"created_at"`
	UpdatedAt   time.Time      `                                         json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"                             json:"-"`

	Genres   []Genre   `gorm:"many2many:manga_genres;" json:"genres,omitempty"`
	Tags     []Tag     `gorm:"many2many:manga_tags;"   json:"tags,omitempty"`
	Chapters []Chapter `gorm:"foreignKey:MangaID"      json:"chapters,omitempty"`
	Ratings  []Rating  `gorm:"foreignKey:MangaID"      json:"-"`

	// Computed by service layer
	AvgRating     float64 `gorm:"-" json:"avg_rating,omitempty"`
	TotalRatings  int     `gorm:"-" json:"total_ratings,omitempty"`
	TotalChapters int64   `gorm:"-" json:"total_chapters,omitempty"`
}

// ── Chapter ───────────────────────────────────────────────────────────────────

type Chapter struct {
	ID        uint           `gorm:"primaryKey"       json:"id"`
	MangaID   uint           `gorm:"not null;index"   json:"manga_id"`
	Number    float64        `gorm:"not null"         json:"number"`
	Title     string         `gorm:"size:500"         json:"title,omitempty"`
	Volume    int            `                        json:"volume,omitempty"`
	Views     uint64         `gorm:"default:0"        json:"views"`
	CreatedAt time.Time      `                        json:"created_at"`
	UpdatedAt time.Time      `                        json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"            json:"-"`

	Pages []Page `gorm:"foreignKey:ChapterID;constraint:OnDelete:CASCADE" json:"pages,omitempty"`
}

// ── Page ──────────────────────────────────────────────────────────────────────

type Page struct {
	ID        uint   `gorm:"primaryKey"       json:"id"`
	ChapterID uint   `gorm:"not null;index"   json:"chapter_id"`
	Number    int    `gorm:"not null"         json:"number"`
	ImageURL  string `gorm:"size:500;not null" json:"image_url"`
	Width     int    `                        json:"width,omitempty"`
	Height    int    `                        json:"height,omitempty"`
}

// ── Bookmark ──────────────────────────────────────────────────────────────────

type BookmarkStatus string

const (
	BookmarkReading   BookmarkStatus = "reading"
	BookmarkCompleted BookmarkStatus = "completed"
	BookmarkPlanning  BookmarkStatus = "plan_to_read"
	BookmarkDropped   BookmarkStatus = "dropped"
	BookmarkOnHold    BookmarkStatus = "on_hold"
)

type Bookmark struct {
	ID        uint           `gorm:"primaryKey"                                         json:"id"`
	UserID    uint           `gorm:"not null;uniqueIndex:idx_user_manga_bookmark"       json:"user_id"`
	MangaID   uint           `gorm:"not null;uniqueIndex:idx_user_manga_bookmark"       json:"manga_id"`
	Status    BookmarkStatus `gorm:"type:varchar(30);default:'plan_to_read'"            json:"status"`
	CreatedAt time.Time      `                                                          json:"created_at"`
	UpdatedAt time.Time      `                                                          json:"updated_at"`

	Manga Manga `gorm:"foreignKey:MangaID" json:"manga,omitempty"`
}

// ── Rating ────────────────────────────────────────────────────────────────────

type Rating struct {
	ID        uint      `gorm:"primaryKey"                                   json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_user_manga_rating"   json:"user_id"`
	MangaID   uint      `gorm:"not null;uniqueIndex:idx_user_manga_rating"   json:"manga_id"`
	Score     int       `gorm:"not null;check:score >= 1 AND score <= 10"    json:"score"`
	CreatedAt time.Time `                                                    json:"created_at"`
	UpdatedAt time.Time `                                                    json:"updated_at"`
}

// ── Reading Progress ──────────────────────────────────────────────────────────

type ReadingProgress struct {
	ID         uint      `gorm:"primaryKey"                                        json:"id"`
	UserID     uint      `gorm:"not null;uniqueIndex:idx_progress_user_chapter"    json:"user_id"`
	ChapterID  uint      `gorm:"not null;uniqueIndex:idx_progress_user_chapter"    json:"chapter_id"`
	MangaID    uint      `gorm:"not null;index"                                    json:"manga_id"`
	PageNumber int       `gorm:"default:1"                                         json:"page_number"`
	UpdatedAt  time.Time `                                                          json:"updated_at"`

	Chapter Chapter `gorm:"foreignKey:ChapterID" json:"chapter,omitempty"`
}

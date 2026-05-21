package validators

// ── Manga ─────────────────────────────────────────────────────────────────────

type CreateMangaRequest struct {
	Title       string `json:"title"       binding:"required,min=1,max=500"`
	AltTitle    string `json:"alt_title"   binding:"omitempty,max=500"`
	Description string `json:"description"`
	Author      string `json:"author"      binding:"omitempty,max=255"`
	Artist      string `json:"artist"      binding:"omitempty,max=255"`
	Status      string `json:"status"      binding:"omitempty,oneof=ongoing completed hiatus cancelled"`
	Year        int    `json:"year"        binding:"omitempty,min=1900,max=2100"`
	AgeRating   string `json:"age_rating"  binding:"omitempty"`
	GenreIDs    []uint `json:"genre_ids"`
	TagIDs      []uint `json:"tag_ids"`
}

type UpdateMangaRequest struct {
	Title       *string `json:"title"       binding:"omitempty,min=1,max=500"`
	AltTitle    *string `json:"alt_title"   binding:"omitempty,max=500"`
	Description *string `json:"description"`
	Author      *string `json:"author"      binding:"omitempty,max=255"`
	Artist      *string `json:"artist"      binding:"omitempty,max=255"`
	Status      *string `json:"status"      binding:"omitempty,oneof=ongoing completed hiatus cancelled"`
	Year        *int    `json:"year"        binding:"omitempty,min=1900,max=2100"`
	AgeRating   *string `json:"age_rating"  binding:"omitempty"`
	GenreIDs    []uint  `json:"genre_ids"`
	TagIDs      []uint  `json:"tag_ids"`
}

// ── Chapter ───────────────────────────────────────────────────────────────────

type CreateChapterRequest struct {
	Number float64 `json:"number" binding:"min=0"`
	Title  string  `json:"title"  binding:"omitempty,max=500"`
	Volume int     `json:"volume" binding:"omitempty,min=0"`
}

type UpdateChapterRequest struct {
	Number *float64 `json:"number" binding:"omitempty,min=0"`
	Title  *string  `json:"title"  binding:"omitempty,max=500"`
	Volume *int     `json:"volume" binding:"omitempty,min=0"`
}

// ── Bookmark ──────────────────────────────────────────────────────────────────

type UpsertBookmarkRequest struct {
	MangaID uint   `json:"manga_id" binding:"required"`
	Status  string `json:"status"   binding:"required,oneof=reading completed plan_to_read dropped on_hold"`
}

// ── Rating ────────────────────────────────────────────────────────────────────

type UpsertRatingRequest struct {
	MangaID uint `json:"manga_id" binding:"required"`
	Score   int  `json:"score"    binding:"required,min=1,max=10"`
}

// ── Progress ──────────────────────────────────────────────────────────────────

type UpdateProgressRequest struct {
	ChapterID  uint `json:"chapter_id"  binding:"required"`
	MangaID    uint `json:"manga_id"    binding:"required"`
	PageNumber int  `json:"page_number" binding:"required,min=1"`
}

// ── Genre / Tag ───────────────────────────────────────────────────────────────

type CreateGenreRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

type CreateTagRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

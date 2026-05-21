package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mangalib/manga-service/internal/config"
	"github.com/mangalib/manga-service/internal/middleware"
	"github.com/mangalib/manga-service/internal/repositories"
	"github.com/mangalib/manga-service/internal/services"
	"github.com/mangalib/manga-service/internal/storage"
	"github.com/mangalib/manga-service/internal/utils"
	"github.com/mangalib/manga-service/internal/validators"
)

func paramUint(c *gin.Context, key string) (uint, error) {
	n, err := strconv.ParseUint(c.Param(key), 10, 64)
	return uint(n), err
}

// ── MangaHandler ──────────────────────────────────────────────────────────────

type MangaHandler struct {
	svc *services.MangaService
	cfg *config.Config
}

func NewMangaHandler(svc *services.MangaService, cfg *config.Config) *MangaHandler {
	return &MangaHandler{svc: svc, cfg: cfg}
}

func (h *MangaHandler) List(c *gin.Context) {
	page, limit, _ := utils.ParsePagination(c)
	f := repositories.MangaFilter{
		Search:  c.Query("search"),
		Status:  c.Query("status"),
		SortBy:  c.Query("sort_by"),
		SortDir: c.Query("sort_dir"),
	}
	if gid := c.Query("genre_id"); gid != "" {
		id, _ := strconv.ParseUint(gid, 10, 64)
		f.GenreID = uint(id)
	}
	list, total, err := h.svc.List(f, page, limit)
	if err != nil {
		utils.InternalError(c)
		return
	}
	utils.Paginated(c, list, total, page, limit)
}

func (h *MangaHandler) Get(c *gin.Context) {
	id, err := paramUint(c, "id")
	if err != nil {
		utils.BadRequest(c, "invalid id")
		return
	}
	m, err := h.svc.GetByID(id)
	if err != nil {
		utils.NotFound(c, "manga")
		return
	}
	utils.OK(c, m)
}

func (h *MangaHandler) GetBySlug(c *gin.Context) {
	m, err := h.svc.GetBySlug(c.Param("slug"))
	if err != nil {
		utils.NotFound(c, "manga")
		return
	}
	utils.OK(c, m)
}

func (h *MangaHandler) Create(c *gin.Context) {
	var req validators.CreateMangaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	m, err := h.svc.Create(req)
	if err != nil {
		utils.InternalError(c)
		return
	}
	utils.Created(c, m)
}

func (h *MangaHandler) Update(c *gin.Context) {
	id, err := paramUint(c, "id")
	if err != nil {
		utils.BadRequest(c, "invalid id")
		return
	}
	var req validators.UpdateMangaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	m, err := h.svc.Update(id, req)
	if err != nil {
		utils.NotFound(c, "manga")
		return
	}
	utils.OK(c, m)
}

func (h *MangaHandler) Delete(c *gin.Context) {
	id, err := paramUint(c, "id")
	if err != nil {
		utils.BadRequest(c, "invalid id")
		return
	}
	if err := h.svc.Delete(id); err != nil {
		utils.NotFound(c, "manga")
		return
	}
	utils.NoContent(c)
}

func (h *MangaHandler) UploadCover(c *gin.Context) {
	id, err := paramUint(c, "id")
	if err != nil {
		utils.BadRequest(c, "invalid id")
		return
	}
	file, err := c.FormFile("cover")
	if err != nil {
		utils.BadRequest(c, "cover file required")
		return
	}
	if file.Size > h.cfg.MaxUploadSizeMB*1024*1024 {
		utils.BadRequest(c, "file too large")
		return
	}
	path, err := storage.SaveFile(c, file, "covers", h.cfg.UploadDir)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	url := h.cfg.AppBaseURL + "/" + path
	if err := h.svc.UpdateCover(id, url); err != nil {
		utils.NotFound(c, "manga")
		return
	}
	utils.OK(c, gin.H{"cover": url})
}

// ── ChapterHandler ────────────────────────────────────────────────────────────

type ChapterHandler struct {
	svc *services.ChapterService
	cfg *config.Config
}

func NewChapterHandler(svc *services.ChapterService, cfg *config.Config) *ChapterHandler {
	return &ChapterHandler{svc: svc, cfg: cfg}
}

func (h *ChapterHandler) ListByManga(c *gin.Context) {
	mangaID, err := paramUint(c, "id")
	if err != nil {
		utils.BadRequest(c, "invalid manga id")
		return
	}
	page, limit, _ := utils.ParsePagination(c)
	list, total, err := h.svc.ListByManga(mangaID, page, limit)
	if err != nil {
		utils.InternalError(c)
		return
	}
	utils.Paginated(c, list, total, page, limit)
}

func (h *ChapterHandler) Get(c *gin.Context) {
	id, err := paramUint(c, "id")
	if err != nil {
		utils.BadRequest(c, "invalid id")
		return
	}
	ch, err := h.svc.GetByID(id)
	if err != nil {
		utils.NotFound(c, "chapter")
		return
	}
	utils.OK(c, ch)
}

func (h *ChapterHandler) Create(c *gin.Context) {
	mangaID, err := paramUint(c, "id")
	if err != nil {
		utils.BadRequest(c, "invalid manga id")
		return
	}
	var req validators.CreateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	ch, err := h.svc.Create(mangaID, req)
	if err != nil {
		utils.NotFound(c, "manga")
		return
	}
	utils.Created(c, ch)
}

func (h *ChapterHandler) Update(c *gin.Context) {
	id, err := paramUint(c, "id")
	if err != nil {
		utils.BadRequest(c, "invalid id")
		return
	}
	var req validators.UpdateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	ch, err := h.svc.Update(id, req)
	if err != nil {
		utils.NotFound(c, "chapter")
		return
	}
	utils.OK(c, ch)
}

func (h *ChapterHandler) Delete(c *gin.Context) {
	id, err := paramUint(c, "id")
	if err != nil {
		utils.BadRequest(c, "invalid id")
		return
	}
	if err := h.svc.Delete(id); err != nil {
		utils.NotFound(c, "chapter")
		return
	}
	utils.NoContent(c)
}

func (h *ChapterHandler) UploadPages(c *gin.Context) {
	id, err := paramUint(c, "id")
	if err != nil {
		utils.BadRequest(c, "invalid id")
		return
	}
	mangaIDStr := c.PostForm("manga_id")
	mangaID, _ := strconv.ParseUint(mangaIDStr, 10, 64)

	form, err := c.MultipartForm()
	if err != nil {
		utils.BadRequest(c, "multipart form required")
		return
	}
	files := form.File["pages"]
	if len(files) == 0 {
		utils.BadRequest(c, "at least one page required")
		return
	}
	paths, err := storage.SaveChapterFiles(c, files, id, h.cfg.UploadDir)
	if err != nil {
		utils.InternalError(c)
		return
	}
	pages, err := h.svc.AddPages(id, uint(mangaID), paths, h.cfg.AppBaseURL)
	if err != nil {
		utils.InternalError(c)
		return
	}
	utils.Created(c, pages)
}

// ── GenreTagHandler ───────────────────────────────────────────────────────────

type GenreTagHandler struct {
	genres *services.GenreService
	tags   *services.TagService
}

func NewGenreTagHandler(genres *services.GenreService, tags *services.TagService) *GenreTagHandler {
	return &GenreTagHandler{genres: genres, tags: tags}
}

func (h *GenreTagHandler) ListGenres(c *gin.Context) {
	list, err := h.genres.List()
	if err != nil {
		utils.InternalError(c)
		return
	}
	utils.OK(c, list)
}

func (h *GenreTagHandler) CreateGenre(c *gin.Context) {
	var req validators.CreateGenreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	g, err := h.genres.Create(req)
	if err != nil {
		utils.InternalError(c)
		return
	}
	utils.Created(c, g)
}

func (h *GenreTagHandler) ListTags(c *gin.Context) {
	list, err := h.tags.List()
	if err != nil {
		utils.InternalError(c)
		return
	}
	utils.OK(c, list)
}

func (h *GenreTagHandler) CreateTag(c *gin.Context) {
	var req validators.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	t, err := h.tags.Create(req)
	if err != nil {
		utils.InternalError(c)
		return
	}
	utils.Created(c, t)
}

// ── BookmarkHandler ───────────────────────────────────────────────────────────

type BookmarkHandler struct{ svc *services.BookmarkService }

func NewBookmarkHandler(svc *services.BookmarkService) *BookmarkHandler {
	return &BookmarkHandler{svc: svc}
}

func (h *BookmarkHandler) List(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	page, limit, _ := utils.ParsePagination(c)
	list, total, err := h.svc.ListByUser(userID, page, limit)
	if err != nil {
		utils.InternalError(c)
		return
	}
	utils.Paginated(c, list, total, page, limit)
}

func (h *BookmarkHandler) Upsert(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	var req validators.UpsertBookmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	b, err := h.svc.Upsert(userID, req)
	if err != nil {
		utils.InternalError(c)
		return
	}
	utils.OK(c, b)
}

func (h *BookmarkHandler) Remove(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	mangaID, err := paramUint(c, "mangaId")
	if err != nil {
		utils.BadRequest(c, "invalid manga id")
		return
	}
	if err := h.svc.Remove(userID, mangaID); err != nil {
		utils.InternalError(c)
		return
	}
	utils.NoContent(c)
}

// ── RatingHandler ─────────────────────────────────────────────────────────────

type RatingHandler struct{ svc *services.RatingService }

func NewRatingHandler(svc *services.RatingService) *RatingHandler { return &RatingHandler{svc: svc} }

func (h *RatingHandler) Upsert(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	var req validators.UpsertRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	r, err := h.svc.Upsert(userID, req)
	if err != nil {
		utils.InternalError(c)
		return
	}
	utils.OK(c, r)
}

func (h *RatingHandler) GetMine(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	mangaID, err := paramUint(c, "mangaId")
	if err != nil {
		utils.BadRequest(c, "invalid manga id")
		return
	}
	r, err := h.svc.GetByManga(userID, mangaID)
	if err != nil {
		utils.NotFound(c, "rating")
		return
	}
	utils.OK(c, r)
}

// ── ProgressHandler ───────────────────────────────────────────────────────────

type ProgressHandler struct{ svc *services.ProgressService }

func NewProgressHandler(svc *services.ProgressService) *ProgressHandler {
	return &ProgressHandler{svc: svc}
}

func (h *ProgressHandler) Update(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	var req validators.UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.Update(userID, req); err != nil {
		utils.InternalError(c)
		return
	}
	utils.OK(c, gin.H{"message": "progress updated"})
}

func (h *ProgressHandler) History(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	list, err := h.svc.History(userID)
	if err != nil {
		utils.InternalError(c)
		return
	}
	utils.OK(c, list)
}

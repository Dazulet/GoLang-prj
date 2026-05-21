package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mangalib/comment-service/internal/middleware"
	"github.com/mangalib/comment-service/internal/services"
	"github.com/mangalib/comment-service/internal/utils"
	"github.com/mangalib/comment-service/internal/validators"
)

type CommentHandler struct {
	svc *services.CommentService
}

func NewCommentHandler(svc *services.CommentService) *CommentHandler {
	return &CommentHandler{svc: svc}
}

// GET /api/comments?manga_id=1&page=1&limit=20
func (h *CommentHandler) List(c *gin.Context) {
	mangaIDStr := c.Query("manga_id")
	if mangaIDStr == "" {
		utils.BadRequest(c, "manga_id is required")
		return
	}
	mangaID, err := strconv.ParseUint(mangaIDStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid manga_id")
		return
	}

	page, limit, _ := utils.ParsePagination(c)
	list, total, err := h.svc.ListByManga(uint(mangaID), page, limit)
	if err != nil {
		utils.InternalError(c)
		return
	}
	utils.Paginated(c, list, total, page, limit)
}

// POST /api/comments
func (h *CommentHandler) Create(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	var req validators.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	comment, err := h.svc.Create(userID, req)
	if err != nil {
		utils.InternalError(c)
		return
	}
	utils.Created(c, comment)
}

// DELETE /api/comments/:id
func (h *CommentHandler) Delete(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	isAdmin := middleware.CurrentUserRole(c) == "admin"

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid comment id")
		return
	}

	if err := h.svc.Delete(userID, uint(id), isAdmin); err != nil {
		switch err.Error() {
		case "forbidden":
			utils.Forbidden(c)
		case "comment not found":
			utils.NotFound(c, "comment")
		default:
			utils.InternalError(c)
		}
		return
	}
	utils.NoContent(c)
}

// POST /api/comments/:id/like
func (h *CommentHandler) Like(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid comment id")
		return
	}

	liked, err := h.svc.ToggleLike(userID, uint(id))
	if err != nil {
		if err.Error() == "comment not found" {
			utils.NotFound(c, "comment")
		} else {
			utils.InternalError(c)
		}
		return
	}

	action := "unliked"
	if liked {
		action = "liked"
	}
	utils.OK(c, gin.H{"action": action})
}

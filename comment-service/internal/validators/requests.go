package validators

type CreateCommentRequest struct {
	MangaID  uint   `json:"manga_id"  binding:"required"`
	Body     string `json:"body"      binding:"required,min=1,max=2000"`
	ParentID *uint  `json:"parent_id"`
}

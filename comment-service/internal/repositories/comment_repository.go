package repositories

import (
	"github.com/mangalib/comment-service/internal/models"
	"gorm.io/gorm"
)

type CommentRepository struct{ db *gorm.DB }

func NewCommentRepository(db *gorm.DB) *CommentRepository { return &CommentRepository{db: db} }

func (r *CommentRepository) Create(c *models.Comment) error { return r.db.Create(c).Error }

func (r *CommentRepository) FindByID(id uint) (*models.Comment, error) {
	var c models.Comment
	return &c, r.db.First(&c, id).Error
}

func (r *CommentRepository) ListByManga(mangaID uint, limit, offset int) ([]models.Comment, int64, error) {
	var total int64
	r.db.Model(&models.Comment{}).
		Where("manga_id = ? AND parent_id IS NULL", mangaID).
		Count(&total)

	var comments []models.Comment
	err := r.db.
		Preload("Replies").
		Where("manga_id = ? AND parent_id IS NULL", mangaID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&comments).Error
	return comments, total, err
}

func (r *CommentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Comment{}, id).Error
}

func (r *CommentRepository) LikeExists(userID, commentID uint) bool {
	var n int64
	r.db.Model(&models.CommentLike{}).
		Where("user_id = ? AND comment_id = ?", userID, commentID).Count(&n)
	return n > 0
}

func (r *CommentRepository) AddLike(userID, commentID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&models.CommentLike{UserID: userID, CommentID: commentID}).Error; err != nil {
			return err
		}
		// Обновляем именно таблицу comments поле likes
		return tx.Exec("UPDATE comments SET likes = likes + 1 WHERE id = ?", commentID).Error
	})
}

func (r *CommentRepository) RemoveLike(userID, commentID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? AND comment_id = ?", userID, commentID).Delete(&models.CommentLike{}).Error; err != nil {
			return err
		}
		return tx.Exec("UPDATE comments SET likes = GREATEST(likes - 1, 0) WHERE id = ?", commentID).Error
	})
}

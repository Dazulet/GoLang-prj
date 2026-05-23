package services

import (
	"github.com/mangalib/comment-service/internal/client"
	"github.com/mangalib/comment-service/internal/models"
	"github.com/mangalib/comment-service/internal/repositories"
	"github.com/mangalib/comment-service/internal/validators"
)

type CommentService struct {
	repo       *repositories.CommentRepository
	authClient *client.AuthClient
}

func NewCommentService(repo *repositories.CommentRepository, authClient *client.AuthClient) *CommentService {
	return &CommentService{repo: repo, authClient: authClient}
}

func (s *CommentService) Create(userID uint, req validators.CreateCommentRequest) (*models.Comment, error) {
	c := models.Comment{
		UserID:   userID,
		MangaID:  req.MangaID,
		ParentID: req.ParentID,
		Body:     req.Body,
	}
	if err := s.repo.Create(&c); err != nil {
		return nil, err
	}

	if info, err := s.authClient.GetUser(userID); err == nil {
		c.Author = &models.UserInfo{ID: info.ID, Username: info.Username, Avatar: info.Avatar}
	}
	return &c, nil
}

// ВАЖНО: добавили параметр currentUserID
func (s *CommentService) ListByManga(mangaID uint, page, limit int, currentUserID uint) ([]models.Comment, int64, error) {
	offset := (page - 1) * limit
	comments, total, err := s.repo.ListByManga(mangaID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	userCache := make(map[uint]*models.UserInfo)

	for i := range comments {
		// 1. Обогащаем данными автора
		comments[i].Author = s.getAuthorWithCache(comments[i].UserID, userCache)

		// 2. Проверяем, лайкнул ли текущий юзер этот коммент
		if currentUserID != 0 {
			comments[i].IsLiked = s.repo.LikeExists(currentUserID, comments[i].ID)
		}

		// Обрабатываем ответы (Replies)
		for j := range comments[i].Replies {
			comments[i].Replies[j].Author = s.getAuthorWithCache(comments[i].Replies[j].UserID, userCache)
			if currentUserID != 0 {
				comments[i].Replies[j].IsLiked = s.repo.LikeExists(currentUserID, comments[i].Replies[j].ID)
			}
		}
	}

	return comments, total, nil
}

// Вспомогательный метод для кеширования авторов (чтобы не спамить Auth-сервис)
func (s *CommentService) getAuthorWithCache(userID uint, cache map[uint]*models.UserInfo) *models.UserInfo {
	if info, ok := cache[userID]; ok {
		return info
	}
	info, err := s.authClient.GetUser(userID)
	if err != nil {
		return nil
	}
	mapped := &models.UserInfo{ID: info.ID, Username: info.Username, Avatar: info.Avatar}
	cache[userID] = mapped
	return mapped
}

func (s *CommentService) ToggleLike(userID, commentID uint) (bool, error) {
	if s.repo.LikeExists(userID, commentID) {
		return false, s.repo.RemoveLike(userID, commentID)
	}
	return true, s.repo.AddLike(userID, commentID)
}

func (s *CommentService) Delete(userID, commentID uint, isAdmin bool) error {
	return s.repo.Delete(commentID)
}

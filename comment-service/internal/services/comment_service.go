package services

import (
	"errors"
	"log"

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
		c.Author = &models.UserInfo{
			ID:       info.ID,
			Username: info.Username,
			Avatar:   info.Avatar,
		}
	} else {
		log.Printf("warn: could not fetch author for comment %d: %v", c.ID, err)
	}

	return &c, nil
}

func (s *CommentService) ListByManga(mangaID uint, page, limit int) ([]models.Comment, int64, error) {
	offset := (page - 1) * limit
	comments, total, err := s.repo.ListByManga(mangaID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	userCache := make(map[uint]*models.UserInfo)

	enrichAuthor := func(userID uint) *models.UserInfo {
		if info, ok := userCache[userID]; ok {
			return info
		}
		info, err := s.authClient.GetUser(userID)
		if err != nil {
			log.Printf("warn: could not fetch user %d: %v", userID, err)
			return nil
		}
		mapped := &models.UserInfo{ID: info.ID, Username: info.Username, Avatar: info.Avatar}
		userCache[userID] = mapped
		return mapped
	}

	for i := range comments {
		comments[i].Author = enrichAuthor(comments[i].UserID)
		for j := range comments[i].Replies {
			comments[i].Replies[j].Author = enrichAuthor(comments[i].Replies[j].UserID)
		}
	}

	return comments, total, nil
}

func (s *CommentService) Delete(userID, commentID uint, isAdmin bool) error {
	c, err := s.repo.FindByID(commentID)
	if err != nil {
		return errors.New("comment not found")
	}
	if !isAdmin && c.UserID != userID {
		return errors.New("forbidden")
	}
	return s.repo.Delete(commentID)
}

func (s *CommentService) ToggleLike(userID, commentID uint) (liked bool, err error) {
	if _, err := s.repo.FindByID(commentID); err != nil {
		return false, errors.New("comment not found")
	}
	if s.repo.LikeExists(userID, commentID) {
		return false, s.repo.RemoveLike(userID, commentID)
	}
	return true, s.repo.AddLike(userID, commentID)
}

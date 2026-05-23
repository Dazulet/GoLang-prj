package services

import (
	"errors"

	"github.com/mangalib/auth-service/internal/config"
	"github.com/mangalib/auth-service/internal/models"
	"github.com/mangalib/auth-service/internal/repositories"
	"github.com/mangalib/auth-service/internal/utils"
	"github.com/mangalib/auth-service/internal/validators"
)

type AuthService struct {
	users *repositories.UserRepository
	cfg   *config.Config
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.users.FindAll()
}
func NewAuthService(users *repositories.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{users: users, cfg: cfg}
}

type AuthResult struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func (s *AuthService) Register(req validators.RegisterRequest) (*AuthResult, error) {
	if s.users.ExistsByEmail(req.Email) {
		return nil, errors.New("email already in use")
	}
	if s.users.ExistsByUsername(req.Username) {
		return nil, errors.New("username already taken")
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hash,
		Role:         models.RoleUser,
	}
	if err := s.users.Create(&user); err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(user.ID, string(user.Role), s.cfg.JWTSecret, s.cfg.JWTExpiryHours)
	if err != nil {
		return nil, err
	}
	return &AuthResult{Token: token, User: user}, nil
}

func (s *AuthService) Login(req validators.LoginRequest) (*AuthResult, error) {
	user, err := s.users.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if !utils.CheckPassword(user.PasswordHash, req.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, string(user.Role), s.cfg.JWTSecret, s.cfg.JWTExpiryHours)
	if err != nil {
		return nil, err
	}
	return &AuthResult{Token: token, User: *user}, nil
}

func (s *AuthService) ValidateToken(tokenStr string) (*utils.Claims, error) {
	return utils.ParseToken(tokenStr, s.cfg.JWTSecret)
}

type UserService struct {
	users *repositories.UserRepository
	cfg   *config.Config
}

func NewUserService(users *repositories.UserRepository, cfg *config.Config) *UserService {
	return &UserService{users: users, cfg: cfg}
}

func (s *UserService) GetProfile(id uint) (*models.User, error) {
	u, err := s.users.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return u, nil
}
func (s *UserService) Update(u *models.User) error {
	return s.users.Update(u)
}

func (s *UserService) UpdateBio(id uint, req validators.UpdateProfileRequest) (*models.User, error) {
	u, err := s.users.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	u.Bio = req.Bio
	return u, s.users.Update(u)
}

func (s *UserService) UpdateAvatar(id uint, avatarURL string) error {
	u, err := s.users.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}
	u.Avatar = avatarURL
	return s.users.Update(u)
}

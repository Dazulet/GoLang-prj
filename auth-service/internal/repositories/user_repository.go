package repositories

import (
	"github.com/mangalib/auth-service/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func (r *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var u models.User
	return &u, r.db.First(&u, id).Error
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var u models.User
	return &u, r.db.Where("email = ?", email).First(&u).Error
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) ExistsByEmail(email string) bool {
	var n int64
	r.db.Model(&models.User{}).Where("email = ?", email).Count(&n)
	return n > 0
}

func (r *UserRepository) ExistsByUsername(username string) bool {
	var n int64
	r.db.Model(&models.User{}).Where("username = ?", username).Count(&n)
	return n > 0
}

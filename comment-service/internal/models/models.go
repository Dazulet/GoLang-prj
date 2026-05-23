package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint           `gorm:"primaryKey"              json:"id"`
	UserID    uint           `gorm:"not null;index"          json:"user_id"`
	MangaID   uint           `gorm:"not null;index"          json:"manga_id"`
	ParentID  *uint          `gorm:"index"                   json:"parent_id,omitempty"`
	Body      string         `gorm:"type:text;not null"      json:"body"`
	Likes     int            `gorm:"default:0"               json:"likes"`
	CreatedAt time.Time      `                               json:"created_at"`
	UpdatedAt time.Time      `                               json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                   json:"-"`

	Replies []Comment `gorm:"foreignKey:ParentID"  json:"replies,omitempty"`
	Author  *UserInfo `gorm:"-"                    json:"author,omitempty"`
}

type CommentLike struct {
	UserID    uint      `gorm:"primaryKey"   json:"user_id"`
	CommentID uint      `gorm:"primaryKey"   json:"comment_id"`
	CreatedAt time.Time `                    json:"created_at"`
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar,omitempty"`
}

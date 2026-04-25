package models

type Manga struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Chapters    []Chapter `json:"chapters" gorm:"foreignKey:MangaID"`
}

type Chapter struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	MangaID uint   `json:"manga_id"`
	Title   string `json:"title"`
	Number  int    `json:"number"`
}

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `json:"username" gorm:"column:username"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
}

type Bookmark struct {
	UserID  uint `gorm:"primaryKey" json:"user_id"`
	MangaID uint `gorm:"primaryKey" json:"manga_id"`
}

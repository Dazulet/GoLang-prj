package main

type Manga struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Price       float64   `json:"price" gorm:"default:0"`
	Chapters    []Chapter `json:"chapters" gorm:"foreignKey:MangaID"`
}

type Chapter struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	MangaID uint   `json:"manga_id"`
	Title   string `json:"title"`
	Number  int    `json:"number"`
}

type Genre struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

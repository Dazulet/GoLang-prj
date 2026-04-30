package main

type Order struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	UserID  uint   `json:"user_id"`
	MangaID uint   `json:"manga_id"`
	Status  string `json:"status"`
}

func (Order) TableName() string {
	return "orders"
}

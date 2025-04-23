package models

type Favorite struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	UserID    uint    `json:"user_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
}

package models

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique;not null" json:"username"` // Уникальное и обязательное поле для имени пользователя
	Password string `json:"password"`                        // Пароль
	Email    string `json:"email"`                           // Добавлено поле для email (если требуется)
}

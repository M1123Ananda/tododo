package model

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `json:"name,omitempty"`
	Email    string `gorm:"uniqueIndex, primaryKey" json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type ToDo struct {
	ID          uint   `gorm:"primaryKey" json:"id,omitempty"`
	UserEmail   string `json:"user_email,omitempty"`
	User        User   `gorm:"foreignKey:UserEmail;references:Email" json:"user,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

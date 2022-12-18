package entity

type User struct {
	UserID   uint   `gorm:"primaryKey"`
	Login    string `gorm:"unique"`
	Password string `gorm:"column:password"`
}

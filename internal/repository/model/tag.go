package model

import "time"

type Tag struct {
	UserID  uint   `gorm:"references:UserID"`
	TagID   string `gorm:"primaryKey"`
	TagName string `gorm:"column:tag_name"`
}

type Note struct {
	UserId uint      `gorm:"references:UserID"`
	TagID  string    `gorm:"references:TagID"`
	Note   string    `json:"note" gorm:"column:note"`
	Time   time.Time `gorm:"column:time"`
}

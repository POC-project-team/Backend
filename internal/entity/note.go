package entity

import (
	"time"
)

type Note struct {
	UserId uint      `gorm:"references:UserID"`
	TagID  string    `gorm:"references:TagID"`
	Note   string    `json:"note" gorm:"column:note"`
	Time   time.Time `json:"time" gorm:"column:time"`
}

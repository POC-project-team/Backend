package entity

import (
	"time"
)

type Note struct {
	UserId uint      `json:"-" gorm:"references:UserID"`
	TagID  string    `json:"-" gorm:"references:TagID"`
	Note   string    `json:"note" gorm:"column:note"`
	Time   time.Time `json:"time" gorm:"column:time"`
}

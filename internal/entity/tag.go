package entity

type Tag struct {
	UserID  uint   `json:"-" gorm:"references:UserID"`
	TagID   string `json:"tagID" gorm:"primaryKey"`
	TagName string `json:"tagName" gorm:"column:tag_name"`
}

package user

import (
	"time"
)

type Note struct {
	Note string    `json:"note" gorm:"column:note"`
	Time time.Time `json:"time" gorm:"column:time"`
}

// NewNote constructor for struct note
func NewNote(note string) Note {
	return Note{note, time.Now()}
}

func (n *Note) ToString() string {
	var result string
	result += "Time: " + n.Time.String() + "\nNote: " + n.Note
	return result
}
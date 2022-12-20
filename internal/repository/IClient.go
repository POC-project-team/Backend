package repository

import (
	"backend/internal/entity"
	"gorm.io/gorm"
)

type IClient interface {
	Connect() (*gorm.DB, error)
	Close() error

	GetUser(login, password string) (entity.User, error)
	GetAllUsers() ([]entity.User, error)
	CreateUser(login, password string) (entity.User, error)
	ChangePassword(userId uint, password string) error
	ChangeLogin(userId uint, login string) error
	DeleteUser(userId uint) error
	GetUserByLogin(login string) (entity.User, error)

	GetUserTags(userId uint) ([]entity.Tag, error)
	CreateTag(userId uint, tag, tagName string) (entity.Tag, error)
	DeleteTag(userId uint, tag string) error
	GetTag(userId uint, tagId string) (entity.Tag, error)
	UpdateTag(userId uint, tagId, tagName string) (entity.Tag, error)
	TransferTag(userId uint, tagId, login string) error

	GetUserNotes(userId uint, tagId string) ([]entity.Note, error)
	AddNote(userId uint, tagId, noteInfo string) (entity.Tag, error)
}

package userService

import (
	"backend/domain"
	"backend/internal/entity"
	"backend/internal/repository"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	db repository.IClient
}

func NewService(database repository.IClient) *Service {
	return &Service{
		db: database,
	}
}

// GetAllUsers func to return all users in the map
func (s *Service) GetAllUsers() ([]entity.User, error) {
	users, err := s.db.GetAllUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}

// CreateUser handler for creating new user
func (s *Service) CreateUser(newUser entity.User) (entity.User, error) {
	if _, err := s.db.GetUserByLogin(newUser.Login); err == nil {
		return entity.User{}, fmt.Errorf(domain.UserAlreadyExists)
	}

	user, err := s.db.CreateUser(newUser.Login, newUser.Password)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

// GetAllUsersTags handler for getting all tags of specific user
func (s *Service) GetAllUsersTags(userId uint) ([]entity.Tag, error) {
	tags, err := s.db.GetUserTags(userId)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (s *Service) GetTag(userId uint, tagId string) (entity.Tag, error) {
	tag, err := s.db.GetTag(userId, tagId)
	if err != nil {
		return entity.Tag{}, err
	}

	return tag, nil
}

func (s *Service) CreateTag(newTag entity.Tag) (entity.Tag, error) {
	tag, err := s.db.CreateTag(newTag.UserID, newTag.TagID, newTag.TagName)
	if err != nil {
		return entity.Tag{}, err
	}
	return tag, nil
}

func (s *Service) UpdateTag(newTag entity.Tag) (entity.Tag, error) {
	tag, err := s.db.UpdateTag(newTag.UserID, newTag.TagID, newTag.TagName)
	if err != nil {
		return entity.Tag{}, err
	}

	log.Info("Tag updated successfully for user: ", tag.UserID)
	return tag, nil
}

func (s *Service) DeleteTag(tag entity.Tag) error {
	err := s.db.DeleteTag(tag.UserID, tag.TagID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) TransferTag(tag entity.Tag, login string) error {
	err := s.db.TransferTag(tag.UserID, tag.TagID, login)
	if err != nil {
		return err
	}
	return nil
}

// GetNotes handler for getting notes for specific tag of user
func (s *Service) GetNotes(userId uint, tagId string) ([]entity.Note, error) {
	notes, err := s.db.GetUserNotes(userId, tagId)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

// AddNote handler for creating new note for specific tag of user
func (s *Service) AddNote(note entity.Note) (entity.Tag, error) {
	response, err := s.db.AddNote(note.UserId, note.TagID, note.Note)
	if err != nil {
		return entity.Tag{}, err
	}

	return response, nil
}

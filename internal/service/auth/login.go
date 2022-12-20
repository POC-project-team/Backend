package auth

import (
	"backend/internal/dto/request"
	"backend/internal/entity"
	"backend/internal/repository"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	db repository.IClient
}

func NewAuthService(db repository.IClient) *Service {
	return &Service{db: db}
}

// Auth generates the token for the user
func (s *Service) Auth(newUser entity.User) (request.Token, error) {
	user, err := s.db.GetUser(newUser.Login, newUser.Password)

	if err != nil {
		return request.Token{}, err
	}

	answer, err := request.GenerateToken(user.UserID)
	if err != nil {
		return request.Token{}, err
	}

	log.Info("New token was created for user ", user)

	return answer, nil
}

func (s *Service) ChangeLogin(user entity.User) error {
	if err := s.db.ChangeLogin(user.UserID, user.Login); err != nil {
		return err
	}

	log.Info("Login was changed for user ", user.UserID)

	return nil
}

func (s *Service) ChangePassword(user entity.User) error {
	if err := s.db.ChangePassword(user.UserID, user.Password); err != nil {
		return err
	}

	log.Info("Password was changed for user ", user.UserID)
	return nil
}

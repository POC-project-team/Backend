package auth

import (
	"backend/internal/controller/rest/APIerror"
	"backend/internal/dto/request"
	"backend/internal/dto/request/userRequest"
	"backend/internal/repository/postgres"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Service struct {
	db *postgres.Client
}

func NewAuthService(db *postgres.Client) *Service {
	return &Service{db: db}
}

// Auth generates the token for the user
func (s *Service) Auth(w http.ResponseWriter, r *http.Request) {
	var response userRequest.AuthRequest

	if err := response.Bind(r); err != nil {
		APIerror.Error(w, err)
		return
	}

	UserID, err := s.db.GetUserID(response.Login, response.Password)

	if err != nil {
		APIerror.Error(w, err)
		return
	}

	answer, err := request.GenerateToken(UserID)
	if err != nil {
		APIerror.Error(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(answer); err != nil {
		APIerror.Error(w, err)
	} else {
		log.Info("New token was created for user ", UserID)
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Service) ChangeLogin(w http.ResponseWriter, r *http.Request) {
	var response userRequest.ChangeLoginRequest

	if err := response.Bind(r); err != nil {
		APIerror.Error(w, err)
		return
	}

	if err := s.db.ChangeLogin(response.UserID, response.Login); err != nil {
		APIerror.Error(w, err)
		return
	}

	log.Info("Login was changed for user ", response.UserID)
	w.WriteHeader(http.StatusOK)
}

func (s *Service) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var response userRequest.ChangePasswdRequest

	if response.Bind(r) != nil {
		return
	}

	if err := s.db.ChangePassword(response.UserID, response.Password); err != nil {
		APIerror.Error(w, err)
		return
	}

	log.Info("Password was changed for user ", response.UserID)
	w.WriteHeader(http.StatusCreated)
}

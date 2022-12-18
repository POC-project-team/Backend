package auth

import (
	"backend/internal/controller/rest/APIerror"
	service "backend/internal/dto/request"
	"backend/internal/dto/request/userRequest"
	"backend/internal/repository/postgres"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Service struct {
	db *postgres.Client
}

func NewAuthService(db *postgres.Client) *Service {
	return &Service{db: db}
}

var jwtKey = []byte("secret_key__DO_NOT_POST_IT_TO_GITHUB")

type Claims struct {
	UserId int `json:"userID"`
	jwt.StandardClaims
}

type Token struct {
	JWTToken string `json:"token"`
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

	claims := Claims{
		UserId: UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	var answer Token
	answer.JWTToken, err = token.SignedString(jwtKey)
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
	var response service.Request

	if response.ParseToken(w, r) != nil {
		return
	}

	if err := response.Bind(w, r); err != nil {
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
	var response service.Request

	if response.ParseToken(w, r) != nil {
		return
	}

	if err := response.Bind(w, r); err != nil {
		APIerror.Error(w, err)
		return
	}

	if err := s.db.ChangePassword(response.UserID, response.Password); err != nil {
		APIerror.Error(w, err)
		return
	}

	log.Info("Password was changed for user ", response.UserID)
	w.WriteHeader(http.StatusCreated)
}

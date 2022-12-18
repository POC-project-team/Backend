package userRequest

import (
	"backend/internal/dto/request"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/mux"
	"net/http"
)

type ChangePasswdRequest struct {
	Password string `json:"password"`
	UserID   int
	Token    string
}

func (cpr *ChangePasswdRequest) Bind(r *http.Request) error {
	if err := cpr.Validate(); err != nil {
		return err
	}

	tokenString := mux.Vars(r)["token"]
	if tokenString == "" {
		return fmt.Errorf("token is empty")
	}

	token, err := request.ParseToken(tokenString)
	if err != nil {
		return err
	}
	cpr.UserID = token.UserId
	cpr.Token = tokenString
	return nil
}

func (cpr *ChangePasswdRequest) Validate() error {
	return validation.ValidateStruct(cpr,
		validation.Field(&cpr.Password, validation.Required),
	)
}

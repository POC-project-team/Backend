package userRequest

import (
	"backend/internal/dto/request"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/mux"
	"net/http"
)

type ChangeLoginRequest struct {
	Login  string `json:"login"`
	UserID int
	Token  string
}

func (clr *ChangeLoginRequest) Bind(r *http.Request) error {
	if err := clr.Validate(); err != nil {
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

	clr.UserID = token.UserId
	clr.Token = tokenString

	return nil
}

func (clr *ChangeLoginRequest) Validate() error {
	return validation.ValidateStruct(clr,
		validation.Field(&clr.Login, validation.Required),
	)
}

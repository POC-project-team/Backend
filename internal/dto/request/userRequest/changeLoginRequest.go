package userRequest

import (
	"backend/internal/dto/request"
	validation "github.com/go-ozzo/ozzo-validation/v4"
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
	token, err := request.ParseToken(r)
	if err != nil {
		return err
	}

	clr.UserID = token.UserId

	return nil
}

func (clr *ChangeLoginRequest) Validate() error {
	return validation.ValidateStruct(clr,
		validation.Field(&clr.Login, validation.Required),
	)
}

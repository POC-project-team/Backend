package userRequest

import (
	"backend/internal/dto/request"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
)

type ChangeLoginRequest struct {
	Login string `json:"login"`
	request.TokenRequest
}

func (clr *ChangeLoginRequest) Bind(r *http.Request) error {
	if err := clr.Validate(); err != nil {
		return err
	}

	if err := clr.ParseToken(r); err != nil {
		return err
	}

	return nil
}

func (clr *ChangeLoginRequest) Validate() error {
	return validation.ValidateStruct(clr,
		validation.Field(&clr.Login, validation.Required),
	)
}

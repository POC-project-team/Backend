package userRequest

import (
	"backend/internal/dto/request"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"io/ioutil"
	"net/http"
)

type ChangePasswdRequest struct {
	Password string `json:"password"`
	UserID   int
}

func (cpr *ChangePasswdRequest) Bind(r *http.Request) error {
	//goland:noinspection ALL
	buff, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buff, cpr)
	if err != nil {
		return err
	}

	if err := cpr.Validate(); err != nil {
		return err
	}

	token, err := request.ParseToken(r)
	if err != nil {
		return err
	}
	cpr.UserID = token.UserId
	return nil
}

func (cpr *ChangePasswdRequest) Validate() error {
	return validation.ValidateStruct(cpr,
		validation.Field(&cpr.Password, validation.Required),
	)
}

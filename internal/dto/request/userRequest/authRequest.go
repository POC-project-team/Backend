package userRequest

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"io/ioutil"
	"net/http"
)

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (a *AuthRequest) Bind(r *http.Request) error {
	buff, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buff, a)
	if err != nil {
		return err
	}
	return a.Validate()
}

func (a *AuthRequest) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Login, validation.Required),
		validation.Field(&a.Password, validation.Required),
	)
}

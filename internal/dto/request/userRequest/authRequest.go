package userRequest

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"io/ioutil"
	"net/http"
)

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func toHash(passwd string) string {
	h := sha1.New()
	return hex.EncodeToString(h.Sum([]byte(passwd)))
}

func (a *AuthRequest) Bind(r *http.Request) error {
	//goland:noinspection ALL
	buff, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buff, a)
	if err != nil {
		return err
	}

	err = a.Validate()

	if err != nil {
		return err
	}

	a.Password = toHash(a.Password)

	return nil
}

func (a *AuthRequest) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Login, validation.Required),
		validation.Field(&a.Password, validation.Required),
	)
}

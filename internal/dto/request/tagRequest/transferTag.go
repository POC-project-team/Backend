package tagRequest

import (
	"backend/internal/dto/request"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type TransferTagRequest struct {
	TagID string `json:"tagID"`
	Login string `json:"login"`
	Token request.Claims
}

func (t *TransferTagRequest) Bind(r *http.Request) error {
	buff, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buff, t)
	if err != nil {
		return err
	}

	t.TagID = mux.Vars(r)["tag_id"]

	if err := t.Validate(); err != nil {
		return err
	}
	token, err := request.ParseToken(r)
	if err != nil {
		return err
	}
	t.Token = token
	return nil
}

func (t *TransferTagRequest) Validate() error {
	return validation.ValidateStruct(t,
		validation.Field(&t.Login, validation.Required),
		validation.Field(&t.TagID, validation.Required),
	)
}

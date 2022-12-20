package tagRequest

import (
	"backend/internal/dto/request"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"io/ioutil"
	"net/http"
)

type TransferTagRequest struct {
	Login string `json:"login"`
	request.BasicRequest
}

func (t *TransferTagRequest) Bind(r *http.Request) error {
	//goland:noinspection ALL
	buff, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buff, t)
	if err != nil {
		return err
	}

	if err := t.BindBasicRequest(r); err != nil {
		return err
	}

	if err := t.Validate(); err != nil {
		return err
	}

	return nil
}

func (t *TransferTagRequest) Validate() error {
	return validation.ValidateStruct(t,
		validation.Field(&t.Login, validation.Required),
		validation.Field(&t.TagId, validation.Required),
	)
}

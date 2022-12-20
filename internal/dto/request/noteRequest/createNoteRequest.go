package noteRequest

import (
	"backend/internal/dto/request"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"io/ioutil"
	"net/http"
)

type CreateNoteRequest struct {
	request.TagRequest
	request.TokenRequest
	Note string `json:"note"`
}

func (c *CreateNoteRequest) Bind(r *http.Request) error {
	//goland:noinspection ALL
	buff, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buff, c)
	if err != nil {
		return err
	}
	if err := c.ParseTagId(r); err != nil {
		return err
	}
	if err = c.ParseToken(r); err != nil {
		return err
	}
	if err := c.Validate(); err != nil {
		return err
	}
	return nil
}

func (c *CreateNoteRequest) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Note, validation.Required),
		validation.Field(&c.TagId, validation.Required),
	)
}

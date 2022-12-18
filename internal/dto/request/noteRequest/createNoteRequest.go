package noteRequest

import (
	"backend/internal/dto/request"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type CreateNoteRequest struct {
	Token request.Claims
	TagID string
	Note  string `json:"note"`
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
	c.TagID = mux.Vars(r)["tag_id"]
	if err := c.Validate(); err != nil {
		return err
	}
	token, err := request.ParseToken(r)
	if err != nil {
		return err
	}
	c.Token = token
	return nil
}

func (c *CreateNoteRequest) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Note, validation.Required),
		validation.Field(&c.TagID, validation.Required),
	)
}

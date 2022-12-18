package tagRequest

import (
	"backend/internal/dto/request"
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type CreateUpdateTagRequest struct {
	TagID   string `json:"tag_id"`
	TagName string `json:"tagName"`
	Token   request.Claims
}

func (c *CreateUpdateTagRequest) Bind(r *http.Request) error {
	//golang:noinspection ALL
	buff, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buff, c)
	if err != nil {
		return err
	}
	if err := c.Validate(); err != nil {
		return err
	}

	tagId := mux.Vars(r)["tag_id"]
	if tagId == "" {
		return fmt.Errorf("tagId is required")
	}
	c.TagID = tagId

	token, err := request.ParseToken(r)
	if err != nil {
		return err
	}
	c.Token = token

	return nil
}

func (c *CreateUpdateTagRequest) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.TagName, validation.Required),
	)
}

package tagRequest

import (
	"backend/internal/dto/request"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"io/ioutil"
	"net/http"
)

type CreateUpdateTagRequest struct {
	TagName string `json:"tagName"`
	request.BasicRequest
}

func (c *CreateUpdateTagRequest) Bind(r *http.Request) error {
	//goland:noinspection ALL
	buff, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buff, c)
	if err != nil {
		return err
	}

	if err := c.BindBasicRequest(r); err != nil {
		return err
	}

	if err := c.Validate(); err != nil {
		return err
	}

	return nil
}

func (c *CreateUpdateTagRequest) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.TagName, validation.Required),
	)
}

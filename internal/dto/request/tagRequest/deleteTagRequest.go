package tagRequest

import (
	"backend/internal/dto/request"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type DeleteTagRequest struct {
	TagID string `json:"tagID"`
	Token request.Claims
}

func (d *DeleteTagRequest) Bind(r *http.Request) error {
	tagId := mux.Vars(r)["tag_id"]
	if tagId == "" {
		return fmt.Errorf("tagId is required")
	}

	token, err := request.ParseToken(r)
	if err != nil {
		return err
	}
	d.TagID = tagId
	d.Token = token
	return nil
}

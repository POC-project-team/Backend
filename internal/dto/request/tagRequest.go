package request

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type TagRequest struct {
	TagId string
}

func (t *TagRequest) ParseTagId(r *http.Request) error {
	tagId := mux.Vars(r)["tag_id"]
	if tagId == "" {
		return fmt.Errorf("tagId is required")
	}
	t.TagId = tagId
	return nil
}

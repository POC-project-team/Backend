package tagRequest

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type GetTagRequest struct {
	TagID string `json:"tagID"`
}

func (g *GetTagRequest) Bind(r *http.Request) error {
	tagID := mux.Vars(r)["tag_id"]
	if tagID == "" {
		return fmt.Errorf("tag_id is required")
	}
	g.TagID = tagID
	return nil
}

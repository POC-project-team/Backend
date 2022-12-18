package noteRequest

import (
	"backend/internal/dto/request"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/mux"
	"net/http"
)

type GetNoteRequest struct {
	Token request.Claims
	TagID string
}

func (g *GetNoteRequest) Bind(r *http.Request) error {
	g.TagID = mux.Vars(r)["tag_id"]
	if err := g.Validate(); err != nil {
		return err
	}

	token, err := request.ParseToken(r)
	if err != nil {
		return err
	}
	g.Token = token
	return nil
}

func (g *GetNoteRequest) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.TagID, validation.Required),
	)
}

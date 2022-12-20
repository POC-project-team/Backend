package request

import "net/http"

type BasicRequest struct {
	TokenRequest
	TagRequest
}

func (br *BasicRequest) BindBasicRequest(r *http.Request) error {
	if err := br.ParseToken(r); err != nil {
		return err
	}

	if err := br.ParseTagId(r); err != nil {
		return err
	}

	return nil
}

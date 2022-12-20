package request

import (
	"net/http"
)

type TokenRequest struct {
	Token Claims
}

func (d *TokenRequest) ParseToken(r *http.Request) error {
	token, err := ParseToken(r)
	if err != nil {
		return err
	}
	d.Token = token
	return nil
}

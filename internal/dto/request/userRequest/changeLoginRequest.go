package userRequest

import "net/http"

type ChangeLoginRequest struct {
	Login string `json:"login"`
	Token string
}

func (clr *ChangeLoginRequest) Bind(r *http.Request) {

}

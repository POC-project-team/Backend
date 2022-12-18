package userRequest

type ChangeLoginRequest struct {
	Login string `json:"login"`
	Token string
}

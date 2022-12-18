package userRequest

type ChangePasswdRequest struct {
	Password string `json:"password"`
	Token    string
}

package tagRequest

type CreateUpdateTagRequest struct {
	TagID   string `json:"tagID"`
	TagName string `json:"tagName"`
	Token   string
}

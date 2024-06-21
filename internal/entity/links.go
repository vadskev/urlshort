package entity

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrSlugExists = errors.New("slug already exists")
)

//easyjson:json
type Links struct {
	RawURL string `json:"url"`
	Slug   string `json:"result,omitempty"`
}

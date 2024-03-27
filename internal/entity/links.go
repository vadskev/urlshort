package entity

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrSlugExists = errors.New("slug already exists")
)

// Links
type Links struct {
	RawURL string
	Slug   string
}

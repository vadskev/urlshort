package storage

import "context"

type Storage interface {
	GetURL(alias string) (URLData, error)
	SaveURL(data URLData) error
	Ping(ctx context.Context) error
}

type URLData struct {
	Alias  string `json:"uuid"`
	ResURL string `json:"short_url"`
	URL    string `json:"original_url"`
}

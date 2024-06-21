package storage

import "context"

type Storage interface {
	GetURL(ctx context.Context, alias string) (URLData, error)
	SaveURL(ctx context.Context, data URLData) error
	Ping(ctx context.Context) error
	SaveBatchURL(ctx context.Context, data []URLData) error
	GetURLbyURL(ctx context.Context, url string) (URLData, bool)
}

type URLData struct {
	Alias  string `json:"uuid,omitempty"`
	ResURL string `json:"short_url"`
	URL    string `json:"original_url"`
}

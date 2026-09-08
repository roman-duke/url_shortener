package shortener

import "errors"

var ErrNotFound = errors.New("short link not found")

type Store interface {
	Save(link Link) error
	Get(code string) (Link, error)
	NextID() (int, error)
	GetLinkFromUrl(longUrl string) (Link, error)
}

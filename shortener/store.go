package shortener

import "errors"

var ErrNotFound = errors.New("short link not found")

type Store interface {
	Save(link ShortLink) error
	Get(code string) (ShortLink, error)
	NextID() (int, error)
	GetLinkFromUrl(longUrl string) (ShortLink, error)
}

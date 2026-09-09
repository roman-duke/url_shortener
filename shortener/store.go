package shortener

import "errors"

var ErrNotFound = errors.New("short link not found")

// Changing the name here because I have been conflating
// the interface defintiion with the implementation just
// because they share very similar terminologies
type Persistence interface {
	Save(link Link) error
	Get(code string) (Link, error)
	NextID() (int, error)
	GetLinkFromUrl(longUrl string) (Link, error)
}

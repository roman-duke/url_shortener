package memory

import (
	"slices"
	"github.com/roman-duke/url_shortener/shortener"
)

// we just keep things simple for now
// we use a struct whose member is a slice
// of shortlinks. or we could literally just
// use a slice directly instead of embedding that
// in the struct? The only benefit this would have
// is the ability to make multiple "fake stores" and
// I guess slap on meta-data (other struct members) on
// each.
type MyStore struct {
	ShortLinks []shortener.Link
}

type Link = shortener.Link

// now we define the methods necessary for this to satisfy the
// store interface
func (fS *MyStore) Save(l Link) error {
	// we simply just append to the end of the array
	// the problem with this is that it's hard to fake
	// a storage failure
	fS.ShortLinks = append(fS.ShortLinks, l)
	return nil
}

func (fS *MyStore) Get(c string) (Link, error) {
	// same thing here, we can't really simulate a storage get failure
	// so the only cases we have here are link found and link not found
	linkIdx := slices.IndexFunc(fS.ShortLinks, func(val Link) bool {
		return val.Code == c
	})

	if linkIdx == -1 {
		return Link{}, shortener.ErrNotFound
	}

	return fS.ShortLinks[linkIdx], nil
}

func (fS *MyStore) NextID() (int, error) {
	// the zero-indexing being used to our advantage here
	return len(fS.ShortLinks), nil
}

func (fS *MyStore) GetLinkFromUrl(url string) (Link, error) {
	linkIdx := slices.IndexFunc(fS.ShortLinks, func(val Link) bool {
		return val.LongUrl == url
	})

	if linkIdx == -1 {
		return Link{}, shortener.ErrNotFound
	}

	return fS.ShortLinks[linkIdx], nil
}

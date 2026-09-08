package shortener

import (
	"errors"
	// "fmt"
	"net/url"
)

type Service struct {
	store Store
}

func encode(val int) string {
	// we have to append from the back so we use a buffer that we fill
	// as we gather the remainders during the conversion process

	buf := make([]byte, 11)

	for i := len(buf) - 1; i >= 0; i-- {
		q, r := val/62, val%62

		switch {
		case r >= 0 && r < 10:
			buf[i] = byte('0' + r)
		case r > 9 && r < 36:
			buf[i] = byte(r + 55)
		case r > 35 && r < 62:
			buf[i] = byte(r + 61)
		}

		if q == 0 {
			// we mutate the "window" of the slice
			buf = buf[i:]
			break
		}
		val = q
	}

	return string(buf)
}

var ErrInvalidUrl = errors.New("unable to parse url")

func (s Service) Shorten(longUrl string) (Link, error) {
	// first, we need to confirm that the longUrl we even receive
	// is actually valid via the url specification.
	parsedUrl, err := url.ParseRequestURI(longUrl)

	// we compress the actual parsing error and our "logical" error
	// into the same thing (we are kind of losing information by doing
	// this though)
	if err != nil ||
		!((parsedUrl.Scheme == "http" || parsedUrl.Scheme == "https") &&
			len(parsedUrl.Host) > 0) {
		return Link{}, ErrInvalidUrl
	}

	// we check if the longUrl already exists in our store, if so
	// then we just resolve it
	sLink, err := s.store.GetLinkFromUrl(longUrl)

	switch {
	case errors.Is(err, ErrNotFound):
		// do nothing, this breaks and continues with normal code execution to
		// go create the Link
	case err != nil:
		return Link{}, errors.New("could not verify if url already exists")
	default:
		return sLink, nil
	}

	// Here we produce an offset and then we call our encode helper function
	// in order to produce the base62 string
	nextId, err := s.store.NextID()
	if err != nil {
		return Link{}, errors.New("could not generate short url")
	}

	// we use this as our offset so that initial codes are not "weird"
	offset := 1250
	c := encode(nextId + offset)

	// build the Link that our domain describes/knows.
	link := Link{
		Code:    c,
		LongUrl: longUrl,
	}

	// We call the method on store to save the short link
	// without caring about the specifics of the "saving"
	saveErr := s.store.Save(link)

	if saveErr != nil {
		return Link{}, errors.New("unexpected error occurred, could not save short url")
	}

	// // Just to indicate that we generate our shortcode successfully
	// fmt.Printf("Generated code ->%s", c)

	return link, nil
}

func (s Service) Resolve(code string) (Link, error) {
	// Check if the code -> longUrl mapping exists
	link, err := s.store.Get(code)

	if err != nil {
		return Link{}, err
	}

	return link, nil
}

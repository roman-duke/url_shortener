package shortener

import (
	"errors"
	"fmt"
	"slices"
	"testing"
)

type EncodingTestCase struct {
	in   int
	want string
}

type ShortenTestCase struct {
	in      string
	want    ShortLink
	wantErr error
}

type ResolveTestCase = ShortenTestCase

// we just keep things simple for now
// we use a struct whose member is a slice
// of shortlinks. or we could literally just
// use a slice directly instead of embedding that
// in the struct? The only benefit this would have
// is the ability to make multiple "fake stores" and
// I guess slap on meta-data (other struct members) on
// each.
type MyFakeStore struct {
	shortLinks []ShortLink
}

// now we define the methods necessary for this to satisfy the
// store interface
func (fS *MyFakeStore) Save(l ShortLink) error {
	// we simply just append to the end of the array
	// the problem with this is that it's hard to fake
	// a storage failure
	fS.shortLinks = append(fS.shortLinks, l)
	return nil
}

func (fS *MyFakeStore) Get(c string) (ShortLink, error) {
	// same thing here, we can't really simulate a storage get failure
	// so the only cases we have here are link found and link not found
	linkIdx := slices.IndexFunc(fS.shortLinks, func(val ShortLink) bool {
		return val.code == c
	})

	if linkIdx == -1 {
		return ShortLink{}, ErrNotFound
	}

	return fS.shortLinks[linkIdx], nil
}

func (fS *MyFakeStore) NextID() (int, error) {
	// the zero-indexing being used to our advantage here
	return len(fS.shortLinks), nil
}

func (fS *MyFakeStore) GetLinkFromUrl(url string) (ShortLink, error) {
	linkIdx := slices.IndexFunc(fS.shortLinks, func(val ShortLink) bool {
		return val.longUrl == url
	})

	if linkIdx == -1 {
		return ShortLink{}, ErrNotFound
	}

	return fS.shortLinks[linkIdx], nil
}

func TestService(t *testing.T) {
	base62ConvCases := []EncodingTestCase{
		{
			10,
			"A",
		},
		{
			0,
			"0",
		},
		{
			61,
			"z",
		},
		{
			62,
			"10",
		},
		{
			27,
			"R",
		},
		{
			83734,
			"LmY",
		},
		{
			543,
			"8l",
		},
		{
			1111,
			"Hv",
		},
		{
			2738578297829239,
			"CXeCoTFbD",
		},
		{
			1985,
			"W1",
		},
		{
			737474823,
			"nuMtD",
		},
	}

	for _, tc := range base62ConvCases {
		t.Run(fmt.Sprintf("convert %d to base62", tc.in), func(t *testing.T) {
			if got := encode(tc.in); got != tc.want {
				t.Errorf("got %+v, want %+v", got, tc.want)
			}
		})
	}

	// Creating a fake store in order to test the real
	// service methods - Shorten and Resolve
	fakeStore := MyFakeStore{
		make([]ShortLink, 0),
	}

	// This was a serious bug (using a value receiver as opposed to a pointer receiver)
	// on each call, a copy of the struct was being made and pass to the method.
	// TODO: Look into value and pointer receivers again.
	newService := Service{&fakeStore}

	// the service shorten method as its own subtest case
	shortenTestCases := []ShortenTestCase{
		{
			"http://youtube.com/eren-jaeger",
			ShortLink{
				"KA",
				"http://youtube.com/eren-jaeger",
			},
			nil,
		},
		{
			"http://",
			ShortLink{},
			ErrInvalidUrl,
		},
		{
			"http://youtube.com/eren-jaeger",
			// If given the same longurl, it should
			// return the same shortlink
			ShortLink{
				"KA",
				"http://youtube.com/eren-jaeger",
			},
			nil,
		},
		{
			"https://google.com",
			ShortLink{
				"KB",
				"https://google.com",
			},
			nil,
		},
	}

	for _, tc := range shortenTestCases {
		t.Run(fmt.Sprintf("shorten the longUrl -> %s", tc.in), func(t *testing.T) {
			got, err := newService.Shorten(tc.in)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("err = %v; want %v", err, tc.wantErr)
			}

			if got.code != tc.want.code || got.longUrl != tc.want.longUrl {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}

	t.Log(fakeStore.shortLinks)

	// the service resolve method as its own subtest case
	resolveTestCases := []ResolveTestCase{
		{
			"KA",
			ShortLink{
				"KA",
				"http://youtube.com/eren-jaeger",
			},
			nil,
		},
		{
			"KB",
			ShortLink{
				"KB",
				"https://google.com",
			},
			nil,
		},
		{
			"KC",
			ShortLink{},
			ErrNotFound,
		},
	}

	for _, tc := range resolveTestCases {
		t.Run(fmt.Sprintf("resolve the given short code -> %v", tc.in), func(t *testing.T) {
			got, err := newService.Resolve(tc.in)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("err = %v, want %v", err, tc.wantErr)
			}

			if got.code != tc.want.code || got.longUrl != tc.want.longUrl {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

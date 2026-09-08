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
	want    Link
	wantErr error
}

type ResolveTestCase = ShortenTestCase

type MyFakeStore struct {
	ShortLinks []Link
}

func (fS *MyFakeStore) Save(l Link) error {
	fS.ShortLinks = append(fS.ShortLinks, l)
	return nil
}

func (fS *MyFakeStore) Get(c string) (Link, error) {
	linkIdx := slices.IndexFunc(fS.ShortLinks, func(val Link) bool {
		return val.Code == c
	})

	if linkIdx == -1 {
		return Link{}, ErrNotFound
	}

	return fS.ShortLinks[linkIdx], nil
}

func (fS *MyFakeStore) NextID() (int, error) {
	return len(fS.ShortLinks), nil
}

func (fS *MyFakeStore) GetLinkFromUrl(url string) (Link, error) {
	linkIdx := slices.IndexFunc(fS.ShortLinks, func(val Link) bool {
		return val.LongUrl == url
	})

	if linkIdx == -1 {
		return Link{}, ErrNotFound
	}

	return fS.ShortLinks[linkIdx], nil
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
	testStore := MyFakeStore{
		make([]Link, 0),
	}

	// This was a serious bug (using a value receiver as opposed to a pointer receiver)
	// on each call, a copy of the struct was being made and pass to the method.
	// TODO: Look into value and pointer receivers again.
	newService := Service{&testStore}

	// the service shorten method as its own subtest case
	shortenTestCases := []ShortenTestCase{
		{
			"http://youtube.com/eren-jaeger",
			Link{
				"KA",
				"http://youtube.com/eren-jaeger",
			},
			nil,
		},
		{
			"http://",
			Link{},
			ErrInvalidUrl,
		},
		{
			"http://youtube.com/eren-jaeger",
			// If given the same longurl, it should
			// return the same shortlink
			Link{
				"KA",
				"http://youtube.com/eren-jaeger",
			},
			nil,
		},
		{
			"https://google.com",
			Link{
				"KB",
				"https://google.com",
			},
			nil,
		},
	}

	for _, tc := range shortenTestCases {
		t.Run(fmt.Sprintf("shorten the LongUrl -> %s", tc.in), func(t *testing.T) {
			got, err := newService.Shorten(tc.in)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("err = %v; want %v", err, tc.wantErr)
			}

			if got.Code != tc.want.Code || got.LongUrl != tc.want.LongUrl {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}

	t.Log(testStore.ShortLinks)

	// the service resolve method as its own subtest case
	resolveTestCases := []ResolveTestCase{
		{
			"KA",
			Link{
				"KA",
				"http://youtube.com/eren-jaeger",
			},
			nil,
		},
		{
			"KB",
			Link{
				"KB",
				"https://google.com",
			},
			nil,
		},
		{
			"KC",
			Link{},
			ErrNotFound,
		},
	}

	for _, tc := range resolveTestCases {
		t.Run(fmt.Sprintf("resolve the given short Code -> %v", tc.in), func(t *testing.T) {
			got, err := newService.Resolve(tc.in)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("err = %v, want %v", err, tc.wantErr)
			}

			if got.Code != tc.want.Code || got.LongUrl != tc.want.LongUrl {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

package shortener

import (
	"fmt"
	"strconv"
	"strings"
)

type Service struct {
	store Store
}

func encode(val int) string {
	var result strings.Builder

	for loop := true; loop; {
		q, r := val/62, val%62

		switch {
		case r >= 0 && r < 10:
			result.WriteString(strconv.Itoa(r))
		case r > 9 && r < 36:
			result.WriteString(string(rune(r + 55)))
		case r > 35 && r < 62:
			result.WriteString(string(rune(r + 61)))
		}

		if q == 0 {
			loop = false
		}
	}

	return result.String()
}

func (s Service) Shorten(longUrl string) (ShortLink, error) {
	// Here we produce an offset and then we call our encode helper function
	// in order to produce the base62 string
	nextId, err := s.store.NextID()
	if err != nil {
	}

	// we use this as our offset so that initial codes are not "weird"
	offset := 1250

	code := encode(nextId + offset)
	fmt.Printf("Generated code ->%s", code)

	return ShortLink{}, nil
}

func (s Service) Resolve(code string) (ShortLink, error) {
	return ShortLink{}, nil
}

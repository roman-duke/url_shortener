package shortener

type ShortenerService struct {}

func (s ShortenerService) Shorten(longUrl string) (ShortLink, error) {
	return ShortLink{}, nil
}

func (s ShortenerService) Resolve(code string) (ShortLink, error) {
	return ShortLink{}, nil
}

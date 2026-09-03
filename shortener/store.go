package shortener

type Store interface {
	Save(link ShortLink) error
	Get(code string) (ShortLink, error)
	NextID() (int, error)
	// Exists(longUrl string)
}

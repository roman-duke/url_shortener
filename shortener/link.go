package shortener

// Domain-level "entity"
// bag of data that describes "what" the app is about (the nouns)
type ShortLink struct {
	code string
	longUrl string

	// these to be added a bit later
	// createdAt string
	// clickCount int
	// expiration string
}

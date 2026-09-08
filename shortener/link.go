package shortener

// Domain-level "entity"
// bag of data that describes "what" the app is about (the nouns)
type Link struct {
	Code string
	LongUrl string

	// these to be added a bit later
	// createdAt string
	// clickCount int
	// expiration string
}

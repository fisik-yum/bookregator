package sources

// Struct for itemized reviews
type ItemizedReview struct {
	Olid       string  `json:"olid" csv:"olid"`
	Source     string  `json:"source" csv:"source"`
	ExternalID string  `json:"external_id" csv:"external_id"`
	Username   string  `json:"username" csv:"username"`
	Rating     float64 `json:"rating" csv:"rating"`
	Text       string  `json:"text" csv:"text"`
}
type Source interface {
	GetReviews(isbn string) ([]ItemizedReview, error)
}

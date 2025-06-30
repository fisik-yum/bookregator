package sources

// Struct for itemized reviews
type ItemizedReview struct {
	Olid       string          `json:"olid"`
	Source     string          `json:"source"`
	ExternalID string          `json:"external_id"`
	Username   string          `json:"username"`
	Rating     float64 `json:"rating"`
	Text       string `json:"text"`
}
type Source interface {
	GetReviews(isbn string) ([]ItemizedReview, error)
}

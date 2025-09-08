package common

type Review struct {
	Olid       string  `json:"olid"`
	Source     string  `json:"source"`
	ExternalID string  `json:"external_id"`
	Username   string  `json:"username"`
	Rating     float64 `json:"rating"`
	Text       string  `json:"text"`
}

type Work struct {
	OLID          string  `json:"olid"`
	Title         string  `json:"title"`
	Author        *string `json:"author"`
	Cover         *string `json:"cover"`
	Description   *string `json:"description"`
	PublishedYear *int64  `json:"published_year"`
}

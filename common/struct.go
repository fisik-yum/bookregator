package common

type Review struct {
	Olid       string  `json:"olid"`
	Source     string  `json:"source"`
	ExternalID string  `json:"external_id"`
	Username   string  `json:"username"`
	Rating     float64 `json:"rating"`
	Text       string  `json:"text"`
}

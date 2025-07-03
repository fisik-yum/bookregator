package hardcover

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"io"
	"net/http"
	"web_pull/sources"
)

const SOURCE_IDENT = "goodreads"

type HardcoverScraper struct{}

/*
type ItemizedReview struct {
	Olid       string          `json:"olid"`
	Source     string          `json:"source"`
	ExternalID string          `json:"external_id"`
	Username   string
	Rating     sql.NullFloat64 `json:"rating"`
	Text       sql.NullString  `json:"text"`
}
*/

func (h HardcoverScraper) GetReviews(isbn string) ([]sources.ItemizedReview, error) {
	var reviews = make([]sources.ItemizedReview, 0)
	reviewCollector := colly.NewCollector(
		colly.AllowedDomains("hardcover.app"),
	)
	reviewCollector.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36"

	return reviews, nil
}

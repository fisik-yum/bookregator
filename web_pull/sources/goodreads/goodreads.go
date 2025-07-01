package goodreads

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"web_pull/sources"

	"github.com/gocolly/colly/v2"
)

const SOURCE_IDENT="goodreads"

type GoodreadsScraper struct{}

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

// NOTE: assign OLID as-is
func (g GoodreadsScraper) GetReviews(isbn string) ([]sources.ItemizedReview, error) {
	var reviews = make([]sources.ItemizedReview,0)
	reviewCollector := colly.NewCollector(
		colly.AllowedDomains("www.goodreads.com"),
	)

	// Example URL pattern: https://www.goodreads.com/book/isbn/9780143127741
	url := fmt.Sprintf("https://www.goodreads.com/book/isbn/%s", isbn)
	reviewCollector.OnHTML(".ReviewCard", func(e *colly.HTMLElement) {
		// Create Review Object
		review := sources.ItemizedReview{}
		// Find Username
		review.Username = e.DOM.Find(".ReviewerProfile__name").Text()
		// Find Review Text
		review.Text = e.DOM.Find(".ReviewCard__content").Find("span.Formatted").Text()
		// Find Rating
		rating_str := e.DOM.Find(".ReviewCard__content").Find("span.RatingStars.RatingStars__small").AttrOr("aria-label", "Rating 0 out of 5")
		val, err:=strconv.ParseFloat(strings.SplitN(rating_str, " ", 3)[1],64)
		if err!=nil{
			log.Fatal(err)
			return
		}
		review.Rating=val
		// Get GR Review ID ?ex https://www.goodreads.com/review/show/
		review_link:= e.DOM.Find(".ReviewCard__content").Find("span.Text.Text__body3").Find("a").AttrOr("href","")
		extid_split:=strings.Split(review_link,"/")
		review.ExternalID=extid_split[len(extid_split)-1]
		// Assign Source ID	
		review.Source=SOURCE_IDENT
		// Append to reviews
		reviews = append(reviews, review)
	})

	err := reviewCollector.Visit(url)
	if err != nil {
		log.Println(err)
	}
	return reviews, nil
}

package htmlbuilder

import (
	"api_front/internal"
	"book.buckminsterfullerene.net/common"
	"encoding/json"
	"io"
	"log"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"strconv"
)

func Review(review common.Review) g.Node {
	return Div(Class("card w-50 mx-auto"),
		Div(Class("card-body"),
			H5(Class("class-title"),
				g.Text(review.Username),
			),
			H6(Class("class-subtitle"),
				g.Text(strconv.FormatFloat(review.Rating, 'f', 0, 64)+" on "+review.Source),
			),
			P(Class("card-text"),
				g.Text(review.Text),
			),
		),
	)
}

func bookInfo(olid string) *common.Review {
	book, err := internal.BookClient.Book.ById(olid + "W")
	if err != nil {
		return nil
	}
	ret:=&common.Review{

	}
}

func getReviews(olid string) []common.Review {
	resp, err := internal.Client.Get("http://127.0.0.1:1024/api/get/work?olid=" + olid)
	if err != nil {
		log.Print(err)
		log.Printf("Review Request for OLID: %s failed", olid)
	}

	reviewJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Review Request for OLID: %s failed", olid)
	}

	var reviewObj []common.Review
	err = json.Unmarshal(reviewJSON, &reviewObj)
	if err != nil {
		log.Printf("Review Request for OLID: %s failed", olid)
		log.Println(err)
	}
	return reviewObj
}

func ReviewPage(olid string) g.Node {
	book := BookInfo(olid)
	reviewNodes := g.Map(getReviews(olid), Review)
	pageContent := []g.Node{book, reviewNodes}
	return Page("", pageContent...)
}

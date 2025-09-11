package htmlbuilder

import (
	"api_front/internal"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"book.buckminsterfullerene.net/db"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Review(review db.Review) g.Node {
	return Div(Class("card w-50 mx-auto"),
		Div(Class("card-body"),
			H4(Class("class-title"),
				g.Text(review.Username),
			),
			H5(Class("class-subtitle"),
				g.Text(strconv.FormatFloat(*review.Rating, 'f', 0, 64)+" on "+review.Source),
			),
			P(Class("card-text"),
				g.Text(*review.Text),
			),
		),
	)
}

//func bookInfo(olid string) *common.Review {
//	book, err := internal.BookClient.Book.ById(olid + "W")
//	if err != nil {
//		return nil
//	}
//	ret:=&common.Review{
//
//	}
//}

func getReviews(olid string) []db.Review {
	reqJSON, err := json.Marshal(db.GetXByOLIDParams{
		OLID: olid})
	if err != nil {
		log.Print(err)
		log.Printf("Review Request for OLID: %s failed", olid)
	}
	req, err := http.NewRequest("GET", "http://127.0.0.1:1024/api/get/reviews", bytes.NewBuffer(reqJSON))
	if err != nil {
		log.Print(err)
		log.Printf("Review Request for OLID: %s failed", olid)
	}

	resp, err := internal.Client.Do(req)
	if err != nil {
		log.Print(err)
		log.Printf("Review Request for OLID: %s failed", olid)
	}

	reviewJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Review Request for OLID: %s failed", olid)
	}

	var reviewObj []db.Review
	err = json.Unmarshal(reviewJSON, &reviewObj)
	if err != nil {
		log.Printf("Review Request for OLID: %s failed", olid)
		log.Println(err)
	}
	return reviewObj
}

func ReviewPage(olid string) g.Node {
	//book := BookInfo(olid)
	reviewNodes := g.Map(getReviews(olid), Review)
	pageContent := []g.Node{reviewNodes}
	return Page("", pageContent...)
}

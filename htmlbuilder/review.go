package htmlbuilder

import (
	"book.buckminsterfullerene.net/db"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"strconv"
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

func ReviewPage(reviews []db.Review) g.Node {
	//book := BookInfo(olid)
	reviewNodes := g.Map(reviews, Review)
	pageContent := []g.Node{reviewNodes}
	return Page("", pageContent...)
}

package htmlbuilder

import (
	"strconv"

	"book.buckminsterfullerene.net/common"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
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

func ReviewPage(reviews []common.Review) g.Node {
	reviewNodes := g.Map(reviews, Review)
	return Page("Test", reviewNodes...)
}

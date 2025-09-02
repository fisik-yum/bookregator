package htmlbuilder

import (
	"book.buckminsterfullerene.net/common"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func makeReview(review common.Review) g.Node {
	return Div(Class("panel panel-default"),
		Div(Class("panel-heading"),
			g.Text(review.Username),
		),
		Div(Class("panel-body"),
			g.Text(review.Text),
	),
	)
}

func ReviewPage(reviews []common.Review) g.Node {
	reviewNodes := g.Map(reviews, makeReview)
	return Page("Test", reviewNodes...)
}

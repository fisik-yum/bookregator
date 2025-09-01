package htmlbuilder

import (
	"book.buckminsterfullerene.net/common"
	g "maragu.dev/gomponents"
	comp "maragu.dev/gomponents/components"
	html "maragu.dev/gomponents/html"
)

func Review(review common.Review) g.Node {
	return html.P(
		html.H3(g.Text(review.Username)),
		g.Text(review.Text))
}

func ReviewPage(reviews []common.Review) g.Node {
	reviewNodes:=make([]g.Node,len(reviews))
	for i,v:=range reviews{
		reviewNodes[i]=Review(v)
	}
	return comp.HTML5(comp.HTML5Props{
		Body: reviewNodes,
	})
}

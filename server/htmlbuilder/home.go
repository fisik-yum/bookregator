package htmlbuilder

import (
	_ "embed"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Home() g.Node {
	return  Page("bookregator: home",Div(Class("jumbotron w-50 mx-auto"),
		H1(Class("display-4"),
			g.Text("bookregator"),
		),
		P(Class("lead"),
			g.Text("the free book review aggregator"),
		),
	),)
}

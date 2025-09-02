package htmlbuilder

import (
	_ "embed"
	g "maragu.dev/gomponents"
	comp "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

//go:embed head.embed
var HTML_IMPORTS_HEAD string

func Navbar() g.Node {
	return Nav(Class("navbar navbar-default"),
		Div(Class("container-fluid"),
			Div(Class("navbar-header"),
				A(Class("navbar-brand"),
					g.Text("bookregator"),
				),
			),
			// NOTE: include navbar stuff *inside* this ul tag
			Div(Class("collapse navbar-collapse"),
				Ul(Class("nav navbar-nav"),
					NavbarLink("Home", "/"),
					NavbarLink("About", "/about"),
					SearchBarOLID(),
				),
			),
		),
	)
}

func NavbarLink(name, path string) g.Node {
	return Li(A(Href(path), g.Text(name)))
}

func SearchBarOLID() g.Node {
	return Form(Class("navbar-form navbar-right"),
		Div(Class("form-group"),
			Input(Type("text"), Class("form-control"), Placeholder("Search OLID")),
		),
		Button(Type("submit"), Class("btn btn-default"),
			g.Text("Submit"),
		),
	)
}

func Disclaimer() g.Node {
	return Footer(g.Text("(c) all non-user generated content on this website is licensed under the CC BY license"))
}

// TODO: abstract some templating stuff over this
func Page(title string, nodes ...g.Node) g.Node {
	head := []g.Node{Navbar()}
	return comp.HTML5(comp.HTML5Props{
		Title: title,
		Head: []g.Node{
			g.Raw(HTML_IMPORTS_HEAD),
		},
		Body: append(head, nodes...),
	})
}

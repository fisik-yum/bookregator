package htmlbuilder

import (
	_ "embed"
	g "maragu.dev/gomponents"
	comp "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

//go:embed head.embed
var HTML_IMPORTS_HEAD string

//go:embed tail.embed
var HTML_IMPORTS_TAIL string

func Navbar() g.Node {
	return Nav(Class("navbar navbar-expand navbar-light bg-light"),
		Div(Class("container-fluid"),
			A(Class("navbar-brand"),
				g.Text("bookregator"),
			),
			// NOTE: include navbar stuff *inside* this ul tag
			Div(Class("collapse navbar-collapse"),
				Ul(Class("navbar-nav me-auto mb-2 mb-lg-0"),
					NavbarLink("Home", "/"),
					NavbarLink("About", "/about"),
					SearchBarOLID(),
				),
			),
		),
	)
}

func NavbarLink(name, path string) g.Node {
	return Li(Class("nav-item"),
		A(Class("nav-link"),
			Href(path),
			g.Text(name),
		),
	)
}

func SearchBarOLID() g.Node {
	return Form(Class("d-flex"),
		Input(Class("form-control me-2"),
			Type("search"),
			Placeholder("Search OLID"),
		),
		Button(Class("btn btn-outline-success"),
			Type("submit"),
			// TODO: add glyph
			g.Text("Search"),
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
		// TODO: figure out better way to do this
		Body: append(head, append(nodes, g.Raw(HTML_IMPORTS_TAIL))...),
	})
}

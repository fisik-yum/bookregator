package htmlbuilder

import (
	 g"maragu.dev/gomponents"
	 html"maragu.dev/gomponents/html"
)

func Navbar() g.Node {
	return html.Nav(html.Class("navbar"),
		html.Ol(
			NavbarItem("Home", "/"),
			NavbarItem("About", "/about"),
		),
	)
}

func NavbarItem(name, path string) g.Node {
	return html.Li(html.A(html.Href(path), g.Text(name)))
}

func Disclaimer() g.Node{
	return html.Footer(g.Text("(c) all non-user generated content on this website is licensed under the CC BY license"))
}


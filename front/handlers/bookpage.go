package handlers

import (
	"api_front/htmlbuilder"
	"net/http"
)

func ByOLID(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	olid := r.Form.Get("olid")
	v := htmlbuilder.ReviewPage(olid)
	v.Render(w)
}

func Home(w http.ResponseWriter, r *http.Request) {
	htmlbuilder.Home().Render(w)
}

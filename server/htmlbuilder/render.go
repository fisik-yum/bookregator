package htmlbuilder

import (
	_ "embed"
	"net/http"
	"server/db"
	"server/htmlbuilder/pages"
	"server/htmlbuilder/shared"
)

func Index(w http.ResponseWriter, r *http.Request) {
	pages.NewIndex(shared.NewBase()).Render(w, r)
}

func ReviewPage(w http.ResponseWriter, r *http.Request, rev []db.Review) {
	pages.NewReview(shared.NewBase(),rev).Render(w, r)
}

package handlers

import (
	"api_front/htmlbuilder"
	"api_front/internal"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"book.buckminsterfullerene.net/common"
)

func ByOLID(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	olid := r.Form.Get("olid")
	resp, err := internal.Client.Get("http://127.0.0.1:1024/api/get/work?olid=" + olid)
	if err != nil {
		log.Print(err)
		log.Printf("Review Request for OLID: %s failed", olid)
	}

	reviewJSON,err:=io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Review Request for OLID: %s failed", olid)
	}

	var reviewObj []common.Review
	err=json.Unmarshal(reviewJSON,&reviewObj)
	if err != nil {
		log.Printf("Review Request for OLID: %s failed", olid)
		log.Println(err)
	}
	v:=htmlbuilder.ReviewPage(reviewObj)
	v.Render(w)
}

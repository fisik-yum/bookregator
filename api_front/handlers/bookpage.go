package handlers

import (
	"api_front/internal"
	"io"
	"log"
	"net/http"
)

func ByOLID(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	olid := r.Form.Get("olid")
	resp, err := internal.Client.Get("http://127.0.0.1:1024/api/get/work?olid=" + olid)
	if err != nil {
		log.Print(err)
		log.Printf("Review Request for OLID: %s failed", olid)
	}
	fval,err:=io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Review Request for OLID: %s failed", olid)
	}
	w.Write(fval)
}

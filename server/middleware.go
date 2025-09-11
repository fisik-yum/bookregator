package main

import (
	"log"
	"net/http"
	"time"
)

type Logger struct {
	handler http.Handler
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %v", r.Method, r.URL.Path, time.Now().Format(time.RFC822))
	l.handler.ServeHTTP(w, r)
}

type APIHandler struct {
	handler http.Handler
}

//func (a *APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	u, p, _ := r.BasicAuth()
//}

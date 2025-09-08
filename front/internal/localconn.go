package internal

import (
	"git.sr.ht/~timharek/openlibrary-go"
	"net/http"
)

var BookClient = openlibrary.New()
var Client = &http.Client{}

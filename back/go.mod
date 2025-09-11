module api_back

go 1.24.6

replace book.buckminsterfullerene.net/db => ../db

require (
	book.buckminsterfullerene.net/db v0.0.0-00010101000000-000000000000
	git.sr.ht/~timharek/openlibrary-go v0.0.0
	github.com/mattn/go-sqlite3 v1.14.32
	github.com/moraes/isbn v0.0.0-20151007102746-e6388fb1bfd5
)

module api_front

go 1.24.6

replace book.buckminsterfullerene.net/db => ../db

require (
	book.buckminsterfullerene.net/db v0.0.0-00010101000000-000000000000
	maragu.dev/gomponents v1.2.0
)

require (
	git.sr.ht/~timharek/openlibrary-go v0.0.0
	github.com/mattn/go-sqlite3 v1.14.32 // indirect
)

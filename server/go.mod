module server

go 1.25.1

replace book.buckminsterfullerene.net/db => ../db

replace book.buckminsterfullerene.net/htmlbuilder => ../htmlbuilder

require (
	book.buckminsterfullerene.net/db v0.0.0-00010101000000-000000000000
	book.buckminsterfullerene.net/htmlbuilder v0.0.0-00010101000000-000000000000
	github.com/moraes/isbn v0.0.0-20151007102746-e6388fb1bfd5
)

require (
	git.sr.ht/~timharek/openlibrary-go v0.0.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.32 // indirect
	maragu.dev/gomponents v1.2.0 // indirect
)

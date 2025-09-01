module api_front

go 1.24.6

replace book.buckminsterfullerene.net/common => ../common

require (
	git.sr.ht/~timharek/openlibrary-go v0.0.0
	github.com/mattn/go-sqlite3 v1.14.28
	github.com/moraes/isbn v0.0.0-20151007102746-e6388fb1bfd5
)

require (
	book.buckminsterfullerene.net/common v0.0.0-00010101000000-000000000000 // indirect
	maragu.dev/gomponents v1.2.0 // indirect
	maragu.dev/gomponents.html v1.2.0 // indirect
)

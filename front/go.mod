module api_front

go 1.24.6

replace book.buckminsterfullerene.net/common => ../common

require (
	book.buckminsterfullerene.net/common v0.0.0-00010101000000-000000000000
	maragu.dev/gomponents v1.2.0
)

require git.sr.ht/~timharek/openlibrary-go v0.0.0 // indirect

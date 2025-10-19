build-debug:
	go build -v && mv server artifacts/
build-test:
	templ generate
	cd db;sqlc generate
	go build -v && mv server artifacts/ 
	cd artifacts;./server -db test.sqlite3 -user test -pass test
generate:
	templ generate
	cd db;sqlc generate
git:
	gofmt -s -w -l .
	go fix
	go mod tidy

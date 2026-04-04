build-debug:
	go build -v && mv server artifacts/
build-test:
	templ generate
	cd db;sqlc generate
	go build -v && mv server artifacts/ 
	cp conf.yml artifacts/conf.yml
	cd artifacts;./server -c conf.yml
generate:
	templ generate
	cd db;sqlc generate
git:
	gofmt -s -w -l .
	go fix
	go mod tidy

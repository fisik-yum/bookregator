build-debug:
	go build -v && mv server artifacts/
build-test:
	templ generate
	cd db;sqlc generate
	go build -v && mv server artifacts/ 
	cd artifacts;./server
generate:
	templ generate
	cd db;sqlc generate

build-debug:
	go build -v && mv server artifacts/
build-test:
	templ generate
	go build -v && mv server artifacts/ 
	cd artifacts;./server

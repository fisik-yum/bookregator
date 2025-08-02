build-debug:
	go build -v && mv api_back artifacts/
prep:
	sqlc generate
	mkdir -p artifacts && cp schema.sql internal/

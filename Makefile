fmt: 
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -o app cmd/blog/main.go 

test: build
	go test ./...

# offer a way to run without DB in case of dependency issue
run-local: build
	ENV=dev ./app

run-integration:
	docker compose up --build --force-recreate
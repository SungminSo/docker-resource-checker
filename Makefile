tidy: go.mod
	@go mod tidy

go.sum: go.mod
	@go mod verify

build: go.sum
	go build -o build/checker cmd/main.go
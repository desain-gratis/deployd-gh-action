build:
	GOOS=linux GOARCH=amd64 go build -o dist/action-linux-amd64 ./cmd/action/main.go

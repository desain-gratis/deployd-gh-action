ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

download-lib:
	curl -LO https://storage.googleapis.com/downloads.webmproject.org/releases/webp/libwebp-1.6.0-linux-x86-64.tar.gz
	tar -xzvf libwebp-1.6.0-linux-x86-64.tar.gz -C dist
	rm libwebp-1.6.0-linux-x86-64.tar.gz

build:
#   https://developers.google.com/speed/webp/docs/using
# 	sudo apt-get install -y musl musl-dev musl-tools
	CC=musl-gcc \
	CGO_ENABLED=1 \
	CGO_CFLAGS="-I$(shell pwd)/dist/libwebp-1.6.0-linux-x86-64/include -I dist/libwebp-1.6.0-linux-x86-64/include/webp" \
	CGO_LDFLAGS="-L$(shell pwd)/dist/libwebp-1.6.0-linux-x86-64/lib -lwebp -lsharpyuv -static" \
	GOOS=linux GOARCH=amd64 go build -ldflags="-extldflags '-static' -s -w" -o dist/action-linux-amd64 ./cmd/action/main.go
	chmod +x dist/action-linux-amd64

run-local:
	GOOS=linux GOARCH=amd64 go run ./cmd/action/main.go

run:
	dist/action-linux-amd64

push-latest:
	git tag latest -f
	git push origin latest -f

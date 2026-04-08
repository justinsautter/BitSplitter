.PHONY: build test clean install

build:
	CGO_ENABLED=0 go build -o bitsplitter .

test:
	CGO_ENABLED=0 go test ./...

clean:
	rm -f bitsplitter

install: build
	cp bitsplitter /usr/local/bin/bitsplitter

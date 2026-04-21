.PHONY: all build test race vet fmt fmt-check tidy check clean

BINARY := grist
PKG := ./...

all: check build

build:
	go build -o $(BINARY) ./cmd/grist

test:
	go test $(PKG)

race:
	go test -race $(PKG)

vet:
	go vet $(PKG)

fmt:
	gofmt -s -w .

fmt-check:
	@out="$$(gofmt -s -l .)"; \
	if [ -n "$$out" ]; then \
		echo "gofmt needs to be run on:"; echo "$$out"; exit 1; \
	fi

tidy:
	go mod tidy

check: fmt-check vet race

clean:
	rm -f $(BINARY)

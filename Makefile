build:
	go build -o zirc ./cmd/zirc

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w cmd internal

clean:
	rm -f zirc

.PHONY: build test vet coverage

build: vet
	go install -v

test:
	go test -v -cover -race -coverprofile=coverage.out

vet:
	go vet

coverage:
	go tool cover -html=coverage.out

clean:
	go clean
	rm -f clog.log
	rm -f coverage.out

.PHONY: clean build test 

.PHONY: build:
	@echo "starging build..."
	go build -o . bin/bian-go-mock.go
	@echo "build completed."

.PHONY: test:
	go test ./...

.PHONY clean:
	rm -rf bin
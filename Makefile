BINARY_NAME=blowe-cli

.PHONY: test clean

run:
	go run main.go $(ARGS)

build:
	GOOS=darwin go build -o ${BINARY_NAME}-darwin main.go
	GOOS=linux go build -o ${BINARY_NAME}-linux main.go

test:
	go test ./...

clean:
	go clean
	rm -f ${BINARY_NAME}*
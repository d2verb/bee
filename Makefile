APP=bee

.PHONY: build
build: clean
	go build -o ${APP} main.go

.PHONY: run
run:
	go run -race main.go

.PHONY: test
test:
	go test ./lexer

.PHONY: clean
clean:
	go clean
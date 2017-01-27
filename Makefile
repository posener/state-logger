all: fmt test

CMD:=bin/fileserver

fmt:
	@go fmt ./...

test:
	@./go.test.sh

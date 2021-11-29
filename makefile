PROJECT_NAME=tee

.PHONY: build
build:
	go build -o ./bin/$(PROJECT_NAME) $(PROJECT_NAME).go

.PHONY: test
test:
	go test ./...

.PHONY: aliastee
aliastee:
	alias tee='bin/tee'

.PHONY: unaliastee
unaliastee:
	unalias tee

all: *.go gorpc/*.go utils/fs/*.go
	go build -gcflags="all=-N -l" -o protoc-gen-gorpc

PHONY: install
PHONY: uninstall

install:
	go install

uninstall:
	go clean
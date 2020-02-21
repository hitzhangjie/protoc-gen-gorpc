all: *.go gorpc/*.go gorpc/utils
	go build -gcflags="all=-N -l" -o protoc-gen-gorpc

.PHONY: install
.PHONY: uninstall

install:
	go install

uninstall:
	go clean -i

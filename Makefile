all: *.go gorpc/*.go utils/fs/*.go
	go build -gcflags="all=-N -l" -o protoc-gen-gorpc

.PHONY: install
.PHONY: uninstall

install:
	go install
	mkdir ~/.gorpc2 && cp -r ./install ~/.gorpc2/go

uninstall:
	go clean -i
	rm -rf ~/.gorpc2

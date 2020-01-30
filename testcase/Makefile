all: *.proto
	rm -rf hello
	protoc --gorpc_out=paths=source_relative:. hello.proto

.PHONY: clean

clean:
	rm -rf hello
	rm *.go



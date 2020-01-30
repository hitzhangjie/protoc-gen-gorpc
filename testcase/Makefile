all: *.proto
	rm -rf hello
	protoc --gorpc_out=. hello.proto

.PHONY: clean

clean:
	rm -rf hello
	rm *.go



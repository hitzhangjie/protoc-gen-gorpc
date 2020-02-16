all: *.proto
	protoc --gorpc_out=paths=source_relative:. hello.proto

.PHONY: clean

clean:
	rm -rf *.go *.mod rpc/



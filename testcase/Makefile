all: *.proto
	protoc --gorpc_out=plugins=gorpc:. hello.proto

.PHONY: clean

clean:
	rm *.go



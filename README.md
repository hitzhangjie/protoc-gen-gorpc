# protoc-gen-gorpc

## Introduction
protoc-gen-gorpc is a protoc plugin for generating [gorpc](https://github.com/hitzhangjie/go-rpc) code.

## Usage
how to use this plugin:

```
protoc --gorpc_out=plugins=gorpc:. ${filename}.proto
```

## Installation

you can install this plugin by run:

```
go install github.com/hitzhangjie/protoc-gen-gorpc
```

## More

please watch `go-rpc` project, it's a simple and beautiful rpc framework written in go.

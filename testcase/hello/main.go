package main

import (
	gorpc "github.com/hitzhangjie/go-rpc"

	pb "github.com/hello"
)

type helloSvrServiceImpl struct{}

func main() {

	s := gorpc.NewServer()

	pb.RegisterHelloSvrService(s, &helloSvrServiceImpl{})

	s.Serve()
}

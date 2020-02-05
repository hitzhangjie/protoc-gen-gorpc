package main_test

import (
	"context"
	"flag"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/hitzhangjie/go-rpc"
	"github.com/hitzhangjie/go-rpc/client"
	_ "github.com/hitzhangjie/go-rpc/http"

	pb "github.com/hello"
)

var (
	helloSvrClientProxy pb.HelloSvrClientProxy

	timeout       = flag.Duration("timeout", time.Second*0, "request timeout")
	network       = flag.String("network", "", "network")
	target        = flag.String("target", "", "target address, like ip://ip:port, cl5://mid:cid")
	svrConfigPath = flag.String("conf", "./gorpc_go.yaml", "config file path")
)

func init() {
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 默认使用配置文件中配置
	err := gorpc.LoadGlobalConfig(*svrConfigPath)
	if err == nil {
		for _, cfg := range gorpc.GlobalConfig().Client.Service {
			client.RegisterClientConfig(cfg.Callee, cfg)
		}
	}

	// 如果配置文件未提供，默认使用如下选项
	opts := []client.Option{
		client.WithProtocol("whisper"),
		client.WithNetwork("tcp4"),
		//client.WithTarget("ip://127.0.0.1:8000"),
		client.WithTimeout(time.Second * 2),
	}

	// 如果命令行选项由指定，覆盖上述选项
	if *timeout != time.Second*0 {
		opts = append(opts, client.WithTimeout(*timeout))
	}
	if *network != "" {
		opts = append(opts, client.WithNetwork(*network))
	}
	if *target != "" {
		opts = append(opts, client.WithTarget(*target))
	}

	helloSvrClientProxy = pb.NewHelloSvrClientProxy(opts...)
}

func Test_HelloSvr_Hello(t *testing.T) {
	ctx := context.TODO()

	tests := []struct {
		name    string
		req     *hello.Request
		wantRsp *hello.Response
		wantErr bool
	}{
		{"1-default", &hello.Request{}, &hello.Response{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := helloSvrClientProxy.Hello(ctx, tt.req)
			if tt.wantErr != (err != nil) {
				t.Errorf("wantErr = %v, err = %v, req:%s, rsp:%s", tt.wantErr, err, tt.req.String(), tt.wantRsp.String())
			}
			if !tt.wantErr && err == nil {
				if !reflect.DeepEqual(got, tt.wantRsp) {
					t.Errorf("got = %s, want = %s", got.String(), tt.wantRsp.String())
				}
			}
		})
	}
}

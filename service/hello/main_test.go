package main_test

import (
	"context"
	"fmt"
	sp_rpc_hello "github.com/socialpoint/envoyproxy-envoy-issues-11099/service/hello/pkg/SP/Rpc/Hello"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"testing"
)

const (
	addrEnvoyL1   = "envoy_l1:11001"
	addrEnvoyL2   = "envoy_l2:10001"
	validName     = "Lucas"
	invalidName   = ""
	expectedError = "rpc error: code = InvalidArgument desc = invalid name"
)

func TestHello_EnvoyProxyLevel2(t *testing.T) {
	fmt.Println("********** client -> L2 (reverse-bridge) -> backend")

	a := assert.New(t)
	r := require.New(t)

	conn, err := grpc.Dial(addrEnvoyL2, grpc.WithInsecure())
	r.NoError(err)

	cli := sp_rpc_hello.NewHelloServiceClient(conn)
	helloResponse, err := cli.Hello(context.Background(), &sp_rpc_hello.HelloRequest{Name: validName})
	a.NoError(err)

	fmt.Println(helloResponse)
	fmt.Print("********** " + "\n\n\n")
}

func TestHello_EnvoyProxyLevel2InvalidArgumentName(t *testing.T) {
	fmt.Println("********** client -> L2 (reverse-bridge) -> backend (invalid name)")

	a := assert.New(t)
	r := require.New(t)

	conn, err := grpc.Dial(addrEnvoyL2, grpc.WithInsecure())
	r.NoError(err)

	cli := sp_rpc_hello.NewHelloServiceClient(conn)
	_, err = cli.Hello(context.Background(), &sp_rpc_hello.HelloRequest{Name: invalidName})
	a.Error(err)

	a.Equal(expectedError, err.Error())
	fmt.Println(err.Error())
	fmt.Print("********** " + "\n\n\n")
}

func TestHello_EnvoyProxyLevel1(t *testing.T) {
	fmt.Println("********** client -> L1 -> L2 (reverse-bridge) -> backend")

	a := assert.New(t)
	r := require.New(t)

	conn, err := grpc.Dial(addrEnvoyL1, grpc.WithInsecure())
	r.NoError(err)

	cli := sp_rpc_hello.NewHelloServiceClient(conn)
	helloResponse, err := cli.Hello(context.Background(), &sp_rpc_hello.HelloRequest{Name: validName})
	a.NoError(err)

	fmt.Println(helloResponse)
	fmt.Print("********** " + "\n\n\n")
}

func TestHello_EnvoyProxyLevel1InvalidArgumentName(t *testing.T) {
	fmt.Println("********** client -> L1 -> L2 (reverse-bridge) -> backend (invalid name)")

	a := assert.New(t)
	r := require.New(t)

	conn, err := grpc.Dial(addrEnvoyL1, grpc.WithInsecure())
	r.NoError(err)

	cli := sp_rpc_hello.NewHelloServiceClient(conn)
	_, err = cli.Hello(context.Background(), &sp_rpc_hello.HelloRequest{Name: invalidName})
	a.Error(err)

	a.Equal(expectedError, err.Error())
	fmt.Println(err.Error())
	fmt.Print("********** " + "\n\n\n")
}

package rpc

import (
	"fmt"
	"google.golang.org/grpc"
)

// NewClient return a new rpc client
// 		NewClient("192.168.0.1:1234")
func NewClient(server string) (client *grpc.ClientConn, err error) {
	client, err = grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		err = fmt.Errorf("[RPC] get client error: %v", err)
		return
	}
	return
}

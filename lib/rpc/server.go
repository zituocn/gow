package rpc

import (
	"fmt"
	"github.com/zituocn/gow/lib/logy"
	"google.golang.org/grpc"
	"net"
)

// Server GRPCServer
type Server struct {
	Listener net.Listener
	Server   *grpc.Server
	Port     int //端口
}

// NewServer return a server
func NewServer(port int) (server *Server, err error) {
	if port == 0 {
		err = fmt.Errorf("[RPC] init failed：need port")
		return
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return
	}

	server = &Server{
		Listener: listener,
		Server:   grpc.NewServer(),
		Port:     port,
	}
	return
}

// Run run rpc server
func (m *Server) Run() {
	go func() {
		err := m.Server.Serve(m.Listener)
		if err != nil {
			logy.Error("[RPC] failed to listen:%v", err)
		}
	}()

}

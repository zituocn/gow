package rpc

import (
	"fmt"
	"github.com/zituocn/gow/lib/logy"
	"google.golang.org/grpc"
	"net"
)

// Server GRPC Server struct
type Server struct {
	Listener   net.Listener
	Server     *grpc.Server
	Port int

}

// NewServer returns a new server
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

func (m *Server) Register() {

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

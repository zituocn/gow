/*
sam
2022-09-18
*/

package gb

import (
	"errors"
	"fmt"
	"github.com/zituocn/gow/lib/gb/rd"
	"github.com/zituocn/gow/lib/logy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
)

// Server GRPC server struct
type Server struct {
	Listener net.Listener
	Server   *grpc.Server

	// grpc service name
	Name string

	// grpc service port
	port int

	// etcd endpoints
	etcdEndPoints []string

	// etcd lease
	lease      int64
	isRegister bool
}

// NewServer returns server
//	server,err:=gb.NewServer(opt)
//	server.Run()
func NewServer(opt *ServerOption) (*Server, error) {
	if opt.Port == 0 {
		return nil, fmt.Errorf("[RPC] init failed: need port")
	}
	listener, err := net.Listen("tcp", opt.address())
	if err != nil {
		return nil, err
	}
	if opt.IsRegister && len(opt.EtcdEndPoints) == 0 {
		return nil, errors.New("[RPC] service register need etcd endpoint")

	}
	if opt.Lease < 1 {
		opt.Lease = 3
	}
	var g *grpc.Server
	if opt.KeyFile != "" && opt.CertFile != "" {
		cred, err := credentials.NewServerTLSFromFile(opt.CertFile, opt.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("[RPC] load cred file error : %s", err.Error())

		}
		g = grpc.NewServer(grpc.Creds(cred))
	} else {
		g = grpc.NewServer()
	}
	server := &Server{
		Listener:      listener,
		Server:        g,
		Name:          opt.Name,
		etcdEndPoints: opt.EtcdEndPoints,
		port:          opt.Port,
		lease:         opt.Lease,
		isRegister:    opt.IsRegister,
	}
	return server, nil
}

// Run start grpc server
//	start service and complete service registration
func (s *Server) Run() {
	go func() {
		logy.Infof("[RPC] [%s] start grpc service listen on :%d", s.Name, s.port)
		err := s.Server.Serve(s.Listener)
		if err != nil {
			logy.Errorf("[RPC] failed to listen: %s", err.Error())
		}
	}()

	go func() {
		if s.isRegister {
			logy.Info("[RPC] start service registration ...")
			_, err := rd.NewServiceRegister(&rd.ServiceOption{
				Endpoints: s.etcdEndPoints,
				Lease:     s.lease,
				Prefix:    s.Name,
				Port:      s.port,
			})
			if err != nil {
				logy.Errorf("[RPC] service register failed : %s", err.Error())
				return
			}
			logy.Info("[RPC] service registration succeeded")
		}
	}()
}

// ServerOption grpc option
type ServerOption struct {
	// * grpc service name
	Name string

	// * grpc service port
	Port int

	// grpc service IP address
	IP string

	// etcd endpoints
	EtcdEndPoints []string

	// etcd lease
	Lease int64

	//  whether registration service is required
	IsRegister bool

	// crt or pem file
	CertFile string

	// private key file
	KeyFile string
}

// Address returns ip:port string
func (o *ServerOption) address() string {
	return fmt.Sprintf("%s:%d", o.IP, o.Port)
}

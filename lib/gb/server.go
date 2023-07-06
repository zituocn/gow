/*
server.go

grpc server 封装
sam
2023-01-30

*/

package gb

import (
	"errors"
	"fmt"
	"github.com/zituocn/logx"
	"google.golang.org/grpc"
	"net"
)

var (
	grpcScheme = "grpc"
)

type GrpcServer struct {
	Listener net.Listener

	Server *grpc.Server

	ServerName string

	port int

	etcdEndpoints []string

	lease int64

	isRegister bool
}

// NewGrpcServer 返回新的grpc register
func NewGrpcServer(opt *GrpcServerOption) (*GrpcServer, error) {
	if opt.Port == 0 {
		return nil, errors.New("[grpc] need port")
	}
	if opt.IsRegister && len(opt.EtcdEndpoints) == 0 {
		return nil, errors.New("[grpc] need etcd endpoint")
	}
	if !opt.IsRegister && opt.IP == "" {
		return nil, errors.New("[grpc] need ip address")
	}
	if opt.Lease < 1 {
		opt.Lease = 3
	}
	listener, err := net.Listen("tcp", opt.address())
	if err != nil {
		return nil, err
	}
	g := grpc.NewServer()
	server := &GrpcServer{
		Listener:      listener,
		Server:        g,
		ServerName:    opt.ServerName,
		etcdEndpoints: opt.EtcdEndpoints,
		port:          opt.Port,
		lease:         opt.Lease,
		isRegister:    opt.IsRegister,
	}
	return server, nil
}

// Run 运行
func (s *GrpcServer) Run() {
	go func() {
		logx.Infof("[grpc] [%s] start grpc service listen on :%d", s.ServerName, s.port)
		err := s.Server.Serve(s.Listener)
		if err != nil {
			logx.Errorf("[grpc] failed to listen: %s", err.Error())
		}
	}()
	go func() {
		if s.isRegister {
			logx.Info("[grpc] start registration ...")
			addr, err := GetLocalIP()
			if err != nil {
				logx.Errorf("[grpc] get location ip error : %s", err.Error())
				return
			}
			val := fmt.Sprintf("%s:%d", addr, s.port)
			_, err = NewRegister(&RegisterOption{
				EtcdEndpoints: s.etcdEndpoints,
				Schema:        grpcScheme,
				Lease:         s.lease,
				Key:           s.ServerName,
				Val:           val,
			})
			if err != nil {
				logx.Errorf("[grpc] register failed: %s", err.Error())
				return
			}
			logx.Info("[grpc] service registered successfully")
		}
	}()

}

// GrpcServerOption grpc register 参数
type GrpcServerOption struct {

	// grpc ServerName 服务名
	ServerName string

	// 端口
	Port int

	// IP地址
	IP string

	// EtcdEndpoints  etcd地址，如果为多个地址时，需要集群部署etcd
	EtcdEndpoints []string

	// Lease etcd 租约续期时间
	Lease int64

	// IsRegister 是否注册到etcd
	//	为true时，可不填写传入ip地址
	IsRegister bool
}

func (o *GrpcServerOption) address() string {
	return fmt.Sprintf("%s:%d", o.IP, o.Port)
}

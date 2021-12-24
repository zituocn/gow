package rpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/zituocn/gow/lib/logy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// Server GRPC Server struct
type Server struct {
	Listener   net.Listener
	Server     *grpc.Server
	Port       int    // 端口
	Name       string // 服务昵称
	EtcdAddr   []string
	isRegister bool
}

// ServerArg grpc 服务端参数
type ServerArg struct {
	IP       string
	Port     int
	Name     string
	EtcdAddr string // 多个地址必须使用","分开; 如: 192.168.0.1:2379,192.168.0.2:2379
	Register bool   // 是否启用服务注册
}

// NewServer returns a new server
func NewServer(grpcAddr ServerArg) (server *Server, err error) {
	if grpcAddr.Port == 0 {
		err = fmt.Errorf("[RPC] init failed：need port")
		return
	}
	listener, err := net.Listen("tcp", grpcAddr.String())
	if err != nil {
		return
	}
	if grpcAddr.Register && len(grpcAddr.EtcdAddr) < 1 {
		err = errors.New("[RPC]服务注册缺少etcd参数")
		return
	}

	newServer := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor), grpc.StreamInterceptor(streamInterceptor))
	server = &Server{
		Listener:   listener,
		Server:     newServer,
		Port:       grpcAddr.Port,
		Name:       grpcAddr.Name,
		EtcdAddr:   strings.Split(grpcAddr.EtcdAddr, ","),
		isRegister: grpcAddr.Register,
	}

	return
}

// String  ServerArg ip -> string
func (addr ServerArg) String() string {
	return fmt.Sprintf("%s:%d", addr.IP, addr.Port)
}

// SetName set grpc server name
func (m *Server) SetName(name string) *Server {
	m.Name = name
	return m
}

// SetPort set grpc port
func (m *Server) SetPort(port int) *Server {
	m.Port = port
	return m
}

// SetEtcdAddr set etcd address
func (m *Server) SetEtcdAddr(etcdAddr string) *Server {
	m.EtcdAddr = strings.Split(etcdAddr, ",")
	return m
}

// OpenRegister set isRegister=true
func (m *Server) OpenRegister() *Server {
	m.isRegister = true
	return m
}

// Run run rpc server
func (m *Server) Run() {
	go func() {
		err := m.Server.Serve(m.Listener)
		if err != nil {
			logy.Errorf("[RPC] failed to listen:%v", err)
		}
	}()

	if m.isRegister {
		m.register()
	}
}

// Register 将grpc服务信息注册到etcd
func (m *Server) register() {
	if len(m.Name) < 1 {
		panic("[RPC] grpc服务注册失败, 未设置服务名称. ")
	}

	ip, err := GetIp()
	if err != nil {
		panic("[RPC] grpc服务注册失败, 获取本机ip失败, err = " + err.Error())
	}

	if len(m.EtcdAddr) < 1 {
		panic("[RPC] grpc服务注册失败, etcd地址为空;")
	}

	NewEtcdCli(m.EtcdAddr)
	key := fmt.Sprintf("%s%s:%d", serverNameKey(m.Name), ip, m.Port)
	err = etcdConn.Register(key, "0")
	if err != nil {
		panic("[致命启动错误] 服务注册失败 err = " + err.Error())
	} else {
		logy.Infof("[RPC] 服务注册 Register Succeed; key = %s ", key)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		// 接收到进程退出信号量,解除租约
		_ = etcdConn.UnRegister(key)
		if i, ok := s.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}

	}()

}

// unaryInterceptor  中间件打印日志
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	startTime := time.Now()
	pr, _ := peer.FromContext(ctx)
	md, _ := metadata.FromIncomingContext(ctx)
	clientName := getValue(md, "clientname")
	serviceName := getValue(md, "servicename")
	serviceIp, _ := GetIp()
	m, err := handler(ctx, req)
	if err != nil {
		logy.Errorf("[GRPC] %s(%v)->%s(%v) | %s | err = %v",
			clientName,
			pr.Addr.String(),
			serviceName,
			serviceIp,
			info.FullMethod,
			err)
	} else {
		logy.Infof("[GRPC] %13v | %s(%v)->%s(%v) | %s ",
			time.Now().Sub(startTime),
			clientName,
			pr.Addr.String(),
			serviceName,
			serviceIp,
			info.FullMethod)
	}
	return m, err
}

// wrappedStream
type wrappedStream struct {
	grpc.ServerStream
}

// RecvMsg  receive message
//	returns error
func (w *wrappedStream) RecvMsg(m interface{}) error {
	return w.ServerStream.RecvMsg(m)
}

// SendMsg send message
//	returns error
func (w *wrappedStream) SendMsg(m interface{}) error {
	return w.ServerStream.SendMsg(m)
}

// newWrappedStream
func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}

// streamInterceptor
func streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return fmt.Errorf("missing metadata")
	}
	if !valid(md["authorization"]) {
		return fmt.Errorf("invalid token")
	}
	err := handler(srv, newWrappedStream(ss))
	if err != nil {
		logy.Errorf("[RPC] failed with error :%v", err)
	}
	return err
}

// valid validates the authorization.
func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	// Perform the token validation here. For the sake of this example, the code
	// here forgoes any of the usual OAuth2 token validation and instead checks
	// for a token matching an arbitrary string.
	return token == "some-secret-token"
}

// getValue
func getValue(md metadata.MD, key string) string {
	if v, ok := md[key]; ok {
		if len(v) > 0 {
			return v[0]
		}
	}
	return ""
}

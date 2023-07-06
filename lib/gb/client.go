/*
client.go

grpc client 封装
sam
2023-01-30

*/

package gb

import (
	"context"
	"errors"
	"fmt"
	"github.com/zituocn/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/resolver"
	"time"
)

type ConnType uint

const (

	// DefaultType 传统模式的grpc连接方式
	DefaultType ConnType = iota + 1

	// BalanceType 负载均衡的grpc连接方式
	//	需要使用etcd做为服务注册与发现
	BalanceType
)

type GrpcClientOption struct {

	// SererAddr grpc服务的地址
	ServerAddr string

	// EtcdEndpoint etcd地址
	EtcdEndpoint []string

	// ServerName 服务名
	ServerName string

	// ClientName 客户端名称
	ClientName string

	// ConnType 连接类型
	ConnType ConnType
}

type GrpcClient struct {
	serverAddr   string
	etcdEndpoint []string
	serverName   string
	clientName   string
	connType     ConnType
}

func NewGrpcClient(opt *GrpcClientOption) (*GrpcClient, error) {
	if opt.ServerName == "" {
		return nil, errors.New("[grpc] need register name")
	}
	if opt.ClientName == "" {
		return nil, errors.New("[grpc] need discovery name")
	}
	if opt.ConnType == DefaultType && opt.ServerAddr == "" {
		return nil, errors.New("[grpc] if conn type is DefaultType, need ServerAddr")

	}
	if opt.ConnType == BalanceType && len(opt.EtcdEndpoint) == 0 {
		return nil, errors.New("[grpc] if conn type is BalanceType, need etcd endpoint")
	}

	client := &GrpcClient{
		serverAddr:   opt.ServerAddr,
		etcdEndpoint: opt.EtcdEndpoint,
		clientName:   opt.ClientName,
		serverName:   opt.ServerName,
		connType:     opt.ConnType,
	}
	return client, nil
}

// GetConn returns grpc clientConn,ctx and error
func (c *GrpcClient) GetConn() (*grpc.ClientConn, context.Context, error) {
	conn, err := c.newClientConn()
	ctx := setCtx(c.serverName, c.clientName, conn)
	return conn, ctx, err
}

func (c *GrpcClient) newClientConn() (conn *grpc.ClientConn, err error) {
	switch c.connType {
	case DefaultType:
		conn, err = grpc.Dial(
			c.serverAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(unaryInterceptorClient),
		)
		if err != nil {
			return nil, err
		}
	case BalanceType:
		key := "/" + grpcScheme + "/" + c.serverName
		rd, err := NewDiscovery(grpcScheme, key, c.etcdEndpoint, 5)
		if err != nil {
			return nil, err
		}
		resolver.Register(rd)
		conn, err = grpc.Dial(
			fmt.Sprintf("%s:///%s", rd.Scheme(), c.serverName),
			grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, err
		}
	default:
		err = errors.New("unknown conn type")
		return nil, err
	}
	return
}

func unaryInterceptorClient(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	for _, o := range opts {
		_, ok := o.(grpc.PerRPCCredsCallOption)
		if ok {
			break
		}
	}
	md, _ := metadata.FromOutgoingContext(ctx)
	clientName := getValue(md, "clientname")
	clientIp, _ := GetLocalIP()
	serviceName := getValue(md, "servicename")
	startTime := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		logx.Errorf("%4s | %13v | %10s:%-10s -> %10s:%-10s | %s | %s",
			"rpc",
			time.Now().Sub(startTime),
			clientName,
			clientIp,
			serviceName,
			cc.Target(),
			method,
			err.Error())
	} else {
		logx.Infof("%4s | %13v | %10s:%-10s -> %10s:%-10s | %s ", "rpc", time.Now().Sub(startTime), clientName, clientIp, serviceName, cc.Target(), method)
	}
	return err
}

func getValue(md metadata.MD, key string) string {
	if v, ok := md[key]; ok {
		if len(v) > 0 {
			return v[0]
		}
	}
	return ""
}

// setCtx set context
func setCtx(serviceName, clientName string, grpcConn *grpc.ClientConn) context.Context {
	if grpcConn == nil {
		return nil
	}
	kv := []string{
		"clientName", clientName,
		"serviceName", serviceName,
	}
	return metadata.NewOutgoingContext(context.Background(), metadata.Pairs(kv...))
}

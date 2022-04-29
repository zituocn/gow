package rpc

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zituocn/gow/lib/logy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Client GRPC Client struct
type Client struct {
	serverAddr  string
	etcdAddr    []string
	clientName  string
	serviceName string
	times       int
	retry       int
	retryTime   time.Duration
}

// ClientArg GRPC Client arg struct
type ClientArg struct {
	ServerAddr  string		// 直连  *必填
	EtcdAddr    string		// 使用发现服务 *必填
	ClientName  string
	ServiceName string		// *必填
}

// NewClient return grpc Client obj
func NewClient(dis ClientArg) (*Client, error) {
	etcdAddr := strings.Split(dis.EtcdAddr, ",")
	if len(etcdAddr) < 1 {
		return nil, errors.New("[RPC] 参数错误: etcdAddress is null")
	}
	if len(dis.ServiceName) < 1 {
		return nil, fmt.Errorf("[RPC] 参数错误: ServiceName is null")
	}

	return &Client{
		serverAddr:  dis.ServerAddr,
		etcdAddr:    etcdAddr,
		clientName:  dis.ClientName,
		serviceName: dis.ServiceName,
		times:       0,		// timestamp
		retry:       20,	// retry 20 times
		retryTime:   50 * time.Millisecond,
	}, nil
}

// Conn 直连
// return grpc.ClientConn, context.Context, err
func (c *Client) Conn() (client *grpc.ClientConn, ctx context.Context, err error) {
	client, err = newClient(c.serverAddr)
	ctx = setCtx(c.serviceName, c.clientName, client)
	return
}

// Min  通过etcd发现服务获取grpc连接; 负载均衡 - 最小连接数法;
// return grpc.ClientConn, context.Context, err
func (c *Client) Min() (client *grpc.ClientConn, ctx context.Context, err error) {
	// 避免一直重试
	if c.times > c.retry {
		err = errors.New("[ETCD] 没有发现服务")
		return
	}

	NewEtcdCli(c.etcdAddr)
	grpcIPKey, _ := etcdConn.GetMinKey(serverNameKey(c.serviceName))
	grpcIPList := strings.Split(grpcIPKey, "/")
	if len(grpcIPList) < 1 {
		time.Sleep(c.retryTime)
		c.times++
		return c.Min()
	}

	grpcIP := grpcIPList[len(grpcIPList)-1]
	client, err = newClient(grpcIP)
	if err != nil || client == nil {
		time.Sleep(c.retryTime)
		c.times++
		return c.Min()
	}

	// 使用GetMinKey方式需要执行GetMinKeyCallBack
	_ = etcdConn.GetMinKeyCallBack(grpcIPKey)
	ctx = setCtx(c.serviceName, c.clientName, client)
	return
}

// Rand  通过etcd发现服务获取grpc连接; 负载均衡 - 随机法;
// return grpc.ClientConn, context.Context, err
func (c *Client) Rand() (client *grpc.ClientConn, ctx context.Context, err error) {
	if c.times > c.retry {
		err = errors.New("[ETCD] 没有发现服务")
		return
	}

	NewEtcdCli(c.etcdAddr)
	serviceNameKey := serverNameKey(c.serviceName)
	grpcIP, _ := etcdConn.GetRandKey(serviceNameKey)
	client, err = newClient(grpcIP)
	if err != nil || client == nil {
		time.Sleep(c.retryTime) // 连不上可能是服务还在启动中, 等待50ms从新获取
		c.times++
		return c.Rand()
	}

	ctx = setCtx(c.serviceName, c.clientName, client)
	return
}

// newClient return a new rpc client
func newClient(server string) (client *grpc.ClientConn, err error) {
	client, err = grpc.Dial(
		server,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(unaryInterceptorClient),
		grpc.WithStreamInterceptor(streamInterceptorClient))
	if err != nil {
		err = fmt.Errorf("[RPC] get client error: %v", err)
		return
	}
	return
}

// unaryInterceptorClient 中间件打印日志
func unaryInterceptorClient(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	for _, o := range opts {
		_, ok := o.(grpc.PerRPCCredsCallOption)
		if ok {
			break
		}
	}

	md, _ := metadata.FromOutgoingContext(ctx)
	clientName := getValue(md, "clientname")
	clientIp, _ := GetIp()
	serviceName := getValue(md, "servicename")
	startTime := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		logy.Errorf("[GRPC] %s(%v)->%s(%v) | %s | err = %v",
			clientName,
			clientIp,
			serviceName,
			cc.Target(),
			method,
			err)
	} else {
		logy.Infof("[GRPC] %13v | %s(%v)->%s(%v) | %s ",
			time.Now().Sub(startTime),
			clientName,
			clientIp,
			serviceName,
			cc.Target(),
			method)
	}
	return err
}

// wrappedStreamClient
type wrappedStreamClient struct {
	grpc.ClientStream
}

// RecvMsg  receive message
//	returns error
func (w *wrappedStreamClient) RecvMsg(m interface{}) error {
	return w.ClientStream.RecvMsg(m)
}

// SendMsg send message
//	returns error
func (w *wrappedStreamClient) SendMsg(m interface{}) error {
	return w.ClientStream.SendMsg(m)
}

// newWrappedStreamClient
func newWrappedStreamClient(s grpc.ClientStream) grpc.ClientStream {
	return &wrappedStreamClient{s}
}

// streamInterceptorClient is an example stream interceptor.
func streamInterceptorClient(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	for _, o := range opts {
		_, ok := o.(*grpc.PerRPCCredsCallOption)
		if ok {
			break
		}
	}

	s, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		return nil, err
	}
	return newWrappedStreamClient(s), nil
}

// setCtx set context
func setCtx(serviceName, myName string, grpcConn *grpc.ClientConn) context.Context {
	if grpcConn == nil {
		return nil
	}
	kv := []string{
		//"RequestId", ID(),
		"ClientName", myName,
		"ServiceName", serviceName,
	}
	return metadata.NewOutgoingContext(context.Background(), metadata.Pairs(kv...))
}

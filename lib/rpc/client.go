package rpc

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/zituocn/gow/lib/logy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	etcdAddrNilErr = fmt.Errorf("[ETCD参数错误] etcd Addr is null; ex: 192.168.0.1:2379,192.168.0.2:2379")
	notServiceErr = fmt.Errorf("[发现服务失败] 没有发现服务")
)

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
	clientIp,_ := GetIp()
	serviceName := getValue(md, "servicename")
	//requestId := getValue(md, "requestid")
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

// RecvMsg
func (w *wrappedStreamClient) RecvMsg(m interface{}) error {
	log.Println("Receive a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ClientStream.RecvMsg(m)
}

// SendMsg
func (w *wrappedStreamClient) SendMsg(m interface{}) error {
	log.Println("Send a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
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

// setCtx
func setCtx(serviceName, myName string, grpcConn *grpc.ClientConn) context.Context {
	if grpcConn == nil {
		return nil
	}
	kv := []string {
		//"RequestId", ID(),
		"ClientName", myName,
		"ServiceName", serviceName,
	}
	return metadata.NewOutgoingContext(context.Background(), metadata.Pairs(kv...))
}


// discovery 发现服务
type discovery struct {
	serverAddr string
	etcdAddr []string
	clientName string
	serviceName string
	isLog bool
	times int
	retry int // 重试次数
	retryTime time.Duration
}

// DiscoveryArg 创建发现服务对象参数
type ClientArg struct {
	ServerAddr  string
	EtcdAddr    string
	ClientName  string
	ServiceName string
	OpenLog     bool
}

// NewClient 创建客户端对象
func NewClient(dis ClientArg) (*discovery, error) {
	etcdAddr := strings.Split(dis.EtcdAddr, ",")
	if len(etcdAddr) < 1 {
		return nil, etcdAddrNilErr
	}
	if len(dis.ServiceName) < 1 {
		return nil, fmt.Errorf("[Grpc] 参数err: ServiceName is null.")
	}

	return &discovery{
		serverAddr: dis.ServerAddr,
		etcdAddr: etcdAddr,
		clientName: dis.ClientName,
		serviceName: dis.ServiceName,
		isLog: dis.OpenLog,
		times: 0,
		retry: 20,
		retryTime: 50*time.Millisecond,
	}, nil
}

// Conn 获取连接
func (c *discovery) Conn() (client *grpc.ClientConn, ctx context.Context, err error) {
	client, err = newClient(c.serverAddr)
	ctx = context.Background()
	return
}


// Min  发现服务获取grpc连接; 负载均衡 - 最小连接数法;
func (c *discovery) Min() (client *grpc.ClientConn, ctx context.Context, err error) {
	// 避免一直重试
	if c.times > c.retry {
		err = notServiceErr
	}
	NewEtcdCli(c.etcdAddr)
	grpcIPKey,_ := etcdConn.GetMinKey(serverNameKey(c.serviceName))
	grpcIPList := strings.Split(grpcIPKey, "/")
	if len(grpcIPList) < 1 {
		time.Sleep(c.retryTime)  // 没有获取到服务地址, 可能是服务还在启动中, 等待50ms从新获取
		c.times++
		return c.Min()
	}
	grpcIP := grpcIPList[len(grpcIPList)-1]
	client, err = newClient(grpcIP)
	if err != nil || client == nil {
		time.Sleep(c.retryTime)  // 连不上可能是服务还在启动中, 等待50ms从新获取
		c.times++
		return c.Min()
	}

	// 使用GetMinKey方式需要执行GetMinKeyCallBack
	_=etcdConn.GetMinKeyCallBack(grpcIPKey)
	if c.isLog {
		ctx = setCtx(c.serviceName, c.clientName, client)
	} else {
		ctx = context.Background()
	}

	return
}

// Rand  发现服务获取grpc连接; 负载均衡 - 随机法;
func (c *discovery) Rand() (client *grpc.ClientConn, ctx context.Context, err error) {
	if c.times > c.retry {
		err = notServiceErr
	}
	NewEtcdCli(c.etcdAddr)
	serviceNameKey := serverNameKey(c.serviceName)
	grpcIP, _ := etcdConn.GetRandKey(serviceNameKey)
	log.Println("grpcIP = ", grpcIP)
	client, err = newClient(grpcIP)
	if err != nil || client == nil {
		time.Sleep(c.retryTime) // 连不上可能是服务还在启动中, 等待50ms从新获取
		c.times++
		return c.Rand()
	}

	if c.isLog {
		ctx = setCtx(c.serviceName, c.clientName, client)
	} else {
		ctx = context.Background()
	}
	return
}

// TODO  Discovery.Poll     发现服务获取grpc连接; 负载均衡 - 轮询法;

// TODO  Discovery.HashIP   发现服务获取grpc连接; 负载均衡 - 源地址哈希法;

// TODO  Discovery.Hash     发现服务获取grpc连接; 负载均衡 - 一致性哈希法;

// TODO  Discovery.WRR      发现服务获取grpc连接; 负载均衡 - 加权轮询法;

// TODO  Discovery.RandWH   发现服务获取grpc连接; 负载均衡 - 加权随机法;

// TODO  Discovery.Fastest  发现服务获取grpc连接; 负载均衡 - 最快响应速度法;
/*
sam
2022-09-18
*/

package gb

import (
	"context"
	"errors"
	"fmt"
	"github.com/zituocn/gow/lib/gb/rd"
	"github.com/zituocn/gow/lib/logy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/resolver"
	"time"
)

// ConnType grpc conn type
type ConnType uint

const (

	// DefaultType get default client conn
	DefaultType ConnType = iota + 1

	// BalanceType  get client conn from etcd
	//	need etcdEndPoint
	BalanceType
)

type Client struct {
	serverAddr   string
	etcdEndPoint []string

	serviceName string
	clientName  string

	certFile string

	connType ConnType
}

type ClientOption struct {
	// grpc server addr host:port *
	ServerAddr string

	// etcd endpoint
	EtcdEndPoint []string

	// grpc server name *
	ServiceName string

	// * grpc client name
	ClientName string

	// cert file if you need
	CertFile string

	// returns conn type
	ConnType ConnType
}

// NewClient returns  new client and error
func NewClient(opt *ClientOption) (*Client, error) {
	if opt.ServiceName == "" {
		return nil, errors.New("[RPC] need ServiceName")
	}
	if opt.ServerAddr == "" {
		return nil, errors.New("[RPC] need ServerAddr")
	}
	if opt.ConnType < 1 {
		return nil, errors.New("[RPC] need conn type")
	}
	client := &Client{
		serverAddr:   opt.ServerAddr,
		etcdEndPoint: opt.EtcdEndPoint,
		clientName:   opt.ClientName,
		serviceName:  opt.ServiceName,
		certFile:     opt.CertFile,
		connType:     opt.ConnType,
	}
	return client, nil
}

// GetConn get the client conn
//	returns conn and error
func (c *Client) GetConn() (*grpc.ClientConn, context.Context, error) {
	conn, err := c.newClientConn()
	ctx := setCtx(c.serviceName, c.clientName, conn)
	return conn, ctx, err
}

func (c *Client) newClientConn() (*grpc.ClientConn, error) {
	switch c.connType {
	case DefaultType:
		if c.certFile != "" {
			// use credential
			cred, err := credentials.NewClientTLSFromFile(c.certFile, c.serviceName)
			if err != nil {
				return nil, fmt.Errorf("[RPC] failed to validate certificate :%s", err.Error())
			}
			conn, err := grpc.Dial(
				c.serverAddr,
				grpc.WithTransportCredentials(cred),
				grpc.WithUnaryInterceptor(unaryInterceptorClient),
				grpc.WithStreamInterceptor(streamInterceptorClient))
			if err != nil {
				return nil, fmt.Errorf("[RPC] get default client conn with credentials from etcd error :%s", err.Error())
			}
			return conn, nil
		} else {
			conn, err := grpc.Dial(
				c.serverAddr,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithUnaryInterceptor(unaryInterceptorClient),
				grpc.WithStreamInterceptor(streamInterceptorClient))
			if err != nil {
				return nil, fmt.Errorf("[RPC] get default client conn error :%s", err.Error())
			}
			return conn, nil
		}
	case BalanceType:
		rb := rd.NewServiceDiscovery(c.etcdEndPoint)
		resolver.Register(rb)
		if c.certFile != "" {
			cred, err := credentials.NewClientTLSFromFile(c.certFile, c.serviceName)
			if err != nil {
				return nil, fmt.Errorf("[RPC] failed to validate certificate : %s", err.Error())
			}
			conn, err := grpc.Dial(
				fmt.Sprintf("%s:///%s", rb.Scheme(), c.serviceName),
				grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
				grpc.WithTransportCredentials(cred),
				grpc.WithUnaryInterceptor(unaryInterceptorClient),
				grpc.WithStreamInterceptor(streamInterceptorClient))

			if err != nil {
				return nil, fmt.Errorf("[RPC] get client conn with credentials from etcd  error : %s", err.Error())
			}
			return conn, nil
		} else {
			conn, err := grpc.Dial(
				fmt.Sprintf("%s:///%s", rb.Scheme(), c.serviceName),
				grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithUnaryInterceptor(unaryInterceptorClient),
				grpc.WithStreamInterceptor(streamInterceptorClient))
			if err != nil {
				return nil, fmt.Errorf("[RPC] get client conn from etcd error :%s", err.Error())
			}
			return conn, nil
		}
	default:
		return nil, fmt.Errorf("unknown conn Type")
	}
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
	clientIp, _ := rd.GetLocalIP()
	serviceName := getValue(md, "servicename")
	startTime := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		logy.Errorf("[GRPC] [%s] %s -> [%s]  %s  | %s | err = %s",
			clientName,
			clientIp,
			serviceName,
			cc.Target(),
			method,
			err.Error())
	} else {
		logy.Infof("[GRPC] %8s | %s (%s) -> %s (%s) | %s ",
			time.Now().Sub(startTime),
			clientName,
			clientIp,
			serviceName,
			cc.Target(),
			method)
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

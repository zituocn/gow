package rd

import (
	"context"
	"fmt"
	"github.com/zituocn/gow/lib/logy"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"math/rand"
	"net"
	"sync"
	"time"
)

var (
	timeOut = 3
)

type ServiceDiscovery struct {
	cli          *clientv3.Client
	etcdEndPoint []string
	serverList   sync.Map
	serviceName  string
	prefix       string
}

func NewClientDiscovery(serviceName string, endpoints []string) (*ServiceDiscovery, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Duration(timeOut) * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &ServiceDiscovery{
		cli:          cli,
		serviceName:  serviceName,
		etcdEndPoint: endpoints,
	}, nil
}

func (s *ServiceDiscovery) Build() error {
	s.prefix = "/" + schema + "/" + s.serviceName + "/"
	resp, err := s.cli.Get(context.Background(), s.prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for _, ev := range resp.Kvs {
		s.SetServiceList(string(ev.Key), string(ev.Value))
	}
	go s.watcher()

	return nil
}

func (s *ServiceDiscovery) watcher() {
	rch := s.cli.Watch(context.Background(), s.prefix, clientv3.WithPrefix())
	for wResp := range rch {
		for _, ev := range wResp.Events {
			switch ev.Type {
			case mvccpb.PUT: //PUT
				s.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE: //DELETE
				s.DelServiceList(string(ev.Kv.Key))
			}
		}
	}
}

// Close 关闭
func (s *ServiceDiscovery) Close() {
	if s.cli != nil {
		err := s.cli.Close()
		if err != nil {
			logy.Errorf("close etcd error :%v", err)
		}
	}
}

func (s *ServiceDiscovery) getServices() []resolver.Address {
	addrs := make([]resolver.Address, 0)
	s.serverList.Range(func(k, v interface{}) bool {
		addrs = append(addrs, v.(resolver.Address))
		return true
	})
	return addrs
}

func (s *ServiceDiscovery) DelServiceList(key string) {
	s.serverList.Delete(key)
}

func (s *ServiceDiscovery) SetServiceList(key, val string) {
	s.serverList.Store(key, resolver.Address{Addr: val})
}

var (
	rd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// GetServerAddr return grpc server ip
func (s *ServiceDiscovery) GetServerAddr() (string, error) {
	addrs := s.getServices()
	if len(addrs) == 0 {
		return "", fmt.Errorf("service address not found from %v", s.etcdEndPoint)
	}
	rr := new(RoundRobin)
	rr.Clear()
	for _, item := range addrs {
		rr.Add(item.Addr)
	}
	rr.length = len(rr.ss)
	addr := rr.Pick()
	return addr, nil
}

/*
RoundRobin
*/

type RoundRobin struct {
	ss     []string
	next   int
	length int
}

func (r *RoundRobin) Clear() {
	r.ss = r.ss[:0]
}

func (r *RoundRobin) Add(s string) {
	if s != "" {
		r.ss = append(r.ss, s)
	}
}

func (r *RoundRobin) Pick() string {
	if r.length == 0 {
		return ""
	}
	index := rd.Intn(r.length)
	return r.ss[index]
}

// tcpHealth simple grpc health check
func tcpHealth(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
	if err != nil || conn == nil {
		return false
	} else {
		_ = conn.Close()
		return true
	}
}

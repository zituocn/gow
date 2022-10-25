package rd

//
//import (
//	"context"
//	"errors"
//	"fmt"
//	"go.etcd.io/etcd/api/v3/mvccpb"
//	"sync"
//	"time"
//
//	"github.com/zituocn/gow/lib/logy"
//
//	clientv3 "go.etcd.io/etcd/client/v3"
//	"google.golang.org/grpc/resolver"
//)
//
//var (
//	timeOut = 3
//)
//
//// ServiceDiscovery implement the Builder interface
//type ServiceDiscovery struct {
//	cli        *clientv3.Client
//	cc         resolver.ClientConn
//	serverList sync.Map
//	prefix     string
//}
//
//// NewServiceDiscovery returns a new resolver.Builder
//func NewServiceDiscovery(endpoints []string) (resolver.Builder, error) {
//	cli, err := clientv3.New(clientv3.Config{
//		Endpoints:   endpoints,
//		DialTimeout: time.Duration(timeOut) * time.Second,
//	})
//	if err != nil {
//		return nil, err
//	}
//	// check etcd timeout
//	timeOutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second/2)
//	defer cancel()
//	_, err = cli.Status(timeOutCtx, endpoints[0])
//	if err != nil {
//		return nil, fmt.Errorf("[ETCD] connection timed out: %s", err.Error())
//	}
//	return &ServiceDiscovery{
//		cli: cli,
//	}, nil
//}
//
//func (s *ServiceDiscovery) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
//	s.cc = cc
//	s.prefix = "/" + target.URL.Scheme + target.URL.Path + "/"
//	resp, err := s.cli.Get(context.Background(), s.prefix, clientv3.WithPrefix())
//	if err != nil {
//		return nil, err
//	}
//	if resp == nil || len(resp.Kvs) == 0 {
//		err = errors.New("no service node available")
//		logy.Errorf("[ETCD] %v", err)
//		return nil, err
//	}
//	for _, ev := range resp.Kvs {
//		s.SetServiceList(string(ev.Key), string(ev.Value))
//	}
//	_ = s.cc.UpdateState(resolver.State{Addresses: s.getServices()})
//	go s.watcher()
//	return s, nil
//}
//
//func (s *ServiceDiscovery) SetServiceList(key, val string) {
//	s.serverList.Store(key, resolver.Address{Addr: val})
//	_ = s.cc.UpdateState(resolver.State{Addresses: s.getServices()})
//}
//
//func (s *ServiceDiscovery) watcher() {
//	rch := s.cli.Watch(context.Background(), s.prefix, clientv3.WithPrefix())
//	for wResp := range rch {
//		for _, ev := range wResp.Events {
//			switch ev.Type {
//			case mvccpb.PUT: //PUT
//				s.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
//			case mvccpb.DELETE: //DELETE
//				s.DelServiceList(string(ev.Kv.Key))
//			}
//		}
//	}
//}
//
//func (s *ServiceDiscovery) ResolveNow(rn resolver.ResolveNowOptions) {
//	logy.Debug("ResolveNow")
//}
//
//// Scheme return schema
//func (s *ServiceDiscovery) Scheme() string {
//	return schema
//}
//
//// Close 关闭
//func (s *ServiceDiscovery) Close() {
//	_ = s.cli.Close()
//}
//
//func (s *ServiceDiscovery) getServices() []resolver.Address {
//	addrs := make([]resolver.Address, 0)
//	s.serverList.Range(func(k, v interface{}) bool {
//		addrs = append(addrs, v.(resolver.Address))
//		return true
//	})
//	return addrs
//}
//
//func (s *ServiceDiscovery) DelServiceList(key string) {
//	s.serverList.Delete(key)
//	_ = s.cc.UpdateState(resolver.State{Addresses: s.getServices()})
//}

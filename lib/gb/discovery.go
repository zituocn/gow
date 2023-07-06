/*
discovery.go

基于etcd的发现 封装
sam
2023-01-30

*/

package gb

import (
	"context"
	"errors"
	"github.com/zituocn/logx"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"sync"
	"time"
)

type Discovery struct {
	cli          *clientv3.Client
	cc           resolver.ClientConn
	etcdEndpoint []string
	values       sync.Map
	schema       string
	key          string
	prefix       string
}

func NewDiscovery(schema, key string, etcdEndpoints []string, dialTimeout int) (*Discovery, error) {
	if schema == "" {
		return nil, errors.New("need schema")
	}

	if key == "" {
		return nil, errors.New("need key")
	}

	if len(etcdEndpoints) == 0 {
		return nil, errors.New("need etcd endpoints")
	}

	if dialTimeout < 1 {
		dialTimeout = 5
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: time.Duration(dialTimeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &Discovery{
		cli:          cli,
		etcdEndpoint: etcdEndpoints,
		schema:       schema,
		key:          key,
	}, nil
}

/*
grpc discovery
*/

func (d *Discovery) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	d.cc = cc
	d.prefix = "/" + target.URL.Scheme + target.URL.Path + "/"
	resp, err := d.cli.Get(context.Background(), d.prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	for _, ev := range resp.Kvs {
		d.SetServiceList(string(ev.Key), string(ev.Value))
	}
	d.updateState()
	go d.watcher()
	return d, nil
}

func (d *Discovery) Scheme() string {
	return grpcScheme
}

func (d *Discovery) ResolveNow(options resolver.ResolveNowOptions) {
	logx.Debug("ResolveNow")
}

func (d *Discovery) Close() {
	if d.cli != nil {
		_ = d.cli.Close()
	}
}

func (d *Discovery) DelServiceList(key string) {
	d.values.Delete(key)
	d.updateState()
}

func (d *Discovery) SetServiceList(key, val string) {
	d.values.Store(key, resolver.Address{Addr: val})
	d.updateState()
}

func (d *Discovery) updateState() {
	_ = d.cc.UpdateState(resolver.State{Addresses: d.getServices()})
}

func (d *Discovery) getServices() []resolver.Address {
	addrs := make([]resolver.Address, 0, 10)
	d.values.Range(func(k, v interface{}) bool {
		addrs = append(addrs, v.(resolver.Address))
		return true
	})
	return addrs
}

func (d *Discovery) watcher() {
	rch := d.cli.Watch(context.Background(), d.prefix, clientv3.WithPrefix())
	for wResp := range rch {
		for _, ev := range wResp.Events {
			switch ev.Type {
			case mvccpb.PUT: //PUT
				d.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE: //DELETE
				d.DelServiceList(string(ev.Kv.Key))
			}
		}
	}
}

/*
discovery
*/

func (d *Discovery) Builder() error {
	d.prefix = "/" + d.schema + "/" + d.key + "/"
	resp, err := d.cli.Get(context.Background(), d.prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for _, ev := range resp.Kvs {
		d.SetDiscoveryServiceList(string(ev.Key), string(ev.Value))
	}
	go d.discoveryWatcher()
	return nil
}

// GetValues returns Values
func (d *Discovery) GetValues() []string {
	vss := make([]string, 0)
	d.values.Range(func(k, v interface{}) bool {
		vss = append(vss, v.(string))
		return true
	})
	return vss
}

func (d *Discovery) SetDiscoveryServiceList(key, val string) {
	d.values.Store(key, val)
}

func (d *Discovery) DelDiscoveryServiceList(key string) {
	d.values.Delete(key)
}

func (d *Discovery) discoveryWatcher() {
	rch := d.cli.Watch(context.Background(), d.prefix, clientv3.WithPrefix())
	for wResp := range rch {
		for _, ev := range wResp.Events {
			switch ev.Type {
			case mvccpb.PUT: //PUT
				d.SetDiscoveryServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE: //DELETE
				d.DelDiscoveryServiceList(string(ev.Kv.Key))
			}
		}
	}
}

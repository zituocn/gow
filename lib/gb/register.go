/*
register.go

基于etcd的注册 封装
sam
2023-01-30

*/

package gb

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type RegisterOption struct {

	// EtcdEndpoints  etcd hosts
	EtcdEndpoints []string

	//Lease etcd lease
	Lease int64

	// DialTimeout etcd dialTimeout second
	DialTimeout time.Duration

	// Schema  etd key prefix
	Schema string

	Key string
	Val string
}

type Register struct {
	cli           *clientv3.Client
	leaseId       clientv3.LeaseID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string
	val           string
}

func NewRegister(opt *RegisterOption) (*Register, error) {
	if opt.DialTimeout < 1 {
		opt.DialTimeout = 5
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   opt.EtcdEndpoints,
		DialTimeout: opt.DialTimeout * time.Second,
	})
	if err != nil {
		return nil, err
	}

	timeOutCtx, cancel := context.WithTimeout(context.Background(), opt.DialTimeout*time.Second/2)
	defer cancel()
	_, err = cli.Status(timeOutCtx, opt.EtcdEndpoints[0])
	if err != nil {
		return nil, fmt.Errorf("[etcd] connection timed out %s", err.Error())
	}

	reg := &Register{
		cli: cli,
		key: "/" + opt.Schema + "/" + opt.Key + "/" + opt.Val,
		val: opt.Val,
	}

	if errEtcd := reg.putKeyWithLease(opt.Lease); errEtcd != nil {
		return nil, errEtcd
	}
	reg.listenExit()
	return reg, nil
}

func (s *Register) Close() error {
	if _, err := s.cli.Revoke(context.Background(), s.leaseId); err != nil {
		return err
	}
	return s.cli.Close()
}

/*
private
*/

func (s *Register) putKeyWithLease(lease int64) error {
	resp, err := s.cli.Grant(context.Background(), lease)
	if err != nil {
		return err
	}
	_, err = s.cli.Put(context.Background(), s.key, s.val, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}
	leaseRespChan, err := s.cli.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		return err
	}
	s.leaseId = resp.ID
	s.keepAliveChan = leaseRespChan
	go s.ListenLeaseRespChan()

	return nil
}

func (s *Register) ListenLeaseRespChan() {
	for {
		<-s.keepAliveChan
	}
}

func (s *Register) listenExit() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		c := <-ch
		s.unRegister()
		if i, ok := c.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}
	}()
}

func (s *Register) unRegister() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := s.cli.Delete(ctx, s.key)
	if err != nil {
		return
	}
}

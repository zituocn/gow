package rd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	schema = "gb"
)

// ServiceOption service register options
type ServiceOption struct {
	Endpoints   []string      // etcd host
	Lease       int64         // etcd lease
	Prefix      string        // prefix,ex:serviceName
	Port        int           // grpc server port
	DialTimeout time.Duration // DialTimeout second
}

type ServiceRegister struct {
	cli           *clientv3.Client
	leaseId       clientv3.LeaseID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string
	val           string
}

// NewServiceRegister return a *ServiceRegister and error
func NewServiceRegister(opt *ServiceOption) (*ServiceRegister, error) {
	if opt.DialTimeout < 1 {
		opt.DialTimeout = 5
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   opt.Endpoints,
		DialTimeout: opt.DialTimeout * time.Second,
	})
	if err != nil {
		return nil, err
	}
	// check etcd timeout
	timeOutCtx, cancel := context.WithTimeout(context.Background(), opt.DialTimeout*time.Second/2)
	defer cancel()
	_, err = cli.Status(timeOutCtx, opt.Endpoints[0])
	if err != nil {
		return nil, fmt.Errorf("[ETCD] connection timed out %s", err.Error())
	}
	addr, err := GetLocalIP()
	if err != nil {
		return nil, err
	}
	val := fmt.Sprintf("%s:%d", addr, opt.Port)
	svc := &ServiceRegister{
		cli: cli,
		key: "/" + schema + "/" + opt.Prefix + "/" + val,
		val: val,
	}
	if err1 := svc.putKeyWithLease(opt.Lease); err1 != nil {
		return nil, err1
	}

	svc.listenExit()

	return svc, nil
}

func (s *ServiceRegister) listenExit() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		c := <-ch
		// 接收到进程退出信号量,解除租约
		s.unRegister()
		if i, ok := c.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}
	}()
}

func (s *ServiceRegister) unRegister() {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	_, err := s.cli.Delete(ctx, s.key)
	if err != nil {
		return
	}
}

func (s *ServiceRegister) putKeyWithLease(lease int64) error {
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

// ListenLeaseRespChan listen the lease chan
func (s *ServiceRegister) ListenLeaseRespChan() {
	for {
		<-s.keepAliveChan
		//fmt.Println("TTL: ", ka.TTL)
	}
}

// Close initiative close the lease
func (s *ServiceRegister) Close() error {
	if _, err := s.cli.Revoke(context.Background(), s.leaseId); err != nil {
		return err
	}
	return s.cli.Close()
}

package rpc

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"go.etcd.io/etcd/client/v3"
)

var (
	etcdConn *EtcdCli
	once     sync.Once
	randEr   = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// EtcdCli etcd cli struct
type EtcdCli struct {
	cli         *clientv3.Client
	EtcdAddr    []string
	ttl         int // 申请租约的时间,单位秒, ttl秒后就会自动移除
	connTimeOut int // 连接etcd的timeout
}

// NewEtcdCli new etcd cli
func NewEtcdCli(etcdAddr []string) *EtcdCli {
	once.Do(
		func() {
			etcdConn = &EtcdCli{
				EtcdAddr:    etcdAddr,
				ttl:         2,
				connTimeOut: 3,
			}
		},
	)
	return etcdConn
}

// SetTTl set ttl
func (etcdCil *EtcdCli) SetTTl(ttl int) *EtcdCli {
	etcdCil.ttl = ttl
	return etcdCil
}

// SetConnTimeOut set etcd conn timeout
func (etcdCil *EtcdCli) SetConnTimeOut(timeOut int) *EtcdCli {
	etcdCil.connTimeOut = timeOut
	return etcdCil
}

// Conn return *EtcdCli error
func (etcdCil *EtcdCli) Conn() (*EtcdCli, error) {
	err := etcdCil.clientv3New()
	return etcdCil, err
}

// clientv3New returns  error
func (etcdCil *EtcdCli) clientv3New() (err error) {
	if len(etcdCil.EtcdAddr) < 1 {
		err = fmt.Errorf("etcd addr is null")
		return
	}
	if etcdCil.cli == nil {
		etcdCil.cli, err = clientv3.New(clientv3.Config{
			Endpoints:   etcdCil.EtcdAddr,
			DialTimeout: time.Duration(etcdCil.connTimeOut) * time.Second,
		})
	}
	return
}

// Register 注册,并创建租约
func (etcdCil *EtcdCli) Register(key, value string) error {
	if err := etcdCil.clientv3New(); err != nil {
		return err
	}
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for {
			getResp, err := etcdCil.cli.Get(context.Background(), key)
			if err != nil {
				log.Printf("[ETCD] Register err : %s", err)
			} else if getResp.Count == 0 {
				err = etcdCil.withAlive(key, value)
				if err != nil {
					log.Printf("[ETCD] keep alive err :%s", err)
				}
			}
			<-ticker.C
		}
	}()
	return nil
}

// withAlive 创建租约
func (etcdCil *EtcdCli) withAlive(key, value string) error {
	leaseResp, err := etcdCil.cli.Grant(context.Background(), int64(etcdCil.ttl))
	if err != nil {
		return err
	}
	_, err = etcdCil.cli.Put(context.Background(), key, value, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		log.Printf("[ETCD] put etcd error:%s", err)
		return err
	}

	ch, err := etcdCil.cli.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		log.Printf("[ETCD] keep alive error:%s", err)
		return err
	}

	// 清空 keep alive 返回的channel
	go func() {
		for {
			<-ch
		}
	}()

	return nil
}

// UnRegister 解除注册
func (etcdCil *EtcdCli) UnRegister(key string) error {
	if err := etcdCil.clientv3New(); err != nil {
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	etcdCil.cli.Delete(ctx, key)
	return nil
}

// Get return *clientv3.GetResponse error
func (etcdCil *EtcdCli) Get(key string) (*clientv3.GetResponse, error) {
	if err := etcdCil.clientv3New(); err != nil {
		return nil, err
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(etcdCil.connTimeOut)*time.Second)
	defer cancelFunc()
	return etcdCil.cli.Get(ctx, key, clientv3.WithPrefix())
}

// GetMinKey 获取被计数最少的key
func (etcdCil *EtcdCli) GetMinKey(key string) (string, error) {
	response, err := etcdCil.Get(key)
	if err != nil {
		return "", err
	}

	if len(response.Kvs) < 1 {
		return "", fmt.Errorf("[ETCD] 在注册表没有找到服务")
	}

	tmp := response.Kvs[0].Key
	tmpValue := response.Kvs[0].Value
	for i := 1; i < len(response.Kvs); i++ {
		if byte2int(tmpValue) > byte2int(response.Kvs[i].Value) {
			tmp = response.Kvs[i].Key
			tmpValue = response.Kvs[i].Value
		}
	}
	return string(tmp), nil
}

// GetMinKeyCallBack 是GetMinKey方法的回调
func (etcdCil *EtcdCli) GetMinKeyCallBack(key string) error {
	if err := etcdCil.clientv3New(); err != nil {
		return err
	}

	getResp, err := etcdCil.cli.Get(context.Background(), key)
	if err != nil {
		return err
	}

	if getResp.Count > 0 {
		v := getResp.Kvs[0].Value
		vInt := byte2int(v) + 1
		_, err = etcdCil.cli.Put(context.Background(), key, fmt.Sprintf("%d", vInt), clientv3.WithPrevKV())
		if err != nil {
			return err
		}
	}
	return nil
}

// GetAllKey 使用前缀key查询值
func (etcdCil *EtcdCli) GetAllKey(key string) ([]string, error) {
	keys := make([]string, 0)
	response, err := etcdCil.Get(key)
	if err != nil {
		return keys, err
	}

	if len(response.Kvs) < 1 {
		return keys, fmt.Errorf("[ETCD] 在注册表没有找到服务")
	}

	for _, ev := range response.Kvs {
		grpcIPList := strings.Split(string(ev.Key), "/")
		grpcIP := grpcIPList[len(grpcIPList)-1]
		keys = append(keys, grpcIP)
	}
	return keys, nil
}

//GetRandKey 使用前缀key查询随机反回一个
func (etcdCil *EtcdCli) GetRandKey(key string) (string, error) {
	keys, err := etcdCil.GetAllKey(key)
	if err != nil {
		return "", err
	}
	l := len(keys)
	if l < 1 {
		return "", fmt.Errorf("[ETCD] 在注册表没有找到服务")
	}
	if l == 1 {
		return keys[0], nil
	}
	r := randEr.Intn(l)
	return keys[r], nil
}

// GetHash 通过传入h把value进行hash,然后后返回
// h : 用于计算hash的值
func (etcdCil *EtcdCli) GetHash(key, h string) {

}

// Delete delete key
//	return error
func (etcdCil *EtcdCli) Delete(key string) error {
	if err := etcdCil.clientv3New(); err != nil {
		return err
	}
	_, err := etcdCil.cli.Delete(context.TODO(), key, clientv3.WithPrevKV())
	return err
}

// DeleteAll delete all key by Prefix
//	returns error
func (etcdCil *EtcdCli) DeleteAll(key string) error {
	if err := etcdCil.clientv3New(); err != nil {
		return err
	}
	_, err := etcdCil.cli.Delete(context.TODO(), key, clientv3.WithPrefix())
	return err
}

# gb

grpc的封装库


### 安装

```sh
go get github.com/zituocn/gow/lib/gb
```

### 特性

* 可使用普通的 grpc server 与 client；
* 可实现 grpc 的负载均衡，需要配合 etcd 使用；


### 服务端使用方法

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/zituocn/gb_demo/pb"
	"github.com/zituocn/gow/lib/gb"
)

var (
	etcdEndPoints = []string{"127.0.0.1:2379"}
	ServiceName   = "grpc-demo"
	Port          = 3333
)

func init() {
	initGrpc()
}

func main() {
	select {}
}

func initGrpc() {
	s, err := gb.NewServer(&gb.ServerOption{
		Name:          ServiceName,   //服务名
		Port:          Port,          //端口
		EtcdEndPoints: etcdEndPoints, //etcd连接配置
		Lease:         3,             //etcd lease时间
		IsRegister:    true,          //是否注册到etcd
		KeyFile:       "",            //私钥
		CertFile:      "",            //公钥
	})

	if err != nil {
		log.Fatalf("get grpc server error :%s", err.Error())
	}
	pb.RegisterUserServer(s.Server, &UserService{})
	s.Run()
}

type UserService struct {
}

func (s *UserService) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	fmt.Println(fmt.Sprintf("from %d", req.Uid))
	return &pb.UserResponse{
		Uid:      req.Uid,
		Username: "username",
		Phone:    "13888889999",
		Address:  "chengdu in china",
	}, nil
}

```

### 客户端代码

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/zituocn/gb_demo/pb"
	"github.com/zituocn/gow/lib/gb"
)

var (
	etcdEndPoints = []string{"127.0.0.1:2379"}
	ServiceName   = "grpc-demo"
	ServerAddr    = "127.0.0.1:1111"
)

func main() {
	
	// ConnType 可以为 
	//  gb.BalanceType 负载均衡方式，可以不指定服务器地址
	//  gb.DefaultType 非负载均衡方式，需要指定服务端的连接地址
	client, err := gb.NewClient(&gb.ClientOption{
		ServerAddr:   ServerAddr,     //服务器地址 host:port
		ServiceName:  ServiceName,    //服务名
		ClientName:   "demo-test",    //客户端名称
		EtcdEndPoint: etcdEndPoints,  //etcd节点
		ConnType:     gb.BalanceType, //连接类型
		CertFile:     "",             //公钥
	})
	if err != nil {
		log.Fatalf("%v", err)
	}

	conn, ctx, err := client.GetConn()
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer conn.Close()

	// grpc调用
	userClient := pb.NewUserClient(conn)
	for i := 0; i < 10000; i++ {
		req := &pb.UserRequest{
			Uid: int32(i),
		}
		resp, err := userClient.GetUser(ctx, req)
		if err != nil {
			log.Printf("Call grpc error: %v \n", err)
		}
		fmt.Println(resp)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		log.Fatalf("net connect  error :%v", err)
	}
}
```
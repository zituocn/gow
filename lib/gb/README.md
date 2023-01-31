# gb

基于etcd的注册与发现


### 注册

```go
package main

import (
	"flag"
	"github.com/zituocn/srd"
	"log"
)

func main() {
	flag.Parse()
	_, err := srd.NewRegister(&srd.RegisterOption{
		EtcdEndpoints: []string{"192.168.0.101:2379"},
		Lease:         3,
		Schema:        "gk100-cache",
		Key:           "pod-ip",
		Val:           "10.10.10.2:10003",
	})
	if err != nil {
		log.Fatal(err)
	}
	select {}
}

```

### 发现

```go
package main

import (
	"fmt"
	"github.com/zituocn/srd"
	"log"
)

func main() {
	dis, err := srd.NewDiscovery("gk100-cache", "pod-ip", []string{"192.168.0.101:2379"}, 3)
	if err != nil {
		log.Fatal(err)
	}
	err = dis.Builder()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dis.GetValues())
}

```


### GRPC服务端

```go
package main

import (
	"context"
	"github.com/zituocn/gow/lib/logy"
	"github.com/zituocn/srd"
	"github.com/zituocn/srd/demo/pb"
	"log"
)

var (
	etcdEndPoints = []string{"192.168.0.101:2379"}
	ServerName    = "user-service"
	Port          = 1111
)

func init() {
	gs, err := srd.NewGrpcServer(&srd.GrpcServerOption{
		ServerName:    ServerName,    //服务名
		Port:          Port,          //端口
		EtcdEndpoints: etcdEndPoints, //etcd连接配置
		Lease:         5,             //etcd lease时间
		IsRegister:    true,          //是否注册到etcd
	})

	if err != nil {
		log.Fatal(err)
	}

	pb.RegisterUserServer(gs.Server, &UserService{})

	gs.Run()
}

func main() {
	select {}
}

type UserService struct {
}

func (s *UserService) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	logy.Debugf("from %d", req.Uid)
	return &pb.UserResponse{
		Uid:      req.Uid,
		Username: "zituocn",
		Phone:    "13888889999",
		Address:  "chengdu in china",
	}, nil
}

```

### GRPC客户端

```go
package main

import (
	"github.com/zituocn/gow"
	"github.com/zituocn/srd"
	"github.com/zituocn/srd/demo/pb"
	"log"
)

var (
	etcdEndPoints = []string{"127.0.0.1:2379"}
	serverName    = "user-service"
)

func main() {
	r := gow.Default()
	r.GET("/user/{id}", GetUser)
	r.Run()
}

func GetUser(c *gow.Context) {
	id, _ := c.ParamInt("id")
	client, err := srd.NewGrpcClient(&srd.GrpcClientOption{
		ServerName:   serverName,
		EtcdEndpoint: etcdEndPoints,
		ClientName:   "test",
		ConnType:     srd.BalanceType,
	})

	if err != nil {
		log.Fatal(err)
	}

	conn, ctx, err := client.GetConn()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	userClient := pb.NewUserClient(conn)

	req := &pb.UserRequest{
		Uid: int32(id),
	}
	resp, err := userClient.GetUser(ctx, req)
	if err != nil {
		c.DataJSON(1, err.Error())
		return
	}
	c.DataJSON(resp)
}


```

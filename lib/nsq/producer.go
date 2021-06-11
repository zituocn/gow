/*
	 pu,err:=NewProducer("192.168.0.197：4150")
	 if err!=nil{
		//handler error
	 }
	 b,_:=json.Marshal(obj)

  	 //publish
	 err = pu.Publish("topic",b)
	 if err!=nil{
		//handler error
	 }

*/

package nsq

import (
	"fmt"
	gnsq "github.com/nsqio/go-nsq"
)

// Producer Producer
type Producer struct {
	P *gnsq.Producer
}

// NewProducer init
//	return producer
func NewProducer(addr string) (producer *Producer, err error) {
	if addr == "" {
		err = fmt.Errorf("[NSQ] init failed：need nsq server addr")
		return
	}
	config := gnsq.NewConfig()
	p, err := gnsq.NewProducer(addr, config)
	if err != nil {
		return
	}
	p.SetLogger(nil, 0)

	producer = &Producer{
		P: p,
	}
	return
}

// Publish publish topic
//	return error
func (m *Producer) Publish(topic string, data []byte) (err error) {
	if m.P == nil {
		err = fmt.Errorf("[NSQ] init failed:%v", err)
	}
	err = m.P.Publish(topic, data)
	defer m.P.Stop()
	if err != nil {
		return fmt.Errorf("[NSQ] publish error:%v", err)
	}
	return
}

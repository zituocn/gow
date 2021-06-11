package nsq

import (
	"fmt"
	gnsq "github.com/nsqio/go-nsq"
)

//MessageHandler MessageHandler
type MessageHandler struct {
	msgChan   chan *gnsq.Message
	stop      bool
	nsqServer string
	Channel   string
}

// NewMessageHandler return new MessageHandler
func NewMessageHandler(nsqServer string, channel string) (mh *MessageHandler, err error) {
	if nsqServer == "" {
		err = fmt.Errorf("[NSQ] need nsq server")
		return
	}
	mh = &MessageHandler{
		msgChan:   make(chan *gnsq.Message, 1024),
		stop:      false,
		nsqServer: nsqServer,
		Channel:   channel,
	}

	return
}

// Registry registry nsq topic
func (m *MessageHandler) Registry(topic string, ch chan []byte) {
	config := gnsq.NewConfig()
	consumer, err := gnsq.NewConsumer(topic, m.Channel, config)
	if err != nil {
		panic(err)
	}
	consumer.SetLogger(nil, 0)
	consumer.AddHandler(gnsq.HandlerFunc(m.handlerMessage))
	err = consumer.ConnectToNSQLookupd(m.nsqServer)
	if err != nil {
		panic(err)
	}
	m.process(ch)

}

// process process
func (m *MessageHandler) process(ch chan<- []byte) {
	m.stop = false
	for {
		select {
		case message := <-m.msgChan:
			ch <- message.Body
			if m.stop {
				close(m.msgChan)
				return
			}
		}
	}
}

// handlerMessage handlerMessage
func (m *MessageHandler) handlerMessage(message *gnsq.Message) error {
	if !m.stop {
		m.msgChan <- message
	}
	return nil
}

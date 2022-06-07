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
	maxInFlight int
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

// SetMaxInFlight set nsq consumer MaxInFlight
func (m *MessageHandler) SetMaxInFlight(val int){
	m.maxInFlight = val
}

// Registry register nsq topic
func (m *MessageHandler) Registry(topic string, ch chan []byte) {
	config := gnsq.NewConfig()
	if m.maxInFlight > 0 {
		config.MaxInFlight = m.maxInFlight
	}
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
	message.Finish()
	return nil
}

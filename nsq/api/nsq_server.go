package api

import (
	"fmt"
	"github.com/ifchange/botKit/nsq/consumer"
	"github.com/ifchange/botKit/nsq/message"
	"github.com/ifchange/botKit/nsq/publisher"
)

var (
	lookupdHost = "127.0.0.1:4161"
	nsqHost     = "localhost:4150"
	topicName   = "test"
	channelName = "channel1"
)

type NsqServerI interface {
	Publish(msg message.Message) error
	AddHandler()
}

type NsqServer struct {
	pi publisher.PublisherI
	ci consumer.ConsumerI
}

//todo 需要拆分为独立服务
func InitStdNsq() (NsqServerI, error) {
	pinew, err := publisher.NewPublisher(nsqHost)
	if err != nil {
		panic(fmt.Errorf("init config error, newPublisher error, %v", err))
	}
	cinew, err := consumer.NewConsumer(topicName, channelName)
	if err != nil {
		panic(fmt.Errorf("init config error, newConsumer error, %v", err))
	}
	ns := &NsqServer{
		pi: pinew,
		ci: cinew,
	}
	return ns, nil
}

func (n *NsqServer) Publish(msg message.Message) error {
	err := n.pi.Publish(msg.TopicName, msg.Msg)
	if err != nil {
		return err
	}
	return nil
}

func (n *NsqServer) AddHandler() {}

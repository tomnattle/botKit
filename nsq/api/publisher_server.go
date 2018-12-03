package api

import (
	"fmt"

	"github.com/ifchange/botKit/nsq/message"
	"github.com/ifchange/botKit/nsq/publisher"
)

type PublisherServer struct {
	pi publisher.PublisherI
}

type PublisherP struct {
	nsqHost string
}

func NewPublisher(pp *PublisherP) (PublisherServerI) {
	pinew, err := publisher.NewPublisher(pp.nsqHost)
	if err != nil {
		panic(fmt.Errorf("NewPublisher error, %v", err))
	}
	ps := &PublisherServer{
		pi: pinew,
	}
	return ps
}

func (ps *PublisherServer) Publish(msg *message.Message) error {
	err := ps.pi.Publish(msg.TopicName, msg.Msg)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PublisherServer) Stop() {
	ps.pi.Stop()
}

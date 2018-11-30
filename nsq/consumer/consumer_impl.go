package consumer

import (
	"github.com/nsqio/go-nsq"
)

type ConsumerImpl struct {
	consumer *nsq.Consumer
}

func NewConsumer(topic string, channel string) (ConsumerI, error) {
	config := nsq.NewConfig()
	c, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		return nil, err
	}
	return &ConsumerImpl{consumer: c,}, nil
}

func (c *ConsumerImpl) ConnectToNSQLookupd(lookupHosts []string) error {
	if err := c.consumer.ConnectToNSQLookupds(lookupHosts); err != nil {
		return err
	}
	return nil
}

func (c *ConsumerImpl) AddHandler(handler nsq.Handler) {
	c.consumer.AddHandler(handler)
}

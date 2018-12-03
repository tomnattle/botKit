package publisher

import (
	"github.com/nsqio/go-nsq"
)

type PublisherImpl struct {
	producer *nsq.Producer
}

func NewPublisher(address string) (*PublisherImpl, error) {
	config := nsq.NewConfig()
	p, err := nsq.NewProducer(address, config)
	if err != nil {
		return nil, err
	}
	return &PublisherImpl{producer: p}, nil
}

func (p *PublisherImpl) Publish(topic string, message []byte) error {
	err := p.producer.Publish(topic, message)
	if err != nil {
		return err
	}
	return nil
}

func (p *PublisherImpl) MultiPublish(topic string, messages [][]byte) error {
	err := p.producer.MultiPublish(topic, messages)
	if err != nil {
		return err
	}
	return nil
}

func (p *PublisherImpl) Stop() {
	if p.producer != nil {
		p.producer.Stop()
	}
}

package api

import (
	"github.com/nsqio/go-nsq"

	"github.com/ifchange/botKit/nsq/message"
)

type PublisherServerI interface {
	Publish(msg *message.Message) error
	Stop()
}

type ConsumerServerI interface {
	StartHandler(handler nsq.Handler)
	StopHandler()
}

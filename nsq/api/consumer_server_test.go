package api

import (
	"log"
	"testing"
	"time"

	"github.com/ifchange/botKit/nsq/mock"
	"github.com/nsqio/go-nsq"
)

func TestConsumerServer_StartHandler(t *testing.T) {
	nc := NewConsumer(&ConsumerP{
		topicName:    mock.TopicName,
		channelName:  mock.ChannelName,
		lookupdHosts: []string{mock.LookupdHost},
	})

	go func() {
		time.Sleep(10 * time.Second)
		nc.StopHandler()
	}()

	nc.StartHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Printf("Got a message: %v", string(message.Body))
		return nil
	}))

	t.Log("test consumer exit")
}

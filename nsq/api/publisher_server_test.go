package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/ifchange/botKit/nsq/message"
	"github.com/ifchange/botKit/nsq/mock"
)

func TestPublisherServer_Publish(t *testing.T) {
	pp := &PublisherP{
		nsqHost: mock.NsqHost,
	}
	np := NewPublisher(pp)

	quit := false
	go func() {
		time.Sleep(20 * time.Minute)
		np.Stop()
		quit = true
	}()

	for {
		np.Publish(&message.Message{
			TopicName: mock.TopicName,
			Msg:       []byte("hello"),
		})

		time.Sleep(1 * time.Second)

		if quit {
			break
		}
	}

	t.Log("test publisher exit")
}

func TestPublisherServer_MultiPublish(t *testing.T) {
	muti := [][]byte{[]byte("{\"json\"}"), []byte("aaa")}
	fmt.Println(muti)
}

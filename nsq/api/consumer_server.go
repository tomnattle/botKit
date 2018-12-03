package api

import (
	"fmt"
	"log"

	"github.com/nsqio/go-nsq"

	"github.com/ifchange/botKit/nsq/consumer"
)

type ConsumerServer struct {
	ci         consumer.ConsumerI
	quitCh     chan string
	lookupHost []string
}

type ConsumerP struct {
	topicName    string
	channelName  string
	lookupdHosts []string
}

func NewConsumer(cp *ConsumerP) (ConsumerServerI) {
	cinew, err := consumer.NewConsumer(cp.topicName, cp.channelName)
	if err != nil {
		log.Fatal("new consumer error ", err.Error())
	}
	cs := &ConsumerServer{
		ci:         cinew,
		lookupHost: cp.lookupdHosts,
	}
	return cs
}

func (cs *ConsumerServer) StartHandler(handler nsq.Handler) {
	cs.ci.AddHandler(handler)
	err := cs.ci.ConnectToNSQLookupd(cs.lookupHost)
	fmt.Println("ch start")
	if err != nil {
		log.Fatal("connect to nsq lookup err ", err.Error())
	}
	cs.quitCh = make(chan string)
	<-cs.quitCh
}

func (cs *ConsumerServer) StopHandler() {
	if cs.quitCh != nil {
		cs.quitCh <- "quit"
	}
}

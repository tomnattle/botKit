package consumer

import "github.com/nsqio/go-nsq"

type ConsumerI interface {
	AddHandler(handler nsq.Handler)
	ConnectToNSQLookupd(lookupHosts []string) error
}

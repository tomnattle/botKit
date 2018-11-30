package consumer

import "github.com/nsqio/go-nsq"

type ConsumerI interface {
	ConnectToNSQLookupd(lookupHosts []string) error
	AddHandler(handler nsq.Handler)
}

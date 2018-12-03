package publisher

type PublisherI interface {
	Publish(topic string, message []byte) error
	MultiPublish(topic string, messages [][]byte) error
	Stop()
}

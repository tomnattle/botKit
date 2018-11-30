package publisher

type PublisherI interface {
	Publish(topic string, message []byte) error
}
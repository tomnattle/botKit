package message

type TopicName struct{}

type Message struct {
	TopicName string
	Msg       []byte
}

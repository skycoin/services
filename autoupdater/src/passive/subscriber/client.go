package subscriber

type Subscriber interface {
	Subscribe(topic string)
}

func NewNatsSubscriber(topic string) Subscriber {
	return newNats(topic)
}
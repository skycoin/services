package subscriber

type Subscriber interface {
	Subscribe(topic string)
}

func NewNats(topic string) Subscriber {
	return newNats(topic)
}
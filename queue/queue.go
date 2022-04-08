package queue

type Queue interface {
	AddTopic(name string) Topic
	TopicByName(name string) (Topic, bool)
	Topics() []Topic

	Producer
	Consumer
}

type Topic interface {
	Name() string
}

type Message interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

type Receipt interface {
	String() string
}

type Producer interface {
	Produce(topic Topic, msg Message) (Receipt, error)
}

type Consumer interface {
	Consume(topic Topic, f func(msg Message) error) error
}

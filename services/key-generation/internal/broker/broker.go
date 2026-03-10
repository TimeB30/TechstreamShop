package broker

type Broker interface {
	Produce(message *KeyMessage) error
	Consume(timeout int) (*OrderMessage, error)
	Close() (int, error)
}

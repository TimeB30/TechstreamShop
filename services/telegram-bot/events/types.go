package events

type Fetcher interface {
	Fetch(limit int64) ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}
type Broker interface {
	Produce(topic string, payload any) error
	Consume() (*Event, error)
	Close() (int, error)
}
type Type int64

const (
	Unknown Type = iota
	Message
	Key
)

type Event struct {
	Type Type
	Text string
	Meta interface{}
}
type MessageMeta struct {
	UserID   int64
	ChatID   int64
	UserName string
}

type KeyMeta struct {
	UserID int64
	ChatID int64
}

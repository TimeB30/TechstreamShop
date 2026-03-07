package events

type Fetcher interface {
	Fetch(limit int64) ([]Event, error)
}
type Listener interface {
	Listen()
}

type Processor interface {
	Process(e Event) error
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

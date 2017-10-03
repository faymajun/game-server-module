package event

type Event interface {
	Source() interface{}
	Type() EventType
}

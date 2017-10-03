package event

type Listener interface {
	OnEvent(Event) bool
}

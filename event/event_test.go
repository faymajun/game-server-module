package event

import (
	"fmt"
	"testing"
)

type FuncListener struct {
	Callback func(Event) bool
}

func (fl *FuncListener) OnEvent(event Event) bool {
	return fl.Callback(event)
}

type LoggingListener struct {
}

func (ll *LoggingListener) OnEvent(e Event) bool {
	fmt.Println(e)
	return false
}

type eventImpl struct {
	eventType EventType
	source    interface{}
}

func (e *eventImpl) String() string {
	return fmt.Sprintf(" type=%v", e.eventType)
}

func (e *eventImpl) Source() interface{} {
	return e.source
}

func (e *eventImpl) Type() EventType {
	return e.eventType
}

func TestNewEventDispatcher(t *testing.T) {
	callbackCalled := false

	ed := NewEventDispatcher()
	ed.AddEventListener(EventNone, &FuncListener{Callback: func(e Event) bool { callbackCalled = true; return false }})
	ed.Dispatch(&eventImpl{EventNone, nil})

	fmt.Println(callbackCalled)

	ed.AddEventListener(EventNone, &LoggingListener{})
	ed.Dispatch(&eventImpl{EventNone, nil})
}

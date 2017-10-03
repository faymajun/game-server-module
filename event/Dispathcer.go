package event

import "fmt"

type EventDispatcher struct {
	eventListeners map[EventType][]Listener
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		eventListeners: make(map[EventType][]Listener),
	}
}

func (ed *EventDispatcher) Dispatch(event Event) {
	list, ok := ed.eventListeners[event.Type()]

	if ok {
		for i := range list {
			list[i].OnEvent(event)
		}
	}
}

func (ed *EventDispatcher) AddEventListener(eventType EventType, listener Listener) {
	_, ok := ed.eventListeners[eventType]
	if !ok {
		ed.eventListeners[eventType] = make([]Listener, 0)
	}

	ed.eventListeners[eventType] = append(ed.eventListeners[eventType], listener)
}

func (ed *EventDispatcher) RemoveEventListener(eventType EventType, listener Listener) {
	list, ok := ed.eventListeners[eventType]

	if ok {
		for i := range list {
			if list[i] == listener {
				list = append(list[:i], list[i+i:]...)
				break
			}
		}
	} else {
		fmt.Printf("CEventDispatcher：RemoveEventListener：当前事件id不包含任何回调结构体。eventType=%v", eventType)
	}
}

func (ed *EventDispatcher) RemoveListener(listener Listener) {
	for eventType := range ed.eventListeners {
		ed.RemoveEventListener(eventType, listener)
	}
}

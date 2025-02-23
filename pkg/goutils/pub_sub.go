package goutils

type EventPublisher[T any] struct {
	subscribers []chan T
}

func NewEventPublisher[T any]() *EventPublisher[T] {
	return &EventPublisher[T]{
		subscribers: make([]chan T, 0),
	}
}

func (e *EventPublisher[T]) Publish(event T) {
	for _, ch := range e.subscribers {
		ch <- event
	}
}

func (e *EventPublisher[T]) Subscribe() <-chan T {
	ch := make(chan T)
	e.subscribers = append(e.subscribers, ch)
	return ch
}

func (e *EventPublisher[T]) Close() {
	for _, ch := range e.subscribers {
		close(ch)
	}
}

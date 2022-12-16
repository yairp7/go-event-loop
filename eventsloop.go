package eventsloop

import (
	"context"
	"errors"
)

type EventLoop[T any] struct {
	Opts      EventLoopOptions
	Queue     chan T
	Finished  bool
	Processor EventProcessor[T]
}

type EventLoopOptions struct {
	MaxQueueSize int
}

type EventProcessor[T any] interface {
	ProcessEvent(event T)
}

func NewEventLoop[T any](
	opts EventLoopOptions,
	processor EventProcessor[T],
) EventLoop[T] {
	return EventLoop[T]{
		Opts:      opts,
		Queue:     make(chan T, opts.MaxQueueSize),
		Finished:  false,
		Processor: processor,
	}
}

func (e *EventLoop[T]) AddEvent(event T) error {
	if e.Finished {
		return errors.New(`queue is closed`)
	}

	e.Queue <- event
	return nil
}

func (e *EventLoop[T]) Run(ctx context.Context) {
	go func() {
		for {
			if e.Finished {
				return
			}

			select {
			case <-ctx.Done():
				e.Close()
				return
			case event := <-e.Queue:
				e.Process(event)
			}
		}
	}()
}

func (e *EventLoop[T]) Close() {
	e.Finished = true
	close(e.Queue)
}

func (e *EventLoop[T]) Process(event T) {
	e.Processor.ProcessEvent(event)
}

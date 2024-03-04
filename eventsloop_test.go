package eventsloop_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	eventsloop "github.com/yairp7/go-event-loop"
)

type EventsLoopTestProcessor struct {
	events []string
}

var strings = []string{
	"dsgdsg",
	"grgegf",
	"sdg23",
	"dsg43t",
	"sdfgh45hgf",
	"sdfad32",
	"fdg5rd",
	"fh456y543g",
	"reg34gfw",
	"dg342r",
	"fd435y54",
	"dsfe56",
	"dsf325",
}

func newEventsLoopTestProcessor() EventsLoopTestProcessor {
	return EventsLoopTestProcessor{}
}

func (p *EventsLoopTestProcessor) ProcessEvent(event string) {
	p.events = append(p.events, event)
}

func TestAddEvents(t *testing.T) {
	processor := newEventsLoopTestProcessor()
	eventLoop := eventsloop.NewEventLoop[string](eventsloop.EventLoopOptions{
		MaxQueueSize: 10,
	}, &processor)

	go func() {
		for _, str := range strings {
			eventLoop.AddEvent(str)
			log.Println(str)
		}
	}()

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(1*float64(time.Second)))

	eventLoop.Run(ctx)

	<-ctx.Done()

	assert.Len(t, processor.events, len(strings))

	i := 0
	for _, event := range processor.events {
		if event != strings[i] {
			t.FailNow()
		}
		i++
	}
}

func TestTimeout(t *testing.T) {
	processor := newEventsLoopTestProcessor()
	eventLoop := eventsloop.NewEventLoop[string](eventsloop.EventLoopOptions{
		MaxQueueSize: 10,
	}, &processor)

	go func() {
		for _, str := range strings {
			eventLoop.AddEvent(str)
			time.Sleep(time.Duration(100 * float64(time.Millisecond)))
		}
	}()

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(500*float64(time.Millisecond)))

	eventLoop.Run(ctx)

	<-ctx.Done()

	assert.Len(t, processor.events, 5)
}

func TestClose(t *testing.T) {
	processor := newEventsLoopTestProcessor()
	eventLoop := eventsloop.NewEventLoop[string](eventsloop.EventLoopOptions{
		MaxQueueSize: 10,
	}, &processor)

	go func() {
		for _, str := range strings {
			eventLoop.AddEvent(str)
			time.Sleep(time.Duration(100 * float64(time.Millisecond)))
		}
	}()

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(1*float64(time.Second)))

	eventLoop.Run(ctx)

	go func() {
		time.Sleep(time.Duration(500 * float64(time.Millisecond)))
		eventLoop.Close()
	}()

	<-ctx.Done()

	assert.Less(t, len(processor.events), 10)
}

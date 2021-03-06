package eventbus

import (
	"context"
	"sync"
)

type Event = interface{}

type EventHandler func(context.Context, Event)

type EventBus struct {
	eventHandlerLists map[string][]EventHandler
	lock              sync.RWMutex
}

var eventBus = &EventBus{}

func Subscribe(eventName string, eventHandlers ...EventHandler) {
	eventBus.Subscribe(eventName, eventHandlers...)
}

func Publish(ctx context.Context, eventName string, event Event) {
	eventBus.Publish(ctx, eventName, event)
}

func (b *EventBus) Subscribe(eventName string, eventHandlers ...EventHandler) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.eventHandlerLists == nil {
		b.eventHandlerLists = make(map[string][]EventHandler)
	}
	eventHandlerList, ok := b.eventHandlerLists[eventName]
	if !ok {
		eventHandlerList = make([]EventHandler, 0, len(eventHandlers))
	}
	b.eventHandlerLists[eventName] = append(eventHandlerList, eventHandlers...)
}

func (b *EventBus) Publish(ctx context.Context, eventName string, event Event) {
	var handlerList []EventHandler
	b.lock.RLock()
	if b.eventHandlerLists != nil {
		handlerList = b.eventHandlerLists[eventName]
	}
	b.lock.RUnlock()
	for _, handler := range handlerList {
		handler(ctx, event)
	}
}

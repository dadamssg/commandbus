package CommandBus

import (
	"reflect"
	"sort"
	"sync"
)

type HandlerFunc func(cmd interface{})

type middlewareFunc func(cmd interface{}, next HandlerFunc)

type middleware struct {
	function middlewareFunc
	priority int
}

type byPriority []middleware

func (slice byPriority) Len() int {
	return len(slice)
}

func (slice byPriority) Less(i, j int) bool {
	return slice[i].priority > slice[j].priority
}

func (slice byPriority) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type commandBus struct {
	handlers    map[reflect.Type]HandlerFunc
	middlewares []middleware
	lock        sync.Mutex
}

func (bus *commandBus) RegisterHandler(cmd interface{}, handler HandlerFunc) {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	bus.handlers[reflect.TypeOf(cmd)] = handler
}

func (bus *commandBus) AddMiddleware(priority int, function middlewareFunc) {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	bus.middlewares = append(bus.middlewares, middleware{function: function, priority: priority})
	sort.Sort(byPriority(bus.middlewares))
}

func (bus commandBus) Handle(cmd interface{}) {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	handler := bus.getNext(0)
	handler(cmd)
}

func (bus commandBus) getNext(index int) HandlerFunc {
	if len(bus.middlewares) >= (index + 1) {
		middleware := bus.middlewares[index]
		return func(cmd interface{}) {
			middleware.function(cmd, bus.getNext(index+1))
		}
	}

	return func(cmd interface{}) {
		handler := bus.getHandler(cmd)
		if handler != nil {
			handler(cmd)
		}
	}
}

func (bus commandBus) getHandler(cmd interface{}) HandlerFunc {
	t := reflect.TypeOf(cmd)

	for kind, handler := range bus.handlers {
		if t == kind {
			return handler
		}
	}

	return nil
}

func New() *commandBus {
	return &commandBus{
		handlers:    make(map[reflect.Type]HandlerFunc),
		middlewares: make([]middleware, 0),
	}
}

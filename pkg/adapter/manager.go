package adapter

import (
	"fmt"
	"sync"

	"yora/pkg/middleware"
)

type AdapterRegistry struct {
	mu          sync.RWMutex
	adapters    map[Protocol]Adapter
	middlewares []middleware.Middleware
}

// Adapters implements Registry.
func (r *AdapterRegistry) Adapters() map[Protocol]Adapter {
	return r.adapters
}

// Middlewares implements Registry.
// func (r *AdapterRegistry) Middlewares() []middleware.Middleware {
// 	return r.middlewares
// }

func NewAdapterRegistry() *AdapterRegistry {
	return &AdapterRegistry{
		adapters:    make(map[Protocol]Adapter),
		middlewares: make([]middleware.Middleware, 0),
	}
}

func (r *AdapterRegistry) Register(adapter Adapter) error {
	if adapter == nil {
		return fmt.Errorf("adapter cannot be nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	proto := adapter.Protocol()
	if _, exists := r.adapters[proto]; exists {
		return fmt.Errorf("adapter for protocol %s already registered", proto)
	}

	r.adapters[proto] = adapter
	return nil
}

func (r *AdapterRegistry) Unregister(protocol Protocol) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.adapters[protocol]; !exists {
		return fmt.Errorf("adapter for protocol %s not found", protocol)
	}
	delete(r.adapters, protocol)
	return nil
}

func (a *AdapterRegistry) RegisterMiddleware(middleware middleware.Middleware) {
	a.middlewares = append(a.middlewares, middleware)
}

package onebot

import (
	"fmt"
	"sync"
	"yora/internal/adapter"
)

var _ adapter.Registry = (*AdapterRegistry)(nil)

type AdapterRegistry struct {
	mu          sync.RWMutex
	adapters    map[adapter.Protocol]adapter.Adapter
	middlewares []adapter.Middleware
}

// Adapters implements adapter.Registry.
func (r *AdapterRegistry) Adapters() []adapter.Adapter {
	adapters := make([]adapter.Adapter, 0, len(r.adapters))
	for _, adapter := range r.adapters {
		adapters = append(adapters, adapter)
	}
	return adapters
}

// Middlewares implements adapter.Registry.
func (r *AdapterRegistry) Middlewares() []adapter.Middleware {
	return r.middlewares
}

func NewAdapterRegistry() *AdapterRegistry {
	return &AdapterRegistry{
		adapters:    make(map[adapter.Protocol]adapter.Adapter),
		middlewares: make([]adapter.Middleware, 0),
	}
}

func (r *AdapterRegistry) Register(adapter adapter.Adapter) error {
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

func (r *AdapterRegistry) Unregister(protocol adapter.Protocol) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.adapters[protocol]; !exists {
		return fmt.Errorf("adapter for protocol %s not found", protocol)
	}
	delete(r.adapters, protocol)
	return nil
}

func (r *AdapterRegistry) GetAdapter(protocol adapter.Protocol) (adapter.Adapter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	adapter, ok := r.adapters[protocol]
	if !ok {
		return nil, fmt.Errorf("adapter for protocol %s not registered", protocol)
	}
	return adapter, nil
}

func (a *AdapterRegistry) AddMiddleware(middleware adapter.Middleware) {
	a.middlewares = append(a.middlewares, middleware)
}

package es

import "sync"

type ConverterFn func(CommandRecord) (ICommand, error)
type ConverterEventFn func(EventRecord) (IEvent, error)

type Registry struct {
	mutex    *sync.RWMutex
	commands map[string]ConverterFn
	events   map[string]ConverterEventFn
}

func NewRegistry() *Registry {
	return &Registry{
		mutex:    &sync.RWMutex{},
		commands: make(map[string]ConverterFn),
		events:   make(map[string]ConverterEventFn),
	}
}

func (r *Registry) RegisterCommand(name string, f ConverterFn) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.commands[name] = f
}

func (r *Registry) GetCommand(name string) (ConverterFn, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	f, ok := r.commands[name]
	return f, ok
}

func (r *Registry) RegisterEvent(name string, f ConverterEventFn) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.events[name] = f
}

func (r *Registry) GetEvent(name string) (ConverterEventFn, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	f, ok := r.events[name]
	return f, ok
}

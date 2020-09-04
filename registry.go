package registry

import (
	"sync"

	"github.com/jealone/sli4go"
)

var (
	defaultRegistry *Registry
	once            sync.Once
)

func GetRegistry() *Registry {
	once.Do(func() {
		sli4go.Fatalf("registry must initialize first")
	})
	return defaultRegistry
}

type Registry struct {
	Drivers map[string]Driver
	mu      sync.RWMutex
}

func (r *Registry) GetDriver(name string) Driver {
	r.mu.RLock()
	driver, ok := r.Drivers[name]
	r.mu.RUnlock()

	if ok {
		return driver
	}
	return nil
}

func (r *Registry) Bootstrap() {
	for _, driver := range r.Drivers {
		driver.Boot()
	}
}

type Driver interface {
	Boot()
	GetName() string
	GetEntry(key string) []byte
}

func InitRegistry(f func() (*Registry, error)) error {
	var err error

	once.Do(func() {
		defaultRegistry, err = f()
		defaultRegistry.Bootstrap()
	})

	return err
}

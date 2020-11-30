package registry

import (
	"fmt"
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

	Closers []Closer
	cmu     sync.RWMutex
}

type Closer interface {
	Close() error
}

type Decoder interface {
	Decode(interface{}) error
}

func (r *Registry) RegisterCloser(cl Closer) error {
	r.cmu.Lock()
	defer r.cmu.Unlock()

	r.Closers = append(r.Closers, cl)
	return nil
}

func (r *Registry) RegisterDriver(typ string, dec Decoder) (Driver, error) {

	var (
		driver Driver
		err    error
	)

	switch typ {
	case "file":
		conf := FileDriverConfig{}
		err = dec.Decode(&conf)
		if nil != err {
			return nil, fmt.Errorf("decode error: %w", err)
		}
		driver, err = NewYamlFileDriver(&conf)
		if nil != err {
			return nil, fmt.Errorf("file driver creation error: %w", err)
		}
	default:
		return nil, ErrUnknownDriver
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.Drivers[driver.GetName()]; ok {
		return nil, ErrDriverExists
	}
	driver.Boot()
	r.Drivers[driver.GetName()] = driver

	return driver, nil
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

func (r *Registry) Close() error {
	r.cmu.RLock()
	defer r.cmu.RUnlock()

	for _, c := range r.Closers {
		c.Close()
	}
}

type Driver interface {
	Boot()
	GetName() string
	GetEntry(key string) []byte
}

func InitRegistry(opts ...func() (*Registry, error)) error {

	if nil == opts || 0 == len(opts) {

		once.Do(func() {
			defaultRegistry = &Registry{}
			defaultRegistry.Drivers = make(map[string]Driver)
		})

		return nil

	} else {
		f := opts[0]
		var err error
		once.Do(func() {
			defaultRegistry, err = f()
			defaultRegistry.Bootstrap()
		})
		return err
	}

}

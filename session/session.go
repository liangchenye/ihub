package session

import (
	"errors"
	"fmt"
	"sync"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"

	"github.com/isula/ihub/config"
)

// Driver provides the session interfaces
type Driver interface {
	// Init setups a driver
	Init(paras map[string]interface{}) error
	// New creates a new session and gets the id string
	New(ctx context.Context, id string) (string, error)
	// Get returns a session by its id
	Get(ctx context.Context, id string) (interface{}, error)
	// Release frees the resource of a session by its id and removes the id
	Release(ctx context.Context, id string) error

	// Should merge to 'Get' and add 'name' paras
	// GetCache gets the session data by its id
	GetCache(ctx context.Context, id string) ([]byte, error)
	// PutCache puts the session data by its id
	PutCache(ctx context.Context, id string, data []byte) error
	// GC starts to garbage collection
	GC() error
}

var (
	sdLock sync.Mutex
	sds    = make(map[string]Driver, 16)

	sysSession Driver
)

func init() {
	sysSession = nil
}

// Register regists a driver
func Register(name string, driver Driver) error {
	if name == "" {
		return errors.New("could not register a session driver with empty name")
	}

	if driver == nil {
		return errors.New("could not register a nil session driver")
	}

	sdLock.Lock()
	defer sdLock.Unlock()

	if _, exists := sds[name]; exists {
		return fmt.Errorf("driver '%s' is already registered", name)
	}

	sds[name] = driver

	return nil
}

// InitSession inits the session defined in the config file
func InitSession(cfg config.SessionConfig) error {
	for n, v := range cfg {
		if d, ok := sds[n]; ok {
			logs.Debug("Init Session Driver: '%s'.", n)
			err := d.Init(v)
			if err == nil {
				sysSession = d
			}
			return err
		}
	}

	return errors.New("cannot find supported session driver")
}

// New creates a new session and gets the id string
func New(ctx context.Context, id string) (string, error) {
	if sysSession == nil {
		return "", errors.New("please init the session driver first")
	}

	return sysSession.New(ctx, id)
}

// Get returns a session by its id
func Get(ctx context.Context, id string) (interface{}, error) {
	if sysSession == nil {
		return nil, errors.New("please init the session driver first")
	}

	return sysSession.Get(ctx, id)
}

// Release frees the resource of a session by its id and removes the id
func Release(ctx context.Context, id string) error {
	if sysSession == nil {
		return errors.New("please init the session driver first")
	}

	return sysSession.Release(ctx, id)
}

// GetCache gets the session data by its id
func GetCache(ctx context.Context, id string) ([]byte, error) {
	if sysSession == nil {
		return nil, errors.New("please init the session driver first")
	}

	// FIXME: id should not be printed after service online
	logs.Debug("Session GetCache '%s'.", id)
	return sysSession.GetCache(ctx, id)
}

// PutCache puts the session data by its id
func PutCache(ctx context.Context, id string, data []byte) error {
	if sysSession == nil {
		return errors.New("please init the session driver first")
	}

	// FIXME: id should not be printed after service online
	logs.Debug("Session PutCache '%s'.", id)
	return sysSession.PutCache(ctx, id, data)
}

// GC starts to garbage collection
func GC() error {
	if sysSession == nil {
		return errors.New("please init the session driver first")
	}

	return sysSession.GC()
}

package memory

import (
	"errors"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/satori/go.uuid"

	"github.com/isula/ihub/session"
)

const (
	memoryPrefix = "inmemory"
)

// Memory provides a session driver with 'memory' backend
type Memory struct {
	store map[string]session.Record
	cache map[string][]byte
}

// Init creates two session maps
func (m *Memory) Init(paras map[string]interface{}) error {
	m.store = make(map[string]session.Record)
	m.cache = make(map[string][]byte)
	return nil
}

// New creates a new session
func (m *Memory) New(ctx context.Context, id string) (string, error) {
	if m.store == nil {
		return "", errors.New("please init the 'session memory' driver before use it")
	}

	if id != "" {
		if _, ok := m.store[id]; ok {
			return "", errors.New("id is already in used")
		}
		m.store[id] = session.NewRecordFromContext(ctx)
		return id, nil
	}

	sessionUUID := uuid.NewV4().String()
	m.store[sessionUUID] = session.NewRecordFromContext(ctx)
	return sessionUUID, nil
}

// Get returns the session interface by the id
//TODO: the return interface is not designed, useless?
func (m *Memory) Get(ctx context.Context, id string) (interface{}, error) {
	if m.store == nil {
		return nil, errors.New("please init the 'session memory' driver before use it")
	}

	r, ok := m.store[id]
	if !ok {
		return nil, errors.New("cannot get the matched sessionid")
	}

	err := r.Match(ctx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Release frees the session data by its id
func (m *Memory) Release(ctx context.Context, id string) error {
	if m.store == nil {
		return errors.New("please init the 'session memory' driver before use it")
	}

	if _, ok := m.store[id]; !ok {
		return errors.New("cannot get the matched sessionid")
	}

	delete(m.store, id)
	delete(m.cache, id)
	return nil
}

// GetCache gets the data from a session by its id
//TODO: id should be same with session id
func (m *Memory) GetCache(ctx context.Context, id string) ([]byte, error) {
	if m.cache == nil {
		return nil, errors.New("please init the 'session memory' driver before use it")
	}

	data, ok := m.cache[id]
	if !ok {
		return nil, errors.New("cannot get the matched cache data")
	}
	return data, nil
}

// PutCache puts the data from a session by its id
func (m *Memory) PutCache(ctx context.Context, id string, data []byte) error {
	if m.cache == nil {
		return errors.New("please init the 'session memory' driver before use it")
	}

	m.cache[id] = data
	return nil
}

// GC starts to garbage collection
func (m *Memory) GC() error {
	if m.store == nil {
		return errors.New("please init the 'session memory' driver before use it")
	}

	var expired []string
	for k, r := range m.store {
		if r.Expired() {
			expired = append(expired, k)
		}
	}

	num := len(expired)
	if num == 0 {
		return nil
	}

	logs.Info("Session GC start: '%d' expired session detected.", num)
	for _, id := range expired {
		delete(m.store, id)
	}

	return nil
}

func init() {
	if err := session.Register(memoryPrefix, &Memory{}); err != nil {
		logs.Error("Failed to register memory session driver.")
	} else {
		logs.Debug("Session driver '%s' registered.", memoryPrefix)
	}
}

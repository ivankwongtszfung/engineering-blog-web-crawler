package kvstore

import "errors"

type InMemoryStore struct {
	store map[string]any
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{store: make(map[string]any)}
}

func (m *InMemoryStore) Ping() error {
	return nil
}

func (m *InMemoryStore) Set(key string, value any) error {
	m.store[key] = value
	return nil
}

func (m *InMemoryStore) Exist(key string) (bool, error) {
	_, err := m.Get(key)
	return err == nil, nil
}

func (m *InMemoryStore) Get(key string) (any, error) {
	value, ok := m.store[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return value, nil
}

func (m *InMemoryStore) Delete(key string) error {
	delete(m.store, key)
	return nil
}

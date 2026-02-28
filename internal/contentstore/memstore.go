package contentstore

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/ashirt-ops/ashirt-server/internal/errors"
	"github.com/google/uuid"
)

// MemStore is the backing structure needed to interact with local memory -- for unit/integration
// testing purposes only
type MemStore struct {
	content map[string][]byte
	mutex   *sync.Mutex
}

// NewMemStore is the constructor for MemStore
func NewMemStore() (*MemStore, error) {
	m := MemStore{
		content: make(map[string][]byte),
		mutex:   new(sync.Mutex),
	}

	return &m, nil
}

// Upload stores content in memory
func (d *MemStore) Upload(data io.Reader) (key string, err error) {
	key = uuid.New().String()
	err = d.UploadWithName(key, data)

	if err != nil {
		err = errors.WrapError("Unable to add to MemStore", err)
	}

	return
}

// UploadWithName writes the given data to the given memory key -- this
// can allow for re-writing/replacing data if names are not unique
//
// Note: to avoid overwriting random keys, DO NOT use uuids as they key
func (d *MemStore) UploadWithName(key string, data io.Reader) error {
	b, err := io.ReadAll(data)
	if err != nil {
		return errors.WrapError("Unable upload with a given name to MemStore", err)
	}

	d.mutex.Lock()
	d.content[key] = b
	d.mutex.Unlock()
	return nil
}

func (d *MemStore) Read(key string) (io.Reader, error) {
	d.mutex.Lock()
	data, ok := d.content[key]
	d.mutex.Unlock()
	if !ok {
		return nil, errors.WrapError("Unable to read from MemStore", fmt.Errorf("No such key"))
	}
	return bytes.NewReader(data), nil
}

// Delete removes files in in your OS's temp directory
func (d *MemStore) Delete(key string) error {
	d.mutex.Lock()
	if _, ok := d.content[key]; !ok { // artificial behavior to match other stores
		return errors.WrapError("Unable to delete from MemStore", fmt.Errorf("No such key"))
	}
	delete(d.content, key)
	d.mutex.Unlock()
	return nil
}

func (d *MemStore) Name() string {
	return "memory"
}

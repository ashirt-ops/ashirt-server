// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package contentstore

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"sync"

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
	b, err := ioutil.ReadAll(data)
	if err != nil {
		return
	}

	key = uuid.New().String()
	d.mutex.Lock()
	d.content[key] = b
	d.mutex.Unlock()
	return
}

func (d *MemStore) Read(key string) (io.Reader, error) {
	data, ok := d.content[key]
	if !ok {
		return nil, fmt.Errorf("No such key")
	}
	return bytes.NewReader(data), nil
}

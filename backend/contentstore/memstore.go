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
	key = uuid.New().String()
	err = d.UploadWithName(key, data)
	return
}

// UploadWithName writes the given data to the given memory key -- this
// can allow for re-writing/replacing data if names are not unique
//
// Note: to avoid overwriting random keys, DO NOT use uuids as they key
// Note 2: This is NOT part of the standard ContentStore interface
func (d *MemStore) UploadWithName(key string, data io.Reader) error {
	b, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}

	d.mutex.Lock()
	d.content[key] = b
	d.mutex.Unlock()
	return nil
}

func (d *MemStore) Read(key string) (io.Reader, error) {
	data, ok := d.content[key]
	if !ok {
		return nil, fmt.Errorf("No such key")
	}
	return bytes.NewReader(data), nil
}

// Delete removes files in in your OS's temp directory
func (d *MemStore) Delete(key string) error {
	if _, ok := d.content[key]; !ok { // artificial behavior to match other stores
		return fmt.Errorf("No such key")
	}
	delete(d.content, key)
	return nil
}

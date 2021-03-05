// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package contentstore

import (
	"bufio"
	"io"
	"os"
	"path"

	"github.com/theparanoids/ashirt-server/backend"
)

// DevStore is the backing structure needed to interact with a local, temporary image store
type DevStore struct {
	dir string
}

// NewDevStore constructs a basic DevStore / local temporary store
func NewDevStore() (*DevStore, error) {
	tmpDir := "/tmp/contentstore"
	tmpDirInfo, err := os.Stat(tmpDir)
	if err != nil || !tmpDirInfo.IsDir() {
		tmpDir, err = os.MkdirTemp("", "ashirt")
		if err != nil {
			return nil, backend.WrapError("Unable to establish a DevStore", err)
		}
	}
	return &DevStore{dir: tmpDir}, nil
}

// Upload stores files in your OS's temp directory
func (d *DevStore) Upload(data io.Reader) (string, error) {
	file, err := os.CreateTemp(d.dir, "")
	if err != nil {
		return "", backend.WrapError("Unable to upload to DevStore", err)
	}
	defer file.Close()
	_, err = bufio.NewReader(data).WriteTo(file)
	if err != nil {
		return "", backend.WrapError("Unable upload to DevStore", err)
	}
	return path.Base(file.Name()), nil
}

func (d *DevStore) Read(key string) (io.Reader, error) {
	reader, err := os.Open(path.Join(d.dir, path.Clean(key)))
	if err != nil {
		return reader, backend.WrapError("Unable to read file from DevStore", err)
	}
	return reader, err
}

// Delete removes files in in your OS's temp directory
func (d *DevStore) Delete(key string) error {
	err := os.Remove(path.Join(d.dir, path.Clean(key)))
	if err != nil {
		return backend.WrapError("Unable to delete file from DevStore", err)
	}
	return err
}

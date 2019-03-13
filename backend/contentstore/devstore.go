// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package contentstore

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path"
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
		tmpDir, err = ioutil.TempDir("", "ashirt")
		if err != nil {
			return nil, err
		}
	}
	return &DevStore{dir: tmpDir}, nil
}

// Upload stores image files in /tmp (or whatever your OS considers to be a temporary file)
func (d *DevStore) Upload(data io.Reader) (string, error) {
	file, err := ioutil.TempFile(d.dir, "")
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = bufio.NewReader(data).WriteTo(file)
	if err != nil {
		return "", err
	}
	return path.Base(file.Name()), nil
}

func (d *DevStore) Read(key string) (io.Reader, error) {
	return os.Open(path.Join(d.dir, path.Clean(key)))
}

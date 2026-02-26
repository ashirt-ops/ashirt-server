package contentstore

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/ashirt-ops/ashirt-server/internal/errors"
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
			return nil, errors.WrapError("Unable to establish a DevStore", err)
		}
	}
	return &DevStore{dir: tmpDir}, nil
}

// Upload stores files in your OS's temp directory
func (d *DevStore) Upload(data io.Reader) (string, error) {
	file, err := os.CreateTemp(d.dir, "")
	if err != nil {
		return "", errors.WrapError("Unable to upload to DevStore", err)
	}
	defer file.Close()
	_, err = bufio.NewReader(data).WriteTo(file)
	if err != nil {
		return "", errors.WrapError("Unable upload to DevStore", err)
	}
	return path.Base(file.Name()), nil
}

// UploadWithName is unsupported for the devstore.
func (d *DevStore) UploadWithName(key string, data io.Reader) error {
	return fmt.Errorf("UploadWithName is Unsupported")
}

func (d *DevStore) Read(key string) (io.Reader, error) {
	reader, err := os.Open(path.Join(d.dir, path.Clean(key)))
	if err != nil {
		return reader, errors.WrapError("Unable to read file from DevStore", err)
	}
	return reader, err
}

// Delete removes files in in your OS's temp directory
func (d *DevStore) Delete(key string) error {
	err := os.Remove(path.Join(d.dir, path.Clean(key)))
	if err != nil {
		return errors.WrapError("Unable to delete file from DevStore", err)
	}
	return err
}

func (d *DevStore) Name() string {
	return "local"
}

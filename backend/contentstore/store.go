// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package contentstore

import (
	"io"
)

// Store provides a generic interface into interacting with the underlying storage service
//
// Upload stores the provided file/bytes into the storage service, returning the location of that
// file or any error that may have occurred
//
// Note that UploadWithName is only intended for development and testing. This should not be used
// directly.
//
// Read retrieves the raw bytes from the storage service, given a key obtained by Upload
type Store interface {
	Upload(data io.Reader) (string, error)
	UploadWithName(key string, data io.Reader) error
	Read(key string) (io.Reader, error)
	Delete(key string) error
	Name() string
}

type ProductionStore interface {
	SendURL(key string) string
}

// ContentKeys stores the location/path of the original content, as well as the thumbnail/preview location
type ContentKeys struct {
	Full      string
	Thumbnail string
}

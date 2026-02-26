package contentstore

import (
	"io"
)

type blobStorable struct {
	data io.Reader
}

// NewBlob returns a Storable that deals with binary/non-binary content that does not lend itself to a preview/proxy
func NewBlob(data io.Reader) Storable {
	return blobStorable{
		data: data,
	}
}

// ProcessPreviewAndUpload uploads the blob as the full/master content, and uses that key for the
// Thumbnail/preview/proxy version
func (blob blobStorable) ProcessPreviewAndUpload(s Store) (ContentKeys, error) {
	contentKeys := ContentKeys{}

	var err error
	contentKeys.Full, err = s.Upload(blob.data)
	if err != nil {
		return contentKeys, err
	}
	contentKeys.Thumbnail = contentKeys.Full

	return contentKeys, nil
}

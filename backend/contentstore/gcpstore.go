// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package contentstore

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	"github.com/ashirt-ops/ashirt-server/backend"
	"github.com/google/uuid"
)

type GCPStore struct {
	bucketName      string
	gcpClient       *storage.Client
	bucketAccess    *storage.BucketHandle
	creationContext context.Context
}

// NewGCPStore provides a mechanism to initialize a GCP client
func NewGCPStore(bucketName string) (*GCPStore, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, backend.WrapError("Unable to establish an gcp session", err)
	}
	return &GCPStore{
		bucketName:      bucketName,
		gcpClient:       client,
		bucketAccess:    client.Bucket(bucketName),
		creationContext: ctx,
	}, nil
}

// Upload stores a file in the Google Cloud bucket configured when the gcpStore was created
func (s *GCPStore) Upload(data io.Reader) (string, error) {
	key := uuid.New().String()

	err := s.UploadWithName(key, data)

	return key, err
}

func (d *GCPStore) SendURL(key string) string {
	return "thing" //path.Join(d.dir, path.Clean(key))
}

// UploadWithName is a test/dev helper that places a file on Google Cloud with a given name
// This is not intended for general use.
func (s *GCPStore) UploadWithName(key string, data io.Reader) error {
	ctx := context.Background()
	wc := s.bucketAccess.Object(key).NewWriter(ctx)

	if _, err := io.Copy(wc, data); err != nil {
		return backend.WrapError("Upload to gcp failed", err)
	}

	if err := wc.Close(); err != nil {
		return backend.WrapError("Unable to close gcp writer", err)
	}

	return nil
}

// Read retrieves the indicated file from Google Cloud
func (s *GCPStore) Read(key string) (io.Reader, error) {
	ctx := context.Background()
	res, err := s.bucketAccess.Object(key).NewReader(ctx)
	if err != nil {
		return nil, backend.WrapError("Unable to read from gcp", err)
	}
	return res, nil
}

// Delete removes the indicated file from GCP
func (s *GCPStore) Delete(key string) error {
	ctx := context.Background()
	err := s.bucketAccess.Object(key).Delete(ctx)

	if err != nil {
		return backend.WrapError("Delete from gcp failed", err)
	}

	return nil
}

func (s *GCPStore) Name() string {
	return "gcp"
}

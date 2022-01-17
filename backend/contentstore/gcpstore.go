// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package contentstore

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/theparanoids/ashirt-server/backend"
)

type gcpStore struct {
	bucketName      string
	gcpClient       *storage.Client
	bucketAccess    *storage.BucketHandle
	creationContext context.Context
}

// NewGCPStore provides a mechanism to initialize a GCP client
func NewGCPStore(bucketName string) (*gcpStore, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, backend.WrapError("Unable to establish an gcp session", err)
	}
	return &gcpStore{
		bucketName:      bucketName,
		gcpClient:       client,
		bucketAccess:    client.Bucket(bucketName),
		creationContext: ctx,
	}, nil
}

// Upload stores a file in the Google Cloud bucket configured when the gcpStore was created
func (s *gcpStore) Upload(data io.Reader) (string, error) {
	key := uuid.New().String()

	ctx := context.Background()
	wc := s.bucketAccess.Object(key).NewWriter(ctx)

	if _, err := io.Copy(wc, data); err != nil {
		return key, backend.WrapError("Upload to gcp failed", err)
	}

	if err := wc.Close(); err != nil {
		return key, backend.WrapError("Unable to close gcp writer", err)
	}
	
	// TODO: figure out how to properly do ACL for gcp
	// acl := s.bucketAccess.Object(key).ACL()
	// err := acl.Set(ctx, storage., storage.ScopeFullControl)
	// if err != nil {
	// 	return key, backend.WrapError("Unable to set GCP ACLs", err)
	// }

	return key, nil
}

// Read retrieves the indicated file from Google Cloud
func (s *gcpStore) Read(key string) (io.Reader, error) {
	ctx := context.Background()
	res, err := s.bucketAccess.Object(key).NewReader(ctx)
	if err != nil {
		return nil, backend.WrapError("Unable to read from gcp", err)
	}
	return res, nil
}

// Delete removes the indicated file from GCP
func (s *gcpStore) Delete(key string) error {
	ctx := context.Background()
	err := s.bucketAccess.Object(key).Delete(ctx)

	if err != nil {
		return backend.WrapError("Delete from gcp failed", err)
	}

	return nil
}

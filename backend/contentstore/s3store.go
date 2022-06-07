// Copyright 2022, Yahoo, Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package contentstore

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/theparanoids/ashirt-server/backend"
)

// S3Store is the backing structure needed to interact with an Amazon S3 storage service
// TODO: this can be unexported
type S3Store struct {
	bucketName string
	s3Client   *s3.Client
}

// NewS3Store provides a mechanism to initialize an S3 bucket in a particular region
func NewS3Store(bucketName string, region string) (*S3Store, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, backend.WrapError("Unable to establish an s3 session", err)
	}
	return &S3Store{
		bucketName: bucketName,
		s3Client:   s3.NewFromConfig(cfg),
	}, nil
}

// Upload stores a file in the Amazon S3 bucket configured when the S3 store was created
func (s *S3Store) Upload(data io.Reader) (string, error) {
	key := uuid.New().String()

	err := s.UploadWithName(key, data)

	return key, err
}

// UploadWithName is a test/dev helper that places a file on S3 with a given name
// This is not intended for general use.
func (s *S3Store) UploadWithName(key string, data io.Reader) error {
	_, err := s.s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		ACL:    "bucket-owner-full-control",
		Body:   data,
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return backend.WrapError("Upload to s3 failed", err)
	}

	return nil
}

// Read retrieves the indicated file from Amazon S3
func (s *S3Store) Read(key string) (io.Reader, error) {
	res, err := s.s3Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, backend.WrapError("Unable to read from s3", err)
	}
	return res.Body, nil
}

// Delete removes files in in your OS's temp directory
func (s *S3Store) Delete(key string) error {
	_, err := s.s3Client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return backend.WrapError("Delete from s3 failed", err)
	}

	return nil
}

func (d *S3Store) Name() string {
	return "s3"
}

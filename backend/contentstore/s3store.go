// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package contentstore

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

// S3Store is the backing structure needed to interact with an Amazon S3 storage service
// TODO: this can be unexported
type S3Store struct {
	bucketName string
	s3Client   *s3.S3
}

// NewS3Store provides a mechanism to initialize an S3 bucket in a particular region
func NewS3Store(bucketName string, region string) (*S3Store, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	return &S3Store{
		bucketName: bucketName,
		s3Client:   s3.New(sess, &aws.Config{Region: &region}),
	}, nil
}

// Upload stores a file in the Amazon S3 bucket configured when the S3 store was created
func (s *S3Store) Upload(data io.Reader) (string, error) {
	key := uuid.New().String()

	_, err := s.s3Client.PutObject(&s3.PutObjectInput{
		ACL:    aws.String("bucket-owner-full-control"),
		Body:   aws.ReadSeekCloser(data),
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})

	return key, err
}

// Read retrieves the indicated file from Amazon S3
func (s *S3Store) Read(key string) (io.Reader, error) {
	res, err := s.s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

// Delete removes files in in your OS's temp directory
func (s *S3Store) Delete(key string) error {
	_, err := s.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	return err
}

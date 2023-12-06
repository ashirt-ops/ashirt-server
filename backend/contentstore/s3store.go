package contentstore

import (
	"io"
	"time"

	"github.com/ashirt-ops/ashirt-server/backend"
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
		return nil, backend.WrapError("Unable to establish an s3 session", err)
	}
	return &S3Store{
		bucketName: bucketName,
		s3Client:   s3.New(sess, &aws.Config{Region: &region}),
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
	_, err := s.s3Client.PutObject(&s3.PutObjectInput{
		ACL:    aws.String("bucket-owner-full-control"),
		Body:   aws.ReadSeekCloser(data),
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
	res, err := s.s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, backend.WrapError("Unable to read from s3", err)
	}
	return res.Body, nil
}

type URLData struct {
	Url            string    `json:"url"`
	ExpirationTime time.Time `json:"expirationTime"`
}

func (s *S3Store) SendURLData(key string) (*URLData, error) {
	contentType := "image/jpeg"
	req, _ := s.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket:              aws.String(s.bucketName),
		Key:                 aws.String(key),
		ResponseContentType: aws.String(contentType),
	})

	minutes := time.Minute * time.Duration(30)
	url, err := req.Presign(minutes)
	if err != nil {
		return nil, backend.WrapError("Unable to get presigned URL", err)
	}
	data := URLData{
		Url:            url,
		ExpirationTime: time.Now().UTC().Add(minutes),
	}

	return &data, nil
}

// Delete removes files in in your OS's temp directory
func (s *S3Store) Delete(key string) error {
	_, err := s.s3Client.DeleteObject(&s3.DeleteObjectInput{
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

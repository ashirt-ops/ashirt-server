package contentstore

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Presigner struct {
	PresignClient *s3.PresignClient
}

func (presigner Presigner) GetObject(
	bucketName string, objectKey string, minutes time.Duration) (*v4.PresignedHTTPRequest, error) {
	contentType := "image/jpeg"
	request, err := presigner.PresignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket:              aws.String(bucketName),
		Key:                 aws.String(objectKey),
		ResponseContentType: aws.String(contentType),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = minutes
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to get %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return request, err
}

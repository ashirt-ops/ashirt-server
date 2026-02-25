package contentstore

import (
	"context"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Presigner struct {
	PresignClient *s3.PresignClient
	Logger        *slog.Logger
}

func (presigner Presigner) GetObject(
	ctx context.Context, bucketName string, objectKey string, minutes time.Duration) (*v4.PresignedHTTPRequest, error) {
	contentType := "image/jpeg"
	request, err := presigner.PresignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket:              aws.String(bucketName),
		Key:                 aws.String(objectKey),
		ResponseContentType: aws.String(contentType),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = minutes
	})
	if err != nil {
		presigner.Logger.ErrorContext(ctx, "Couldn't get a presigned request", "bucket", bucketName, "key", objectKey, "error", err)
	}
	return request, err
}
